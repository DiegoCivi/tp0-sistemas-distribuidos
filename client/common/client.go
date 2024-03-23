package common

import (
	"net"
	"os"
	"time"
	"io"

	log "github.com/sirupsen/logrus"
)

const (
	BATCH_SIZE = 7000
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

	reader, err := getReader(c.config.ID)
	if err != nil {
		log.Errorf("action: open_file | result: fail | client_id: %v | error: %v", c.config.ID, err)
		return
	}

	batch := []byte("")
	for {
		line, isPrefix, err := reader.ReadLine() // TODO: Use the isPrefix
		if err != nil {
			if err != io.EOF { // Handle any errors other than EOF 
				log.Errorf("action: read_line | result: fail | client_id: %v | error: %v", c.config.ID, err)
				c.conn.Close()
				return
			} else if len(batch) > 0 { // If an EOF was received, but theres still bytes on the batch
				if sendBatch(c.conn, batch, c.config.ID) != nil {
					log.Errorf("action: send_last_batch | result: fail | client_id: %v | error: %v", c.config.ID, err)
					c.conn.Close()
					return
				}
			}
			break
		} else if isPrefix { // If isPrefix is set, the line didnt enter so we have to read again
			batch = append(batch, line...)
			continue
		}

		// If adding the new line to the batch exceeds its maximum size,
		// the batch is sent and emptied.
		if len(line) + len(batch) > BATCH_SIZE {
			if sendBatch(c.conn, batch, c.config.ID) != nil {
				log.Errorf("action: send_batch | result: fail | client_id: %v | error: %v", c.config.ID, err)
				c.conn.Close()
				return
			}
			batch = []byte("")
		}

		// The read line is appended to the match, with a separator to diferentiate lines
		// The last byte will represent a '/'
		batch = append(batch, line...)
		batch = append(batch, []byte(SEPARATOR)...)
	}
	
	// Send the message with the END-FLAG set to true
	err = closeSocket(c.conn)
	if err != nil {
		log.Errorf("action: close_socket | result: fail | client_id: %v | error: %v", c.config.ID, err)
		return
	}

	log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
}
