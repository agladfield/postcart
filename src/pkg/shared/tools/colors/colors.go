// Package colors wraps an RGBA implementation with helper methods
package colors

import (
	"fmt"
	"image"
	"image/color"
	"io/fs"
	"math"
	"strings"

	color_extractor "github.com/marekm4/color-extractor"
)

func to8Bit(c color.Color) (uint8, uint8, uint8, uint8) {
	r16, g16, b16, a16 := c.RGBA()
	return uint8(r16 / 257), uint8(g16 / 257), uint8(b16 / 257), uint8(a16 / 257)
}

type Color struct {
	R uint8
	G uint8
	B uint8
	A uint8
}

func (c Color) HexString(alpha ...bool) string {
	if len(alpha) > 1 && alpha[0] {
		return strings.ToUpper(fmt.Sprintf("#%02x%02x%02x%02x", c.R, c.G, c.B, c.A))
	} else {
		return strings.ToUpper(fmt.Sprintf("#%02x%02x%02x", c.R, c.G, c.B))
	}
}

func (c Color) RGB() (uint8, uint8, uint8) {
	return c.R, c.G, c.B
}

func (c Color) RGBA() (uint8, uint8, uint8, uint8) {
	return c.R, c.G, c.B, c.A
}

func ExtractColors(imgFile fs.File) ([]Color, error) {
	image, _, imageErr := image.Decode(imgFile)
	if imageErr != nil {
		return nil, imageErr
	}

	colors16bit := color_extractor.ExtractColorsWithConfig(image, color_extractor.Config{
		SmallBucket: .05,
		DownSizeTo:  256.,
	})

	colors8bit := make([]Color, len(colors16bit))

	for c, color := range colors16bit {
		r, g, b, a := to8Bit(color)
		colors8bit[c] = Color{r, g, b, a}
	}

	return colors8bit, nil
}

func isCloseToWhite(color Color) bool {
	const threshold = 10
	const minBrightness = 200

	r, g, b := float64(color.R), float64(color.G), float64(color.B)

	if math.Abs(r-g) <= threshold && math.Abs(g-b) <= threshold && math.Abs(b-r) <= threshold {
		return r >= minBrightness && g >= minBrightness && b >= minBrightness
	}
	return false
}

func isCloseToBlack(color Color) bool {
	const threshold = 10
	const maxBrightness = 60

	r, g, b := float64(color.R), float64(color.G), float64(color.B)

	if math.Abs(r-g) <= threshold && math.Abs(g-b) <= threshold && math.Abs(b-r) <= threshold {
		return r <= maxBrightness && g <= maxBrightness && b <= maxBrightness
	}
	return false
}

func FindFirstNonWhiteColors(colors []Color) (Color, Color, error) {
	if len(colors) < 1 {
		return Color{}, Color{}, fmt.Errorf("at least one color is required")
	}

	var firstNonWhite Color
	firstNonWhiteIndex := -1
	for i, color := range colors {
		if !isCloseToWhite(color) {
			firstNonWhite = color
			firstNonWhiteIndex = i
			break
		}
	}

	if firstNonWhiteIndex == -1 {
		return Color{}, Color{}, fmt.Errorf("no non-white color found")
	}

	for i := firstNonWhiteIndex + 1; i < len(colors); i++ {
		if !isCloseToWhite(colors[i]) && !isCloseToBlack(colors[i]) {
			return firstNonWhite, colors[i], nil
		}
	}

	if !isCloseToBlack(firstNonWhite) {
		return firstNonWhite, firstNonWhite, nil
	}

	for i := firstNonWhiteIndex + 1; i < len(colors); i++ {
		if !isCloseToWhite(colors[i]) {
			return firstNonWhite, colors[i], nil
		}
	}

	return firstNonWhite, firstNonWhite, nil
}

func GetDesiredColors(colors []Color) (Color, Color, error) {
	primColor := colors[0]
	secColor := colors[0]
	if len(colors) > 2 {
		var leastWhiteErr error
		primColor, secColor, leastWhiteErr = FindFirstNonWhiteColors(colors)
		if leastWhiteErr != nil {
			return Color{}, Color{}, leastWhiteErr
		}
	} else if len(colors) > 1 {
		if isCloseToWhite(colors[0]) {
			primColor = colors[1]
		}
		if !isCloseToWhite(colors[1]) && !isCloseToBlack(colors[1]) {
			secColor = colors[1]
		}
	}

	return primColor, secColor, nil
}

// © Arthur Gladfield
