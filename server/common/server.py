import socket
import logging
import signal
from multiprocessing import Process, Pipe, Semaphore
from common import communication
from common.utils import store_bets, load_bets, has_won, handle_client

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

    def accept_clients(self):
        """
        For each client accepted, a new process is created.
        Pipes are used to communicate with them
        rec_EOF and send_EOF pipe are for the server to know when each of clients finished sending their data.
        rec_winner and send_winner are for the "father" process to communicate the others so they can send it to the client.
        """
        clients_accepted = 0
        sem = Semaphore()
        rec_EOF, send_EOF = Pipe(False)
        while not self._stop_server and clients_accepted != 5:
            try:
                agency, client_sock = self.__accept_new_connection()
                rec_winner, send_winner = Pipe(False)
                p = Process(target=handle_client, args=(agency,client_sock, send_EOF, rec_winner, sem,))
                p.start()
            except OSError as e:
                self._server_socket.close()
                rec_EOF.close()
                send_EOF.close()
                return None, None, Exception("SIGTERM RECEIVED")

            self.clients[int(agency)] = (send_winner, p, rec_winner, client_sock)
            clients_accepted += 1

        return rec_EOF, send_EOF, None

    def run(self):
        """
        Responsible for the whole logic of the server.
        """
        rec_EOF, send_EOF, err = self.accept_clients()
        if err is not None:
            logging.error(f'action: accept_clients | result: fail | error: {err}')
            return
        
        # Wait till all the clients finished sending the bets
        EOF_received = 0
        while EOF_received != len(self.clients):
            rec_EOF.recv()
            EOF_received += 1
        rec_EOF.close()
        send_EOF.close()

        # Look for winners and send them to their corresponding process
        self.start_lottery()
        logging.info('action: sorteo | result: success')

        self.__close_clients()

    def start_lottery(self):
        """
        For each winner found, a message with the document is sent to the corresponding process that handles the agency.
        """
        for bet in load_bets():
            if self._stop_server:
                return
            elif has_won(bet):
                send_winner = self.clients[bet.agency][0]
                send_winner.send(bet.document)      

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
        
    def __close_clients(self):
        """
        Close the las resources. The rec_winner and the client_sock
        are close inside the other processes.
        """
        for (send_winner, process, _, _) in self.clients.values():
            send_winner.send("EOF")
            send_winner.close()
            process.join()
            process.close()


    def __exit_gracefully(self, *args):
        """
        Handles SIGTERM
        """
        for (send_winner, p, rec_winner, client_sock) in self.clients.values():
            # Close pipes
            send_winner.close()
            rec_winner.close()
            # Close socket
            client_sock.close()
            # End process
            p.close()

        self._server_socket.shutdown(socket.SHUT_RDWR)
        self._stop_server = True
