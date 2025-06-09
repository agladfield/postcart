package img

import "github.com/davidbyttow/govips/v2/vips"

func Flatten(backdrop *vips.ImageRef, layers ...*vips.ImageRef) (*vips.ImageRef, error) {
	w := backdrop.Width()
	h := backdrop.Height()
	sceneBlank, sceneBlankErr := New(w, h, true)
	if sceneBlankErr != nil {
		return nil, sceneBlankErr
	}

	compositeOntoNewErr := sceneBlank.Composite(backdrop, vips.BlendModeOver, 0, 0)
	if compositeOntoNewErr != nil {
		return nil, compositeOntoNewErr
	}

	for _, layer := range layers {
		layerCopy, copyErr := layer.Copy()
		if copyErr != nil {
			sceneBlank.Close()
			return nil, copyErr
		}
		defer layerCopy.Close()

		layerCompositeErr := sceneBlank.Composite(layerCopy, vips.BlendModeOver, 0, 0)
		if layerCompositeErr != nil {
			sceneBlank.Close()
			return nil, layerCompositeErr
		}
	}

	return sceneBlank, nil
}
