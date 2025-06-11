package img

import (
	"fmt"

	"github.com/davidbyttow/govips/v2/vips"
)

const imgFlattenErrFmtStr = "img flatten err: %w"

func Flatten(backdrop *vips.ImageRef, layers ...*vips.ImageRef) (*vips.ImageRef, error) {
	w := backdrop.Width()
	h := backdrop.Height()
	sceneBlank, sceneBlankErr := New(w, h, true)
	if sceneBlankErr != nil {
		return nil, fmt.Errorf(imgFlattenErrFmtStr, sceneBlankErr)
	}

	compositeOntoNewErr := sceneBlank.Composite(backdrop, vips.BlendModeOver, 0, 0)
	if compositeOntoNewErr != nil {
		return nil, fmt.Errorf(imgFlattenErrFmtStr, compositeOntoNewErr)
	}

	for _, layer := range layers {
		layerCopy, copyErr := layer.Copy()
		if copyErr != nil {
			sceneBlank.Close()
			return nil, fmt.Errorf(imgFlattenErrFmtStr, copyErr)
		}
		defer layerCopy.Close()

		layerCompositeErr := sceneBlank.Composite(layerCopy, vips.BlendModeOver, 0, 0)
		if layerCompositeErr != nil {
			sceneBlank.Close()
			return nil, fmt.Errorf(imgFlattenErrFmtStr, layerCompositeErr)
		}
	}

	return sceneBlank, nil
}

func FlattenNoAlpha(backdrop *vips.ImageRef, layers ...*vips.ImageRef) (*vips.ImageRef, error) {
	bg, copyErr := backdrop.Copy()
	if copyErr != nil {
		return nil, fmt.Errorf(imgFlattenErrFmtStr, copyErr)
	}

	for _, layer := range layers {
		layerCopy, copyErr := layer.Copy()
		if copyErr != nil {
			bg.Close()
			return nil, fmt.Errorf(imgFlattenErrFmtStr, copyErr)
		}
		defer layerCopy.Close()

		layerCompositeErr := bg.Composite(layerCopy, vips.BlendModeOver, 0, 0)
		if layerCompositeErr != nil {
			bg.Close()
			return nil, fmt.Errorf(imgFlattenErrFmtStr, layerCompositeErr)
		}
	}

	flattenErr := bg.Flatten(&vips.Color{})
	if flattenErr != nil {
		return nil, fmt.Errorf(imgFlattenErrFmtStr, flattenErr)
	}

	return bg, nil
}

// © Arthur Gladfield
