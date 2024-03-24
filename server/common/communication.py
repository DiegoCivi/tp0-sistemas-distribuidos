from common.utils import Bet
import logging


""" Indexes for bet atributes in the deserialization """
BYTES_LEN_INDEX = 0
FIRST_NAME_INDEX = 0
SECOND_NAME_INDEX = 1
DOCUMENT_INDEX = 2
BIRTHDATE_INDEX = 3
NUMBER_INDEX = 4
""" Lenght of the communication headers protcol"""
HEADER_LENGHT = 5
MSG_SIZE_LENGTH = 4
""" Separators """
LINE_SEPARATOR = '/'
INFO_SEPARATOR = ','

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

"""
Reads from the received socket. It supports short-read.
"""
def read_socket(socket):
    try: 
        # Read header
        header = _handle_short_read(socket, HEADER_LENGHT)

        # Read message
        msg_len = int(header[:-1])
        end_flag = int(header[-1])
        if end_flag == 1:
            return "EOF", None
        bet_msg = _handle_short_read(socket, msg_len)

        return bet_msg, None
    
    except Exception as e:
        return None, e

"""
Handler of the short-read. Called by read_socket().
"""
def _handle_short_read(socket, bytes_to_read):
    bytes_read = 0
    msg = ""
    while bytes_read < bytes_to_read:
        msg_bytes = socket.recv(bytes_to_read - bytes_read) #.rstrip()
        bytes_read += len(msg_bytes)
        msg += msg_bytes.decode('utf-8')
    
    return msg

"""
Writes into the received socket. It supports short-write.
"""
def write_socket(socket, msg):
    try: 
        # Add header
        header = get_header(msg, "0")
        complete_msg = header + msg
        
        _handle_short_write(socket, complete_msg)

        return None
    
    except Exception as e:
        return e

"""
If the socket.send() call does not write the whole message, 
it sends again from the first byte it did not sent.
"""
def _handle_short_write(socket, msg):
    msg_bytes = msg.encode("utf-8")
    bytes_to_write = len(msg_bytes)
    sent_bytes = socket.send(msg_bytes)
    while sent_bytes < bytes_to_write:
        sent_bytes += socket.send(msg_bytes[sent_bytes:])

"""
Returns the protocols header for a message
"""
def get_header(msg, end_flag):
    header = str(len(msg))
    msg_len_bytes = len(header)

    for _ in range(0, MSG_SIZE_LENGTH - msg_len_bytes):
        header = '0' + header

    header += end_flag

    return header

def sendEOF(socket):
    try: 
        header = get_header("", "1")        
        _handle_short_write(socket, header)
        return None
    except Exception as e:
        return e