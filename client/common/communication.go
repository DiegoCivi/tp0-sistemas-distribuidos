package communication

import {
	"net"
	"fmt"
	"reflect"
}

const (
    READ_BUF_SIZE  = 1024
    WRITE_BUF_SIZE = 1024
	NO_ERROR = nil
)

func (addr string) createSocket() (net.Conn, error) {
	conn, err := net.Dial("tcp", c.config.ServerAddress)
	if err != NO_ERROR {
		return NO_ERROR, err
	} 

	err := setBuffers(conn)
	if err != NO_ERROR {
		conn.Close()
		return nil, err
	}

	return conn, NO_ERROR
}

func setBuffers(conn net.Conn) error {
	err := conn.SetReadBuffer(READ_BUF_SIZE)
	if err != NO_ERROR {
		return err
	}

	err := conn.SetWriteBuffer(WRITE_BUF_SIZE)
	if err != NO_ERROR {
		return err
	}

	return NO_ERROR
}

func writeSocket(conn net.Conn, bet Bet) (int, error) {
	msg := serialize(bet)

	// Send the serialized msg to the server
	bytes_read, err = fmt.Fprintf(conn, msg)

	return bytes_read, err
}

func serialize(bet Bet) string {
	msg := "-"
	v := reflect.ValueOf(bet)

	// Iterate over Bet fields and add them to the message
	for i := 0; i < v.NumField(); i++ {

		val := v.Field(i).Interface()

		msg += val + "-"
	}

	
}