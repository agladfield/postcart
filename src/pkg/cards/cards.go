// Package cards exposes code for generating postcard images and ascii text
package cards

import (
	"context"
	"embed"
	"fmt"
	"runtime"
	"sync"

	"github.com/agladfield/postcart/pkg/shared/env"
	"github.com/davidbyttow/govips/v2/vips"
)

// list embedded directories here
// countries
// postcards
// stamps

//go:embed res/artwork
var placeholderArtwork embed.FS

//go:embed res/countries/circle
var circleFlags embed.FS

//go:embed res/countries/square
var squareFlags embed.FS

//go:embed res/fonts
var fontResources embed.FS

//go:embed res/postcards
var postcardResources embed.FS

//go:embed res/stamps
var stampResources embed.FS

const (
	cardWidth  = 1280
	cardHeight = 853
)

func Close() {
	// libvips close
	vips.Shutdown()
}

const (
	cardsErrFmtStr        = "cards err: %w"
	cardsPrepareErrFmtStr = "cards prepare err: %w"
)

func Prepare(ctx context.Context, wg *sync.WaitGroup) error {
	if env.InstallFonts() {
		fontsErr := installFonts()
		if fontsErr != nil {
			return fmt.Errorf(cardsPrepareErrFmtStr, fontsErr)
		}
	}

	vips.LoggingSettings(nil, vips.LogLevelError) // don't log anything unless an error with vips occurs
	vips.Startup(&vips.Config{
		ConcurrencyLevel: runtime.NumCPU(),
	})

	// create queue
	createQueue(ctx, wg)
	// start backlogger
	// retry jobs that failed every 2 minutes unless failed over three
	// times

	return nil
}
