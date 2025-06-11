package img

import (
	"errors"
	"fmt"

	"github.com/davidbyttow/govips/v2/vips"
)

const imgMaskErrFmtStr = "img mask err: %w"

func Mask(source, mask *vips.ImageRef) (*vips.ImageRef, error) {
	if source.Width() != mask.Width() || source.Height() != mask.Height() {
		// error size difference
		return nil, fmt.Errorf(imgMaskErrFmtStr, errors.New("mask must be same size as source"))
	}
	maskCopy, maskCopyErr := mask.Copy()
	if maskCopyErr != nil {
		return nil, fmt.Errorf(imgMaskErrFmtStr, maskCopyErr)
	}

	maskErr := maskCopy.Composite(source, vips.BlendModeIn, 0, 0)
	if maskErr != nil {
		return nil, fmt.Errorf(imgMaskErrFmtStr, maskErr)
	}

	return maskCopy, nil
}

// © Arthur Gladfield
