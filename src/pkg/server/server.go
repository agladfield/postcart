// Package server contains the HTTP server code
package server

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/agladfield/postcart/pkg/server/hooks"
)

var (
	// port
	server *http.Server
)

func checkPortAvailability(port string) error {
	host := "127.0.0.1"
	timeout := time.Second
	conn, err := net.DialTimeout("tcp", host+":"+port, timeout)
	if err != nil {
		return nil
	} else {
		conn.Close()
		return fmt.Errorf("port %s is in use", port)
	}
}

const (
	httpServerErrFmtStr        = "http server err: %w"
	httpServerPrepareErrFmtStr = "http server prepare err: %w"
)

func Prepare() error {
	portErr := checkPortAvailability("8080")
	if portErr != nil {
		return fmt.Errorf(httpServerPrepareErrFmtStr, portErr)
	}

	hooksConfigErr := hooks.Configure()
	if hooksConfigErr != nil {
		return fmt.Errorf(httpServerPrepareErrFmtStr, hooksConfigErr)
	}

	return nil
}

func Close() {
	//
}

func Listen() {
	// add hook routes
	hooks.Routes()

	// var serverErr error

	server = &http.Server{
		Addr: ":8080",
	}

	// Start the server in a goroutine
	go func() {
		if listenErr := server.ListenAndServe(); listenErr != http.ErrServerClosed {
			log.Fatalf("http server ListenAndServe(): %v", listenErr)
		}
	}()

	log.Printf("listening for postmark webhook requests on port %s", "8080")
}
