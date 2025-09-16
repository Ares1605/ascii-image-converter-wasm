/*
Copyright © 2021 Zoraiz Hassan <hzoraiz8@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package image_conversions

import (
	"fmt"
	"image"
	"errors"
	"image/color"

	"github.com/disintegration/imaging"
	gookitColor "github.com/gookit/color"
	"github.com/makeworld-the-better-one/dither/v2"
)

func ditherImage(img image.Image) image.Image {

	palette := []color.Color{
		color.Black,
		color.White,
	}

	d := dither.NewDitherer(palette)
	d.Matrix = dither.FloydSteinberg

	return d.DitherCopy(img)
}

func resizeImage(img image.Image, isBraille bool, dimensions []int, width, height int) (image.Image, error) {

	var asciiWidth, asciiHeight int
	var smallImg image.Image

	imgWidth := float64(img.Bounds().Dx())
	imgHeight := float64(img.Bounds().Dy())
	aspectRatio := imgWidth / imgHeight

	if (width != 0 || height != 0) && len(dimensions) == 0 {
		// If either width or height is set and dimensions aren't given

		if width != 0 && height == 0 {
			// If width is set and height is not set, use width to calculate aspect ratio

			asciiWidth = width
			asciiHeight = int(float64(asciiWidth) / aspectRatio)
			asciiHeight = int(0.5 * float64(asciiHeight))

			if asciiHeight == 0 {
				asciiHeight = 1
			}

		} else if height != 0 && width == 0 {
			// If height is set and width is not set, use height to calculate aspect ratio

			asciiHeight = height
			asciiWidth = int(float64(asciiHeight) * aspectRatio)
			asciiWidth = int(2 * float64(asciiWidth))

			if asciiWidth == 0 {
				asciiWidth = 1
			}

		} else {
			return nil, fmt.Errorf("error: both width and height can't be set. Use dimensions instead")
		}

	} else if len(dimensions) == 0 {
		return nil, errors.New("Either dimensions, width, or height must be passed.")
	} else {
		// Else, set passed dimensions

		asciiWidth = dimensions[0]
		asciiHeight = dimensions[1]
	}

	// Because one braille character has 8 dots (4 rows and 2 columns)
	if isBraille {
		asciiWidth *= 2
		asciiHeight *= 4
	}
	smallImg = imaging.Resize(img, asciiWidth, asciiHeight, imaging.Lanczos)

	return smallImg, nil
}

func reverse(imgSet [][]AsciiPixel, flipX, flipY bool) [][]AsciiPixel {

	if flipX {
		for _, row := range imgSet {
			for i, j := 0, len(row)-1; i < j; i, j = i+1, j-1 {
				row[i], row[j] = row[j], row[i]
			}
		}
	}

	if flipY {
		for i, j := 0, len(imgSet)-1; i < j; i, j = i+1, j-1 {
			imgSet[i], imgSet[j] = imgSet[j], imgSet[i]
		}
	}

	return imgSet
}

// This functions calculates terminal color level between rgb colors and 256-colors
// and returns the character with escape codes appropriately
func getColoredCharForTerm(r, g, b uint8, char string, background bool, colorLevel ColorLevel) (string, gookitColor.RGBColor, error) {
	var coloredChar string

	switch colorLevel {
	case Millions:
		colorRenderer := gookitColor.RGB(uint8(r), uint8(g), uint8(b), background)
		coloredChar = colorRenderer.Sprintf("%v", char)
		return coloredChar, colorRenderer, nil
	case Hundreds:
		colorRenderer := gookitColor.RGB(uint8(r), uint8(g), uint8(b), background).C256()
		coloredChar = colorRenderer.Sprintf("%v", char)
		// after converting the RGB to C256, convert it back to RGB so we can have a 256 color normalized as an RGB color
		return coloredChar, colorRenderer.RGB(), nil
	default:
		return "", gookitColor.RGBColor{}, fmt.Errorf("%d-bit color level is unsupported.", colorLevel)
	}
}
