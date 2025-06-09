package cards

import (
	"context"
	"fmt"

	"github.com/agladfield/postcart/pkg/shared/tools/img"
	"github.com/davidbyttow/govips/v2/vips"
)

const cardsBackErrFmtStr = "cards create back err: %w"

func createBack(ctx context.Context, email *EmailParams) (*sideOutput, error) {
	ascii := createBackASCII(email.Artwork, email.Border)
	// fmt.Println(ascii)
	// panic("")
	// get the art image early (large-ish network call so it could fail and it's better to fail early)
	artImage, artErr := getArtwork(ctx, email.Artwork, email.ArtStyle, email.Attachment)
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

	return &sideOutput{
		image: cropped,
		ascii: ascii,
	}, nil
}
