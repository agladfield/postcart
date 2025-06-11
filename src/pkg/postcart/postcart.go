// Package postcart exposes the main program loop and functionality
package postcart

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/agladfield/postcart/pkg/cards"
	"github.com/agladfield/postcart/pkg/jdb"
	"github.com/agladfield/postcart/pkg/postmark"
	"github.com/agladfield/postcart/pkg/server"
	"github.com/agladfield/postcart/pkg/shared/env"
)

func Close() {
	jdb.Close()
	cards.Close()
	server.Close()
}

const postcartProgramErrFmtStr = "postcart program err: %w"

func Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var wg sync.WaitGroup

	// ensure valid environment for operation
	envErr := env.Configure()
	if envErr != nil {
		return fmt.Errorf(postcartProgramErrFmtStr, envErr)
	}

	defer Close()
	jdbErr := jdb.Load()
	if jdbErr != nil {
		return jdbErr
	}

	// configure postmark
	postmarkErr := postmark.Configure()
	if postmarkErr != nil {
		return fmt.Errorf(postcartProgramErrFmtStr, postmarkErr)
	}

	cardsErr := cards.Prepare(ctx, &wg)
	if cardsErr != nil {
		return fmt.Errorf(postcartProgramErrFmtStr, cardsErr)
	}

	serverErr := server.Prepare()
	if serverErr != nil {
		return fmt.Errorf(postcartProgramErrFmtStr, serverErr)
	}

	// splash happens
	postcartSplash()

	server.Listen()

	// Listen for sig interrupt or sig termination
	exitChannel := make(chan os.Signal, 1)
	signal.Notify(exitChannel, syscall.SIGINT, syscall.SIGTERM)

	<-exitChannel

	return nil
}

// © Arthur Gladfield
