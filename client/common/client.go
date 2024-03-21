package common

import (
	"net"
	"os"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	ERR = -1
	NO_ERR = 0
)

// Contains the info about the clients bet
type Bet struct {
	Name		string
	Surname		string
	Id			string
	Birth		string
	Number		string
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
	config 		ClientConfig
	conn   		net.Conn
	bet			Bet	
}

// Creates a Bet from the env variables
func CreateBet() Bet {
	bet := Bet{
		Name: os.Getenv("NAME"),
		Surname: os.Getenv("SURNAME"),
		Id: os.Getenv("ID"),
		Birth: os.Getenv("BIRTH"),
		Number: os.Getenv("NUMBER"),
	}
	return bet
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig, bet Bet) *Client {
	client := &Client {
		config: config,
		bet: bet,
	}
	return client
}

// CreateClientSocket Initializes client socket. In case of
// failure, error is printed in stdout/stderr and exit 1
// is returned
func (c *Client) createClientSocket() error {
	conn, err := createSocket(c.config.ServerAddress)
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

func (c *Client) sendServer(msg string) int {
	// Send serialized message to server, handling short read
	bytes_wrote := 0
	bytes_to_write := len(msg)
	for bytes_wrote < bytes_to_write {
		nbytes, err := writeSocket(c.conn, msg)
		if err != nil {
			log.Errorf("action: send_message | result: fail | client_id: %v | error: %v",
            	c.config.ID,
				err,
			)
			return ERR
		}
		bytes_wrote += nbytes
	}

	return NO_ERR
}

func (c *Client) readServer(bytes_to_read int) (string, error) {
	bytes_read := 0
	msg := ""
	for bytes_read < bytes_to_read {
		buf := make([]byte, bytes_to_read - bytes_read)
		nbytes, err := c.conn.Read(buf)
		if err != nil {
			return "", err
		}
		msg += string(buf)
		bytes_read += nbytes
	}
	return msg, nil
}


// StartClientLoop Send messages to the client until some time threshold is met
func (c *Client) StartClientLoop() {
	// Create the connection the server in every loop iteration. Send an
	// Skip the rest if the socket was not created 
	err := c.createClientSocket()
	if err != nil {
		log.Errorf("action: create_socket | result: fail | client_id: %v | error: %v",
            c.config.ID,
			err,
		)
		return
	} 
	
	// Send Bet to server
	msg := c.serialize()
	if c.sendServer(msg) == ERR {
		return
	}

	
	// Read header
	header, err := c.readServer(HEADER_LENGTH)
	if err != nil {
		log.Errorf("action: receive_message | result: fail | client_id: %v | error: %v",
            c.config.ID,
			err,
		) 
		return
	}

	// Read message
	msg_len, _ := strconv.Atoi(header)
	bet_msg, err := c.readServer(msg_len)
	if err != nil {
		log.Errorf("action: receive_message | result: fail | client_id: %v | error: %v",
            c.config.ID,
			err,
		) 
		return
	}
	
	log.Infof("action: receive_message | result: success | client_id: %v | msg: %v",
		c.config.ID,
		bet_msg,
	)

	c.conn.Close()

	log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
}
