package common

import (
	"bufio"
	//"errors"
	//"io/ioutil"
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
	bet_msg, err := readSocket(conn)
	if err != nil {
		log.Errorf("action: receive_message | result: fail | client_id: %v | error: %v",
			id,
			err,
		)
		return err
	}
	
	log.Infof("action: receive_message | result: success | client_id: %v | msg: %s",
		id,
		bet_msg,
	)
	return nil
}

func getReader(id string) (*bufio.Reader, error) {
	file_path := "./data/agency-" + id + ".csv"

	log.Infof("[CLIENT-%s] Voy a abrir: %s",
		id,
		file_path,
	)

	file, err := os.Open(file_path) // AGREGAR BIEN EL PATH
    if err != nil {
        return nil, err
    }
	
	reader := bufio.NewReader(file)

	return reader, nil
}

//func handleReadLine(reader *bufio.Reader, batch []byte) ([]byte, error) {
//	line, _, err := reader.ReadLine() // TODO: Use the isPrefix
//	if err != nil {
//		if err != errors.New("EOF") { // Handle any errors other than EOF 
//			log.Errorf("action: send_message | result: fail | client_id: %v | error: %v", c.config.ID, err)
//			return nil, err
//		} else if len(batch) > 0 { // If an EOF was received, but theres still bytes on the batch
//			//if sendBatch(c.conn, batch, c.config.ID) != nil {
//			//	log.Errorf("action: send_last_batch | result: fail | client_id: %v | error: %v", c.config.ID, err)
//			//	c.conn.Close()
//			//	return
//			//}
//		}
//	}
//	return line, nil
//}