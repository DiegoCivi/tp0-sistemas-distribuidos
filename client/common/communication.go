package common

import (
	"net"
	"fmt"
	"strconv"
	log "github.com/sirupsen/logrus"
)

const (
    READ_BUF_SIZE  = 1024
    WRITE_BUF_SIZE = 1024
	HEADER_LENGTH = 4
)

func createSocket(addr string) (net.Conn, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	} 

	//err = setBuffers(conn)
	//if err != nil {
	//	conn.Close()
	//	return nil, err
	//}

	return conn, nil
}

// func setBuffers(conn net.Conn) error {
// 	err := conn.SetReadBuffer(READ_BUF_SIZE)
// 	if err != nil {
// 		return err
// 	}
// 
// 	err = conn.SetWriteBuffer(WRITE_BUF_SIZE)
// 	if err != nil {
// 		return err
// 	}
// 
// 	return nil
// }

func writeSocket(conn net.Conn, msg string) (int, error) {
	// Add header
	header := getHeader(msg)
	complete_msg := header + msg

	log.Infof("[WRITE-SOCKET] El mensaje enviado es: %s", complete_msg)

	// Send the serialized msg to the server
	return fmt.Fprintf(conn, "%s", complete_msg)
}


//func readSocket(conn net.Conn) string {
//	buf := make([]byte, HEADER_LENGTH)
//	conn.Read(buf) // do smth with the bytes read
//	header := string(buf)
//	log.Infof("[READ-SOCKET] El header recibido es: %s", header)
//	i, _ := strconv.Atoi(header) // handle if this fails
//
//	buf = make([]byte, i)
//	conn.Read(buf)
//	msg := string(buf) 
//	return msg
//}

//func readSocket(conn net.Conn, buf []byte) string {
//	read_bytes, err := conn.Read(buf)
//	msg := string(buf) 
//	return msg
//}

func getHeader(msg string) string {
	msg_len := strconv.Itoa(len(msg))
	msg_len_bytes := len(msg_len)
	for i := 0; i < HEADER_LENGTH - msg_len_bytes; i++ {
		msg_len = "0" + msg_len
	}
	return msg_len
}
