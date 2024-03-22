package common

import (
	"net"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
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
		log.Errorf("action: create_socket | result: fail | client_id: %v | error: %v",
            c.config.ID,
			err,
		)
		return
	} 
	
	// Send Bet to the server
	msg := c.serialize()
	err = writeSocket(c.conn, msg)
	if err != nil {
		log.Errorf("action: send_message | result: fail | client_id: %v | error: %v",
            c.config.ID,
			err,
		)
		c.conn.Close()
		return
	} 

	// Read Bet ack from server
	bet_msg, err := readSocket(c.conn)
	if err != nil {
		log.Errorf("action: receive_message | result: fail | client_id: %v | error: %v",
            c.config.ID,
			err,
		)
		c.conn.Close()
		return
	}
	
	log.Infof("action: receive_message | result: success | client_id: %v | msg: %s",
		c.config.ID,
		bet_msg,
	)

	c.conn.Close()

	log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
}
