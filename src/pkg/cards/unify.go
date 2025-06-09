package cards

import (
	"github.com/agladfield/postcart/pkg/shared/tools/img"
	"github.com/davidbyttow/govips/v2/vips"
)

type unifiedOutput struct {
	UnifiedImage *vips.ImageRef
	UnifiedText  string
}

func unify(bordered *borderedOutput) (*unifiedOutput, error) {
	width := bordered.front.image.Width()
	// widthDifference := 0
	if bordered.back.image.Width() > width {
		// widthDifference = back.backImage.Width() - width
		width += bordered.back.image.Width() - width
	} else if bordered.back.image.Width() < width {
		// widthDifference = width - back.backImage.Width()
	}

	height := bordered.front.image.Height() + bordered.back.image.Height()

	empty, emptyErr := img.New(width, height, true)
	if emptyErr != nil {
		return nil, emptyErr
	}

	// composite the two on
	compositeFrontErr := empty.Composite(bordered.front.image, vips.BlendModeOver, 0, 0)
	if compositeFrontErr != nil {
		return nil, compositeFrontErr
	}
	compositeBackErr := empty.Composite(bordered.back.image, vips.BlendModeOver, 0, height/2)
	if compositeBackErr != nil {
		return nil, compositeBackErr
	}

	return &unifiedOutput{
		UnifiedImage: empty,
		UnifiedText:  bordered.back.ascii + "\n" + bordered.front.ascii,
	}, nil
}
