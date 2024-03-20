package common

import (
	"bufio"
	"fmt"
	"net"
	"time"
	"os"
	"os/signal"
	"syscall"

	"./communication.go"

	log "github.com/sirupsen/logrus"
)

const EMPTY_ENV = ""

// Contains the info about the clients bet
type Bet struct {
	name		string
	surname		string
	id			string
	birth		string
	number		string
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
func createBet() Bet {
	//bet := Bet{}
	//bet_type := reflect.TypeOf(bet)
	//for i:= 0, i < bet_type.NumField(), i++ {
	//	env_var := bet_type.Field(i).Name
	//	val, ok := os.LookupEnv(env_var) 
	//	
	//	// Checks if the env var was set. If it was not or it has an empty value, nil is set.
	//	if val == EMPTY_ENV {
	//		
	//	}
	//}
	bet = common.Bet{
		name: os.Getenv("NAME"),
		surname: os.Getenv("SURNAME"),
		id: os.Getenv("ID"),
		birth: os.Getenv("BIRTH"),
		number: os.Getenv("NUMBER"),
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
	conn, err = communication.createSocket(c.config.ServerAddress)
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
	// autoincremental msgID to identify every message sent
	msgID := 1

	// Channel to receive SIGTERM signal
	signal_chan := make(chan os.Signal, 1)
	signal.Notify(signal_chan, syscall.SIGTERM)

loop:
	// Send messages if the loopLapse threshold has not been surpassed or a SIGTERM has not been received
	for timeout := time.After(c.config.LoopLapse); ; {
		select {
		case <-timeout:
	        log.Infof("action: timeout_detected | result: success | client_id: %v",
                c.config.ID,
            )
			break loop
		case <- signal_chan: // Check if a SIGTERM was received before starting the iteration
			log.Infof("action: SIGTERM_detected | result: success | client_id: %v",
                c.config.ID,
            )
			break loop
		default:
		}

		// Create the connection the server in every loop iteration. Send an
		// Skip the rest of the iteration if thee socket was not created  
		if c.createClientSocket() != nil {
			continue
		} 

		// TODO: Modify the send to avoid short-write. Get the msg from somewhere and the number of bytes
		communication.writeSocket(c.conn, c.bet)
		//fmt.Fprintf(
		//	c.conn,
		//	"[CLIENT %v] Message N°%v\n",
		//	c.config.ID,
		//	msgID,
		//)
		msg, err := bufio.NewReader(c.conn).ReadString('\n')
		msgID++
		c.conn.Close()

		if err != nil {
			log.Errorf("action: receive_message | result: fail | client_id: %v | error: %v",
                c.config.ID,
				err,
			)
			return
		}
		log.Infof("action: receive_message | result: success | client_id: %v | msg: %v",
            c.config.ID,
            msg,
        )

		// Wait a time between sending one message and the next one
		time.Sleep(c.config.LoopPeriod)
	}

	log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
}
