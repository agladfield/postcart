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
	closeCache()
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
	cacheErr := createCaches()
	if cacheErr != nil {
		return cacheErr
	}

	vips.LoggingSettings(nil, vips.LogLevelWarning) // don't log anything unless an error with vips occurs
	vips.Startup(&vips.Config{
		ConcurrencyLevel: runtime.NumCPU(),
	})

	templatesErr := checkTemplatesAreAvailable()
	if templatesErr != nil {
		return templatesErr
	}

	// create queues
	createBlockQueues(ctx, wg)
	createQueue(ctx, wg)

	return nil
}

// © Arthur Gladfield
