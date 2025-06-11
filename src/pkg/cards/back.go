package cards

import (
	"context"
	"fmt"

	"github.com/agladfield/postcart/pkg/shared/tools/img"
	"github.com/davidbyttow/govips/v2/vips"
)

const cardsBackErrFmtStr = "cards create back err: %w"

func createBack(ctx context.Context, params *Params) (*sideOutput, error) {
	ascii := createBackASCII(params.Artwork, params.Border)
	// get the art image early (large-ish network call so it could fail and it's better to fail early)
	artImage, artErr := getArtwork(ctx, params.Artwork, params.Style, params.Attachment)
	if artErr != nil {
		return nil, fmt.Errorf(cardsBackErrFmtStr, artErr)
	}
	defer artImage.Close()

	// resize if too small or too big width wise
	if artImage.Width() != cardWidth {
		resizeFactor := float64(cardWidth) / float64(artImage.Width())
		resizeErr := artImage.Resize(resizeFactor, vips.KernelAuto)
		if resizeErr != nil {
			return nil, fmt.Errorf(cardsBackErrFmtStr, resizeErr)
		}
	}

	// mask the back
	cropWidth := (artImage.Width() - cardWidth) / 2
	cropHeight := (artImage.Height() - cardHeight) / 2

	cropped, cropErr := img.Extract(artImage, cropWidth, cropHeight, cardWidth, cardHeight)
	if cropErr != nil {
		return nil, fmt.Errorf(cardsBackErrFmtStr, cropErr)
	}

	backBuff, _, backBuffErr := cropped.ExportPng(&vips.PngExportParams{Quality: 90})
	if backBuffErr != nil {
		return nil, backBuffErr
	}
	defer cropped.Close()
	back, backErr := img.LoadFromBuffer(backBuff)
	if backErr != nil {
		return nil, backErr
	}

	return &sideOutput{
		image: back,
		ascii: ascii,
	}, nil
}

// © Arthur Gladfield
