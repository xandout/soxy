package proxy

import (
	"net"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/xandout/soxy/wsconnadapter"
)

func chanFromConn(conn net.Conn) chan []byte {
	c := make(chan []byte)

	go func() {
		b := make([]byte, 1024)

		for {
			n, err := conn.Read(b)
			if n > 0 {
				res := make([]byte, n)
				// Copy the buffer so it doesn't get changed while read by the recipient.
				copy(res, b[:n])
				c <- res
			}
			if err != nil {
				//c <- nil //fix more 0x00 data bug
				break
			}
		}
	}()

	return c
}

// Copy accepts a websocket connection and TCP connection and copies data between them
func Copy(gwsConn *websocket.Conn, tcpConn net.Conn) {
	wsConn := wsconnadapter.New(gwsConn)
	wsChan := chanFromConn(wsConn)
	tcpChan := chanFromConn(tcpConn)

	defer wsConn.Close()
	defer tcpConn.Close()
	for {
		select {
		case wsData := <-wsChan:
			if wsData == nil {
				log.Infof("TCP connection closed: D: %v, S: %v", tcpConn.LocalAddr(), wsConn.RemoteAddr())
				return
			} else {
				tcpConn.Write(wsData)
			}
		case tcpData := <-tcpChan:
			if tcpData == nil {
				log.Infof("TCP connection closed: D: %v, S: %v", tcpConn.LocalAddr(), wsConn.LocalAddr())
				return
			} else {
				wsConn.Write(tcpData)
			}
		}
	}

}
