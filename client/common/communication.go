package common

import (
	"errors"
	"net"
	"strconv"

	//log "github.com/sirupsen/logrus"
)

// The header is HEADER_LENGTH long, but MSG_SIZE_LENGTH bytes are for the part of the header that
// tell how long in bytes is the message. One byte will occupy the end_flag in the header.
const HEADER_LENGTH = 5
const MSG_SIZE_LENGTH = 4


// Called by writeSocket(). It returns the protocols header for a message.
func getHeader(msg []byte, end_flag string) []byte {
	header := strconv.Itoa(len(msg))
	msg_len_bytes := len(header)
	for i := 0; i < MSG_SIZE_LENGTH - msg_len_bytes; i++ {
		header = "0" + header
	}
	header += end_flag
	return []byte(header)
}

// Writes the message into the received socket
func writeSocket(conn net.Conn, msg []byte) error {
	// Add header
	header := getHeader(msg, "0")
	complete_msg := append(header, msg...)

	//log.Infof("Se va a mandar un batch")

	err := handleShortWrite(conn, complete_msg)
	if err != nil {
		return err
	}
	return nil
}

// Called by writeSocket(). It makes sure that if a short-write happens,
// the rest of the message is also sent.
func handleShortWrite(conn net.Conn, msg []byte) error {
	// Send serialized message to server, handling short read
	bytes_to_write := len(msg)
	bytes_wrote := 0
	for bytes_wrote < bytes_to_write {
		nbytes, err := conn.Write(msg[bytes_wrote:])
		if err != nil {
			return err
		}
		bytes_wrote += nbytes
	}

	return nil
}

// Reads from the received socket.
// It returns the message received or the error.
func readSocket(conn net.Conn) (string, error) {
	// Read header
	header, err := handleShortRead(conn, HEADER_LENGTH)
	if err != nil {
		return header, err
	}

	// Read message
	msg_len, _ := strconv.Atoi(header[:len(header) - 1])
	end_flag := string(header[len(header) - 1])
	if end_flag == "1" {
		return "", errors.New("EOF")
	}
	msg, err := handleShortRead(conn, msg_len)
	
	return msg, err
}

// Called by writeSocket(). It makes sure that if a short-read happens,
// the rest of the message is also read.
func handleShortRead(conn net.Conn, bytes_to_read int) (string, error) {
	bytes_read := 0
	msg := ""
	for bytes_read < bytes_to_read {
		buf := make([]byte, bytes_to_read - bytes_read)
		nbytes, err := conn.Read(buf)
		if err != nil {
			return "", err
		}
		msg += string(buf)
		bytes_read += nbytes
	}
	return msg, nil
}

func sendEOF(conn net.Conn) error {
	// Send the header with the end flag set on 1
	header := getHeader([]byte(""), "1")
	err := handleShortWrite(conn, header)
	return err
}
