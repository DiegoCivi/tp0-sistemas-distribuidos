from common.utils import Bet
import logging


""" Indexes for bet atributes in the deserialization """
BYTES_LEN_INDEX = 0
AGENCY_INDEX = 0
FIRST_NAME_INDEX = 1
SECOND_NAME_INDEX = 2
DOCUMENT_INDEX = 3
BIRTHDATE_INDEX = 4
NUMBER_INDEX = 5
""" Lenght of the communication headers protcol"""
HEADER_LENGHT = 4

"""
Deserealizes a message and creates a Bet
"""
def deserialize(msg):
    splitted_msg = msg.split('/')

    # bytes_to_read = splitted_msg[BYTES_LEN_INDEX]

    agency = splitted_msg[AGENCY_INDEX]
    first_name = splitted_msg[FIRST_NAME_INDEX]
    second_name = splitted_msg[SECOND_NAME_INDEX]
    document = splitted_msg[DOCUMENT_INDEX]
    birthdate = splitted_msg[BIRTHDATE_INDEX]
    number = splitted_msg[NUMBER_INDEX]

    try:
        return (Bet(agency, first_name, second_name, document, birthdate, number), None)
    except ValueError as e:
        return (None, e)
    

def read_socket(socket):
    header = socket.recv(HEADER_LENGHT).decode('utf-8') # lE BORRE EL rstrip(), fijjasrse si eso no rompe
    logging.info(f'Se recibio el header: {header}')
    msg_len = int(header)

    msg = socket.recv(msg_len).rstrip().decode('utf-8')
    
    logging.info(f'action: receive_message | result: success | msg: {msg}')

    # Deserialization of the message
    return deserialize(msg)

