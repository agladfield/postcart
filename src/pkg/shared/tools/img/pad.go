package img

import (
	"fmt"

	"github.com/davidbyttow/govips/v2/vips"
)

const imgPadErrFmtStr = "img pad err: %w"

func Pad(src *vips.ImageRef, padding [2]int, vertical bool) (*vips.ImageRef, error) {
	srcCopy, srcCopyErr := src.Copy()
	if srcCopyErr != nil {
		return nil, fmt.Errorf(imgPadErrFmtStr, srcCopyErr)
	}
	defer srcCopy.Close()

	blank, blankErr := New(srcCopy.Width(), srcCopy.Height(), true)
	if blankErr != nil {
		return nil, fmt.Errorf(imgPadErrFmtStr, blankErr)
	}

	difference := padding[0] + padding[1]

	var resizeFactor float64

	if vertical {
		resizeFactor = float64(srcCopy.Height()-difference) / float64(srcCopy.Height())
	} else {
		resizeFactor = float64(srcCopy.Width()-difference) / float64(srcCopy.Width())
	}

	resizeErr := srcCopy.Resize(resizeFactor, vips.KernelAuto) // changed from vips.KernelNearest to reduce sharpness
	if resizeErr != nil {
		return nil, fmt.Errorf(imgPadErrFmtStr, resizeErr)
	}

	var (
		compX int
		compY int
	)

	if vertical {
		compX = (blank.Width() - srcCopy.Width()) / 2
		compY = padding[0]
	} else {
		compX = padding[0]
		compY = (blank.Height() - srcCopy.Height()) / 2
	}

	compositeErr := blank.Composite(srcCopy, vips.BlendModeOver, compX, compY)
	if compositeErr != nil {
		return nil, fmt.Errorf(imgPadErrFmtStr, compositeErr)
	}

	return blank, nil
}

// © Arthur Gladfield
