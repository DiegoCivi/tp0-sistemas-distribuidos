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
        while not self._stop_server and clients_accepted != 5:
            try:
                agency, client_sock = self.__accept_new_connection()
                self.__handle_client_connection(client_sock, agency)
            except OSError as e:
                # In case the client_sock wasn't closed because of the exception
                client_sock.close()
                self.__close_clients_socks()
                self._server_socket.close()
                continue

            self.clients[int(agency)] = client_sock
            clients_accepted += 1

        if self._stop_server:
            return

        logging.info('action: sorteo | result: success')

        self.start_lottery()

        self.__close_clients_socks()
        logging.info('action: server_finished | result: success')

    def start_lottery(self):
        """
        For each winner found, a message with the document is sent to the corresponding agency.
        Afeter thath, we notify the agencys there are no more winners to send.
        """

        for bet in load_bets():
            if has_won(bet):
                client_sock = self.clients[bet.agency]
                err = communication.write_socket(client_sock, bet.document)
                if err != None:
                    logging.error(f"action: send_winner | result: fail | error: {err}")
                    return        

        for sock in self.clients.values():
            communication.sendEOF(sock)

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
                eof = True
                continue
                

            # Deserialize message
            bets, err = communication.deserialize(msg, agency)
            if err is not None:
                logging.error(f'action: deserialize | result: fail | ip: {addr[0]} | message: {msg} |error: {err}')
                client_sock.close()
                return

            # Store the bet
            store_bets(bets)

            # Send ack
            msg = f'ACK'
            err = communication.write_socket(client_sock, msg)
            if err is not None:
                logging.error(f'action: send_ack | result: fail | ip: {addr[0]} | error: {err}')
                client_sock.close()
                return
        

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
        
    def __close_clients_socks(self):
        for sock in self.clients.values():
            sock.close()


    def __exit_gracefully(self, *args):
        """
        Handles SIGTERM

        By setting self._stop_server to False, the server will continue with the iteration
        it was working, but it will be his last one before stopping gracefully. 
        """
        self._server_socket.shutdown(socket.SHUT_RDWR)
        self._stop_server = True
