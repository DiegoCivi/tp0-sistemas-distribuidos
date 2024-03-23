package common

import (
	"bufio"
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

func getReader(id string) (*bufio.Reader, error) {
	file_path := "./data/agency-" + id + ".csv"

	file, err := os.Open(file_path)
	if err != nil {
		return nil, err
	}

	reader := bufio.NewReader(file)

	return reader, nil
}
