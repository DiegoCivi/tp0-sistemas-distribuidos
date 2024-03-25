""" Lenght of the communication headers protcol"""
HEADER_LENGHT = 5
MSG_SIZE_LENGTH = 4


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

"""
Sends a short meesage to the client. Only the header with the end_flag set to True
"""
def sendEOF(socket):
    try: 
        header = get_header("", "1")        
        _handle_short_write(socket, header)
        return None
    except Exception as e:
        return e