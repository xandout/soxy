package server

import (
	"crypto/tls"
	"net"
	"net/http"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/xandout/soxy/proxy"
)

// Start starts the http server
func Start(c *cli.Context) error {
	port := c.String("port")
	http.Handle("/", socketHandler{apiKey: c.String("api-key")})
	err := http.ListenAndServe(port, nil)
	log.Errorf("HTTP SERVER: %v", err.Error())
	return err

}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type socketHandler struct {
	apiKey string
}

func (h socketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.apiKey != "" {
		if h.apiKey != r.Header.Get("X-Api-Key") {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Forbidden"))
			log.Errorf("HTTP SERVER: %v", "Invalid API Key")
			return
		}
	}
	q := r.URL.Query()
	var useTLS bool
	if q.Get("useTLS") != "" {
		useTLS = true
	}
	remote := q.Get("remote")
	if remote == "" {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("remote not set"))
		log.Errorf("HTTP SERVER: %v", "remote not set")
		return
	}
	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		log.Errorf("HTTP SERVER, WS Connection Upgrade: %v", err.Error())
		return
	}
	var remoteTCPConn net.Conn
	if useTLS {
		remoteTCPConn, err = tls.Dial("tcp", remote, &tls.Config{
			InsecureSkipVerify: true,
		})
	} else {
		remoteTCPConn, err = net.Dial("tcp", remote)
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		log.Errorf("HTTP SERVER, TCP Write: %v", err.Error())
		return
	}
	log.Infof("Proxying traffic to %v on behalf of %v", remoteTCPConn.RemoteAddr(), wsConn.RemoteAddr())
	go proxy.Copy(wsConn, remoteTCPConn)
}
