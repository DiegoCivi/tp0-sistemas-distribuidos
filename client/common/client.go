package common

import (
	"fmt"
	"net"
	"os"
	"reflect"
	"strconv"
	"strings"
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

func (c *Client) serialize() string {
	msg := c.config.ID + "/"
	v := reflect.ValueOf(c.bet)

	// Iterate over Bet fields and add them to the message
	for i := 0; i < v.NumField(); i++ {

		val := v.Field(i).Interface()

		// t :=  fmt.Sprintf("%s", val)
		// log.Infof("[SERIALIZE] VALOR: %s", t)

		msg += fmt.Sprintf("%s/", val)
	}

	msg = strings.TrimSuffix(msg, "/")
	// log.Infof("[SERIALIZE] El mensaje serializado es: %s", msg)

	return msg
}


// StartClientLoop Send messages to the client until some time threshold is met
func (c *Client) StartClientLoop() {
	// Create the connection the server in every loop iteration. Send an
	// Skip the rest of the iteration if thee socket was not created 
	err := c.createClientSocket()
	if err != nil {
		log.Errorf("action: create_socket | result: fail | client_id: %v | error: %v",
            c.config.ID,
			err,
		)
		return
	} 
	
	msg := c.serialize()

	log.Infof("[START CLIENT LOOP] Se serializo el mensaje y quedo: %s", msg)

	//bytes_wrote := 0
	// TO-DO: Handle short write
	//bytes_wrote, err := communication.writeSocket(c.conn, msg)
	writeSocket(c.conn, msg)
	
	// Read header
	buf := make([]byte, HEADER_LENGTH)
	nbytes, err := c.conn.Read(buf)
	if nbytes != HEADER_LENGTH {
		log.Infof("action: receive_HEADER | result: SHORT-REEAD | client_id: %v", c.config.ID)
	} else if err != nil {
		log.Errorf("action: receive_HEADER | result: fail | client_id: %v | error: %v",
            c.config.ID,
			err,
		)
	}

	// Read message
	msg_len, _ := strconv.Atoi(string(buf))
	bytes_read := 0
	bet_msg := ""
	for bytes_read < msg_len {
		buf = make([]byte, msg_len - bytes_read)
		nbytes, err = c.conn.Read(buf)
		if err != nil {
			log.Errorf("action: receive_message | result: fail | client_id: %v | error: %v",
            c.config.ID,
			err,
			) 
			return
		}
		bet_msg += string(buf)
		bytes_read += nbytes
	}
	
	log.Infof("action: receive_message | result: success | client_id: %v | msg: %v",
	c.config.ID,
	bet_msg,
)

	c.conn.Close()

	log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
}
