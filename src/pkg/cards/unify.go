package cards

import (
	"fmt"

	"github.com/agladfield/postcart/pkg/shared/tools/img"
	"github.com/davidbyttow/govips/v2/vips"
)

const cardsUnificationErrFmtStr = "cards unify err: %w"

type unifiedOutput struct {
	UnifiedImage *vips.ImageRef
	UnifiedText  string
}

func unify(bordered *borderedOutput) (*unifiedOutput, error) {
	width := bordered.front.image.Width()

	height := bordered.front.image.Height() + bordered.back.image.Height()

	empty, emptyErr := img.New(width, height, false)
	if emptyErr != nil {
		return nil, fmt.Errorf(cardsUnificationErrFmtStr, emptyErr)
	}
	_, bErr := empty.ToBytes()
	if bErr != nil {
		return nil, bErr
	}

	// // composite the two sides onto the one image
	compositeFrontErr := empty.Composite(bordered.front.image, vips.BlendModeOver, 0, 0)
	if compositeFrontErr != nil {
		return nil, fmt.Errorf(cardsUnificationErrFmtStr, compositeFrontErr)
	}
	compositeBackErr := empty.Composite(bordered.back.image, vips.BlendModeOver, 0, height/2)
	if compositeBackErr != nil {
		return nil, fmt.Errorf(cardsUnificationErrFmtStr, compositeBackErr)
	}

	_, _, bytesErr := empty.ExportPng(&vips.PngExportParams{})
	if bytesErr != nil {
		return nil, bytesErr
	}

	return &unifiedOutput{
		UnifiedImage: empty,
		UnifiedText:  bordered.back.ascii + "\n" + bordered.front.ascii,
	}, nil
}

// © Arthur Gladfield
