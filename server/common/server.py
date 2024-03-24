import socket
import logging
import signal
import time
from common import communication
from common.utils import store_bets, load_bets, has_won

HEADER_LENGHT = 4

class Server:
    def __init__(self, port, listen_backlog):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)

        # Dictionary that will have agency number as keys and client sockets as values
        self.clients = dict()

        # Boolean to stop the server gracefully
        self._stop_server = False 

        # Declare the SIGTERM handler
        signal.signal(signal.SIGTERM, self.__exit_gracefully)

    def run(self):
        """
        Dummy Server loop

        Server that accept a new connections and establishes a
        communication with a client. After client with communucation
        finishes, servers starts to accept new connections again
        """
        clients_accepted = 0
        while not self._stop_server or clients_accepted != 5:
            try:
                agency, client_sock = self.__accept_new_connection()
                self.__handle_client_connection(client_sock, agency)
            except socket.timeout:
                # In case the client_sock wasn't closed because of the exception
                client_sock.close()
                self._server_socket.close()
                continue

            self.clients[agency] = client_sock
            clients_accepted += 1

        logging.info('action: sorteo | result: success')

        self.start_lottery()

        logging.info('action: server_finished | result: success')

    def __handle_client_connection(self, client_sock, agency):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        addr = client_sock.getpeername()

        eof = False
        while eof != True:
            # Read message
            msg, err = communication.read_socket(client_sock)
            if err is not None:
                logging.error(f'action: read_socket | result: fail | ip: {addr[0]} | error: {err}')
                client_sock.close()
                return
            elif msg == "EOF":
                logging.info(f'action: finish_loop | result: success | ip: {addr[0]}')
                break
                

            # Deserialize message
            bets, err = communication.deserialize(msg, agency)
            if err is not None:
                logging.error(f'action: deserialize | result: fail | ip: {addr[0]} | message: {msg} |error: {err}')
                client_sock.close()
                return

            # Store the bet
            store_bets(bets)
            logging.info(f'action: apuestas_almacenada | result: success | ip: {addr[0]}')

            # Send ack
            msg = f'ACK'
            err = communication.write_socket(client_sock, msg)
            if err is not None:
                logging.error(f'action: send_ack | result: fail | ip: {addr[0]} | error: {err}')
                client_sock.close()
                return
            logging.info(f'action: send_ack | result: success | ip: {addr[0]}')
        

    def __accept_new_connection(self):
        """
        Accept new connections

        Function blocks until a connection to a client is made.
        Then connection created is printed and returned
        """

        # Connection arrived
        logging.info('action: accept_connections | result: in_progress')
        c, addr = self._server_socket.accept()
        logging.info(f'action: accept_connections | result: success | ip: {addr[0]}')

        # Read the number of the agency
        agency, _ = communication.read_socket(c)

        return agency, c
    
    def start_lottery(self):

        for bet in load_bets():
            if has_won(bet):
                #notify its agency
                communication.write_socket()
        



    def __exit_gracefully(self, *args):
        """
        Handles SIGTERM

        By setting self._stop_server to False, the server will continue with the iteration
        it was working, but it will be his last one before stopping gracefully. 
        """
        self._stop_server = True
