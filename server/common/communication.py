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
    

#def read_socket(socket):
#    header = socket.recv(HEADER_LENGHT).decode('utf-8') # lE BORRE EL rstrip(), fijarse si eso no rompe
#    logging.info(f'Se recibio el header: {header}')
#    msg_len = int(header)
#
#    msg = socket.recv(msg_len).rstrip().decode('utf-8')
#    
#    logging.info(f'action: receive_message | result: success | msg: {msg}')
#
#    # Deserialization of the message
#    return deserialize(msg)

def read_socket(socket, bytes_ro_read):
    msg = socket.recv(bytes_ro_read).rstrip().decode('utf-8')
    
    logging.info(f'action: receive_message | result: success | msg: {msg}')

    return msg

def write_socket(socket, msg):
    # Add header
    header = get_header(msg)
    complete_msg = header + msg

    temp = socket.send(complete_msg.encode("utf-8"))

    logging.info(f'action: write_message | result: success | msg: {complete_msg}')
    
    return temp

def get_header(msg):
    msg_len = str(len(msg))
    msg_len_bytes = len(msg_len)

    for _ in range(0, HEADER_LENGHT - msg_len_bytes):
        msg_len = '0' + msg_len

    return msg_len