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
Deserealizes a message and creates a Bet.
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

"""
Reads from the received socket. It supports short-read.
"""
def read_socket(socket):
    try: 
        # Read header
        header = _handle_short_read(socket, HEADER_LENGHT)

        # Read message
        msg_len = int(header)
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
        msg += socket.recv(bytes_to_read - bytes_read).rstrip().decode('utf-8')
        bytes_read += len(msg)
    
    return msg


"""
Writes into the received socket. It supports short-write.
"""
def write_socket(socket, msg):
    try: 
        # Add header
        header = get_header(msg)
        complete_msg = header + msg
        
        _handle_short_write(socket, complete_msg, len(complete_msg))

        return None
    
    except Exception as e:
        return e

"""
If the socket.send() call does not write the whole message, 
it sends again from the first byte it did not sent.
"""
def _handle_short_write(socket, msg, bytes_to_write):
    sent_bytes = socket.send(msg.encode("utf-8"))
    while sent_bytes < bytes_to_write:
        sent_bytes += socket.send(msg[sent_bytes + 1:].encode("utf-8"))

"""
Returns the protocols header for a message
"""
def get_header(msg):
    msg_len = str(len(msg))
    msg_len_bytes = len(msg_len)

    for _ in range(0, HEADER_LENGHT - msg_len_bytes):
        msg_len = '0' + msg_len

    return msg_len