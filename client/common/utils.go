package common

import (
	"bufio"
	"io"
	"net"
	"os"

	log "github.com/sirupsen/logrus"
)

func sendBatch(conn net.Conn, batch []byte, id string) error {
	// Send batch to the server
	err := writeSocket(conn, batch)
	if err != nil {
		log.Errorf("action: send_message | result: fail | client_id: %v | error: %v",
			id,
			err,
		)

		return err
	}

	// Read ack from server
	_, err = readSocket(conn)
	if err != nil {
		log.Errorf("action: receive_message | result: fail | client_id: %v | error: %v",
			id,
			err,
		)
		return err
	}
	return nil
}

func getReader(id string) (*bufio.Reader, *os.File, error) {
	file_path := "./data/agency-" + id + ".csv"

	file, err := os.Open(file_path)
	if err != nil {
		return nil, nil, err
	}

	reader := bufio.NewReader(file)

	return reader, file, nil
}

func handleFileErrors(err error, conn net.Conn, id string, batch []byte) bool {
	if err != io.EOF { // Handle any errors other than EOF 
		log.Errorf("action: read_line | result: fail | client_id: %v | error: %v", id, err)
		return true 
	} else if len(batch) > 0 { // If an EOF was received, but theres still bytes on the batch
		if sendBatch(conn, batch, id) != nil {
			log.Errorf("action: send_last_batch | result: fail | client_id: %v | error: %v", id, err)
			return true
		}
	}
	return false
}
