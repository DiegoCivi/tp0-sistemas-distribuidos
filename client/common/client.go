package common

import (
	"bufio"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	BATCH_MAX_SIZE = 8000
	BETS_IN_BATCH = 150
	SEPARATOR = "/"
)

// Contains the info about the clients bet
type Bet struct {
	Name    string
	Surname string
	Id      string
	Birth   string
	Number  string
}

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID            string
	ServerAddress string
	LoopLapse     time.Duration
	LoopPeriod    time.Duration
}

// Client Entity that encapsulates how
type Client struct {
	config ClientConfig
	conn   net.Conn
	bet    Bet
}

// Creates a Bet from the env variables
func CreateBet() Bet {
	bet := Bet{
		Name:    os.Getenv("NAME"),
		Surname: os.Getenv("SURNAME"),
		Id:      os.Getenv("ID"),
		Birth:   os.Getenv("BIRTH"),
		Number:  os.Getenv("NUMBER"),
	}
	return bet
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig, bet Bet) *Client {
	client := &Client{
		config: config,
		bet:    bet,
	}
	return client
}

// CreateClientSocket Initializes client socket. In case of
// failure, error is printed in stdout/stderr and exit 1
// is returned
func (c *Client) createClientSocket() error {
	conn, err := net.Dial("tcp", c.config.ServerAddress)
	if err != nil {
		log.Fatalf(
			"action: connect | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
		return err
	}
	c.conn = conn
	return nil
}

func sendBets(reader *bufio.Reader, conn net.Conn, id string, file *os.File) error {
	log.Infof("Entro a sendBets")
	// Channel to receive SIGTERM signal
	signal_chan := make(chan os.Signal, 1)
	signal.Notify(signal_chan, syscall.SIGTERM)

	batch := []byte("")
	bets_in_msg := 0
loop:
	for {
		select {
		case <-signal_chan:
			log.Errorf("action: sigterm_received | result: success | client_id: %v", id)
			conn.Close()
			file.Close()
			break loop
		default:
		}

		line, isPrefix, err := reader.ReadLine()
		if err != nil {
			finish_conn := handleFileErrors(err, conn, id, batch)
			if finish_conn {
				conn.Close()
				file.Close()
				return err
			}
			break
		} else if isPrefix { // If isPrefix is set, the line didnt enter so we have to read again
			batch = append(batch, line...)
			continue
		}

		// If adding the new line to the batch exceeds its maximum size, 
		// or there are BETS_IN_BATCH bets in the message are the batch is sent and emptied.
		if len(line) + len(batch) > BATCH_MAX_SIZE || bets_in_msg == BETS_IN_BATCH {
			if sendBatch(conn, batch, id) != nil {
				log.Errorf("action: send_batch | result: fail | client_id: %v | error: %v", id, err)
				conn.Close()
				file.Close()
				return err
			}
			batch = []byte("")
			bets_in_msg = 0
		}

		// The read line is appended to the match, with a separator to diferentiate lines
		// The last byte will represent a '/'
		batch = append(batch, line...)
		batch = append(batch, []byte(SEPARATOR)...)
		bets_in_msg += 1
	}
	log.Infof("TERMINO SENDBETS")
	return nil
}

// Client reads the socket, receiving the different winners until the server tells to stop reading.
func readLotteryWinners(conn net.Conn, id string) error {
	winners := 0
loop:
	for {
		_, err := readSocket(conn)
		if err != nil {
			if err.Error() == "EOF" {
				break loop
			}
			log.Errorf("action: receive_winner | result: fail | client_id: %v | error: %v", id, err)
			return err
		}
		winners += 1
	}

	log.Infof("action: consulta_ganadores | result: success | cant_ganadores: %d", winners)

	return nil
}

// StartClientLoop Send messages to the client until some time threshold is met
func (c *Client) StartClientLoop() {
	// Create the connection the server in every loop iteration. Send an
	// Skip the rest if the socket was not created
	err := c.createClientSocket()
	if err != nil {
		log.Errorf("action: create_socket | result: fail | client_id: %v | error: %v", c.config.ID, err)
		return
	}

	// Send the client id
	err = writeSocket(c.conn, []byte(c.config.ID))
	if err != nil {
		log.Errorf("action: send_ID | result: fail | client_id: %v | error: %v", c.config.ID, err)
		return
	}

	reader, file, err := getReader(c.config.ID)
	if err != nil {
		log.Errorf("action: open_file | result: fail | client_id: %v | error: %v", c.config.ID, err)
		return
	}

	// Send all the bets from the file to the server
	err = sendBets(reader, c.conn, c.config.ID, file)
	if err != nil {
		return
	}

	// Send the message with the END-FLAG set to true
	err = sendEOF(c.conn)
	if err != nil {
		log.Errorf("action: close_socket | result: fail | client_id: %v | error: %v", c.config.ID, err)
		return
	}

	// Read loterry winners
	err = readLotteryWinners(c.conn, c.config.ID)
	if err != nil {
		log.Errorf("action: read_lottery_winners | result: fail | client_id: %v | error: %v", c.config.ID, err)
		return
	}

	c.conn.Close()
	file.Close()

	log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
}
