// Package img wraps govips to perform image and buffer operations
package img

import (
	"github.com/davidbyttow/govips/v2/vips"
)

func Extract(target *vips.ImageRef, x, y, w, h int) (*vips.ImageRef, error) {
	copy, copyErr := target.Copy()
	if copyErr != nil {
		return nil, copyErr
	}
	extractErr := copy.ExtractArea(x, y, w, h)

	if extractErr != nil {
		return nil, extractErr
	}

	return copy, nil
}
