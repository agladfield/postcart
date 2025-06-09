package img

import (
	"os"

	"github.com/davidbyttow/govips/v2/vips"
)

func DebugSave(target *vips.ImageRef, path string) {
	bytes, _, exportErr := target.ExportPng(&vips.PngExportParams{
		Quality:     75,
		Compression: 2,
	})
	if exportErr != nil {
		panic(exportErr)
	}

	writeErr := os.WriteFile(path, bytes, 0677)
	if writeErr != nil {
		panic(writeErr)
	}

	return
}

func DebugSaveJPG(target *vips.ImageRef, path string) {
	bytes, _, exportErr := target.ExportJpeg(&vips.JpegExportParams{
		Quality: 70,
	})
	if exportErr != nil {
		panic(exportErr)
	}

	writeErr := os.WriteFile(path, bytes, 0677)
	if writeErr != nil {
		panic(writeErr)
	}
}
