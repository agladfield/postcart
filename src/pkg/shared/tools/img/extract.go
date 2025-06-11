package img

import (
	"fmt"

	"github.com/davidbyttow/govips/v2/vips"
)

const imgExtractErrFmtStr = "img extract err: %w"

func Extract(target *vips.ImageRef, x, y, w, h int) (*vips.ImageRef, error) {
	copy, copyErr := target.Copy()
	if copyErr != nil {
		return nil, fmt.Errorf(imgExtractErrFmtStr, copyErr)
	}
	extractErr := copy.ExtractArea(x, y, w, h)

	if extractErr != nil {
		return nil, fmt.Errorf(imgExtractErrFmtStr, extractErr)
	}

	return copy, nil
}

// © Arthur Gladfield
