package img

import (
	"errors"

	"github.com/davidbyttow/govips/v2/vips"
)

func Mask(source, mask *vips.ImageRef) (*vips.ImageRef, error) {
	if source.Width() != mask.Width() || source.Height() != mask.Height() {
		// error size difference
		return nil, errors.New("mask must be same size as source")
	}
	maskCopy, maskCopyErr := mask.Copy()
	if maskCopyErr != nil {
		return nil, maskCopyErr
	}

	maskErr := maskCopy.Composite(source, vips.BlendModeIn, 0, 0)
	if maskErr != nil {
		return nil, maskErr
	}

	return maskCopy, nil
}
