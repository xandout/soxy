package client

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"
	"os"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/xandout/soxy/proxy"
)

// Start starts a soxy client
func Start(c *cli.Context) error {

	l, err := net.Listen("tcp", c.String("local"))
	if err != nil {
		log.Errorf("TCP LISTENER: %v", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	log.Info("Listening on " + c.String("local"))

	insecureSkipVerify := c.String("InsecureSkipVerify")
	caPath := c.String("caPath")
	clientCertPath := c.String("clientCertPath")
	clientKeyPath := c.String("clientKeyPath")
	tslConfig := &tls.Config{}

	if "yes" == insecureSkipVerify {
		tslConfig.InsecureSkipVerify = true
	} else {
		tslConfig.InsecureSkipVerify = false
	}

	if caPath != "" {
		pool := x509.NewCertPool()
		caCrt, err := ioutil.ReadFile(caPath)
		if err != nil {
			log.Fatal("read ca.crt file error:", err.Error())
		}
		pool.AppendCertsFromPEM(caCrt)
		tslConfig.RootCAs = pool
	}

	if clientCertPath != "" && clientKeyPath != "" {
		cliCrt, err := tls.LoadX509KeyPair(clientCertPath, clientKeyPath)
		if err != nil {
			fmt.Println("Loadx509keypair err:", err)
			return err
		}
		tslConfig.Certificates = []tls.Certificate{cliCrt}
	}

	websocket.DefaultDialer.TLSClientConfig = tslConfig

	for {
		// Listen for an incoming connection.
		tcpConn, err := l.Accept()
		if err != nil {
			log.Errorf("TCP ACCEPT: %v", err.Error())
		}
		fmtString := "%s/?remote=%s"
		fmted := fmt.Sprintf(fmtString, c.String("soxy-url"), c.String("remote"))

		clientWsConn, _, err := websocket.DefaultDialer.Dial(fmted, nil)
		if err != nil {
			log.Errorf("DIALER: %v", err.Error())
			//return err  never stop client process
		}
		// Handle connections in a new goroutine.
		log.Infof("Proxying traffic to %v via %v for %v", c.String("remote"), clientWsConn.RemoteAddr(), tcpConn.RemoteAddr())
		go proxy.Copy(clientWsConn, tcpConn)

	}

}
