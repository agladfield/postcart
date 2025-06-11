package cards

import (
	"context"
	"fmt"
	"sync"

	"github.com/agladfield/postcart/pkg/jdb"
	"github.com/agladfield/postcart/pkg/postmark"
)

var (
	blockSenderCh chan string
)

func createBlockQueues(ctx context.Context, wg *sync.WaitGroup) {
	blockSenderCh = make(chan string, 100)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case sender, ok := <-blockSenderCh:
				if !ok {
					return
				}
				// block sender email
				res, err := postmark.CreateInboundTriggerRule(sender)
				if err != nil {
					fmt.Printf("failed to block sender %q with error: %v", sender, err)
				}
				jdb.BlockSender(sender, res.ID)
				jdb.RecordBlockedSender()

			case <-ctx.Done():
				return
			}
		}
	}()
}

func BlockSender(email string) {
	blockSenderCh <- email
}

// © Arthur Gladfield
