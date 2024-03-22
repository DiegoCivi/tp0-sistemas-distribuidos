package common

import (
	"errors"
	"net"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	BATCH_SIZE = 4096
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
		log.Errorf("action: create_socket | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
		return
	}

	reader, err := getReader(c.config.ID)
	if err != nil {
		log.Errorf("action: open_file | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
		return
	}

	batch := []byte("")
	for {
		line, _, err := reader.ReadLine() // TODO: Use the isPrefix
		if err != nil {
			if err != errors.New("EOF") {
				log.Errorf("action: send_message | result: fail | client_id: %v | error: %v",
					c.config.ID,
					err,
				)
				c.conn.Close()
				return
			}
			break
		}

		if len(line)+len(batch) > BATCH_SIZE {
			if sendBatch(c.conn, batch, c.config.ID) != nil {
				//ver lo que hay que cerrar
				return
			}
			batch = []byte("")
		}

		batch = append(batch, line...)
	}

	c.conn.Close()

	log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
}
