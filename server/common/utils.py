import csv
import datetime
from common.communication import read_socket, write_socket
import logging

""" Bets storage location. """
STORAGE_FILEPATH = "./bets.csv"
""" Simulated winner number in the lottery contest. """
LOTTERY_WINNER_NUMBER = 7574
""" Separators """
LINE_SEPARATOR = '/'
INFO_SEPARATOR = ','
""" Indexes for bet atributes in the deserialization """
BYTES_LEN_INDEX = 0
FIRST_NAME_INDEX = 0
SECOND_NAME_INDEX = 1
DOCUMENT_INDEX = 2
BIRTHDATE_INDEX = 3
NUMBER_INDEX = 4



""" A lottery bet registry. """
class Bet:
    def __init__(self, agency: str, first_name: str, last_name: str, document: str, birthdate: str, number: str):
        """
        agency must be passed with integer format.
        birthdate must be passed with format: 'YYYY-MM-DD'.
        number must be passed with integer format.
        """
        self.agency = int(agency)
        self.first_name = first_name
        self.last_name = last_name
        self.document = document
        self.birthdate = datetime.date.fromisoformat(birthdate)
        self.number = int(number)

    def serialize(self):
        return f'{self.agency}/{self.first_name}/{self.last_name}/{self.document}/{self.birthdate}/{self.number}'

"""
Deserealizes a message and creates a list with Bets.
"""
def deserialize(msg, agency):
    # The message is first splitted to have the different bets
    splitted_lines = msg.split(LINE_SEPARATOR)
    bets = []

    # The last item in the splitted_lines list is ignored since its a empty string
    # This happens because the last character of the message received is a '/' 
    for bet_msg in splitted_lines[:-1]:
        splitted_msg = bet_msg.split(INFO_SEPARATOR)

        first_name = splitted_msg[FIRST_NAME_INDEX]
        second_name = splitted_msg[SECOND_NAME_INDEX]
        document = splitted_msg[DOCUMENT_INDEX]
        birthdate = splitted_msg[BIRTHDATE_INDEX]
        number = splitted_msg[NUMBER_INDEX]

        try:
            bet = Bet(agency, first_name, second_name, document, birthdate, number)
            bets.append(bet)
        except ValueError as e:
            return (None, e)
    
    return (bets, None)

""" Checks whether a bet won the prize or not. """
def has_won(bet: Bet) -> bool:
    return bet.number == LOTTERY_WINNER_NUMBER

"""
Persist the information of each bet in the STORAGE_FILEPATH file.
Not thread-safe/process-safe.
"""
def store_bets(bets: list[Bet]) -> None:
    with open(STORAGE_FILEPATH, 'a+') as file:
        writer = csv.writer(file, quoting=csv.QUOTE_MINIMAL)
        for bet in bets:
            writer.writerow([bet.agency, bet.first_name, bet.last_name,
                             bet.document, bet.birthdate, bet.number])

"""
Loads the information all the bets in the STORAGE_FILEPATH file.
Not thread-safe/process-safe.
"""
def load_bets() -> list[Bet]: # type: ignore
    with open(STORAGE_FILEPATH, 'r') as file:
        reader = csv.reader(file, quoting=csv.QUOTE_MINIMAL)
        for row in reader:
            yield Bet(row[0], row[1], row[2], row[3], row[4], row[5])

def handle_message(client_sock, addr, agency, sem):
    # Read message
    #logging.info(f"Se va a leer un mensaje de la agencia: {agency}")
    msg, err = read_socket(client_sock)
    #logging.info(f"Se leyo un mensaje de la agencia: {agency} mensaje de largo: {len(msg)}")
    if err is not None:
        logging.error(f'action: read_socket | result: fail | ip: {addr[0]} | error: {err}')
        client_sock.close()
        return False, err
    elif msg == "EOF":
        logging.info(f'action: finish_loop | result: success | ip: {addr[0]}')
        return True, None
    #logging.info(f"Se va a des-serializar un mensaje de la agencia: {agency}")
    # Deserialize message
    bets, err = deserialize(msg, agency)
    if err is not None:
        logging.error(f'action: deserialize | result: fail | ip: {addr[0]} | message: {msg} |error: {err}')
        client_sock.close()
        return False, err
    #logging.info(f"Se va escribir las apuestas de la agencia: {agency}")
    sem.acquire()
    # Store the bet
    store_bets(bets)
    sem.release()
    #logging.info(f"Se escribieron las apuestas de la agencia: {agency}")

    # Send ack
    msg = f'ACK'
    err = write_socket(client_sock, msg)
    if err is not None:
        logging.error(f'action: send_ack | result: fail | ip: {addr[0]} | error: {err}')
        client_sock.close()
        return False, err
    
    return False, None

def handle_client(agency, client_sock, send_EOF, rec_winner, sem):
    #logging.info(f"Soy el proceso de la agencia: {agency}")
    addr = client_sock.getpeername()

    eof = False
    while eof != True:
        eof, err = handle_message(client_sock, addr, agency, sem)
    
    # After reading all the bets, the "father" process is notified
    send_EOF.send("EOF")
    send_EOF.close()

    # Winners are received from the "father" process and sent to the client agency
    winner_doc = None
    winners_quantity = 0
    while winner_doc != "EOF":
        winner_doc = rec_winner.recv()
        err = write_socket(client_sock, winner_doc)
        if err is not None:
            logging.error(f'action: send_ack | result: fail | ip: {addr[0]} | error: {err}')
            client_sock.close()
            return
        
        winners_quantity += 1
    
    # Close socket and pipe
    rec_winner.close()
    client_sock.close()

    
