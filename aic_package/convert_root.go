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

package aic_package

import (
	"fmt"
	"net/http"

	// Image format initialization
	_ "image/jpeg"
	_ "image/png"

	// Image format initialization
	image_conversions "github.com/Ares1605/ascii-image-converter-wasm/image_manipulation"
	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"
)

var pipedInputTypes = []string{
	"image/png",
	"image/jpeg",
	"image/webp",
	"image/tiff",
	"image/bmp",
}

// Return default configuration for flags.
// Can be sent directly to Convert() for default ascii art
func DefaultFlags() Flags {
	return Flags{
		Complex:             false,
		Dimensions:          nil,
		Width:               0,
		Height:              0,
		Negative:            false,
		Colored:             false,
		CharBackgroundColor: false,
		Grayscale:           false,
		CustomMap:           "",
		FlipX:               false,
		FlipY:               false,
		FontColor:           [3]int{255, 255, 255},
		Braille:             false,
		Threshold:           128,
		Dither:              false,
		ColorLevel:          image_conversions.Millions,
	}
}

/*
Convert() takes a bytes array of the image/gif as its first argument
and a aic_package.Flags literal as the second argument, with which it alters
the returned ascii art string.
*/
func parseMetadata (inputBytes []byte, flags Flags) error {
	if flags.Dimensions == nil {
		dimensions = nil
	} else {
		dimensions = flags.Dimensions
	}
	width = flags.Width
	height = flags.Height
	complex = flags.Complex
	negative = flags.Negative
	colored = flags.Colored
	colorBg = flags.CharBackgroundColor
	grayscale = flags.Grayscale
	customMap = flags.CustomMap
	flipX = flags.FlipX
	flipY = flags.FlipY
	jsonOutput = flags.JsonOutput
	fontColor = flags.FontColor
	braille = flags.Braille
	threshold = flags.Threshold
	dither = flags.Dither
	colorLevel = flags.ColorLevel

	fileType := http.DetectContentType(inputBytes)
	invalidInput := true

	if fileType == "image/gif" {
		inputIsGif = true
		invalidInput = false
	} else {
		for _, inputType := range pipedInputTypes {
			if fileType == inputType {
				invalidInput = false
				break
			}
		}
	}

	if invalidInput {
		return fmt.Errorf("File type of piped input could not be determined, input may be malformed or not be one of the supported file types")
	}
	return nil
}
func Convert(inputBytes []byte, flags Flags) (string, error) {
	// Force JsonOutput to false
	flags.JsonOutput = false
	// parseMetadata mutates the inputIsGif global
	if err := parseMetadata(inputBytes, flags); err != nil {
		return "", err
	}
	if inputIsGif {
		return "", pathIsGif(inputBytes)
	} else {
		return pathIsImage(inputBytes, flattenToAscii)
	}
}
func ConvertJSON(inputBytes []byte, flags Flags) ([][]ColoredChar, error) {
	// Force Jsonoutput to true
	flags.JsonOutput = true
	// parseMetadata mutates the inputIsGif global
	if err := parseMetadata(inputBytes, flags); err != nil {
		return [][]ColoredChar{}, err
	}
	if inputIsGif {
		return [][]ColoredChar{}, fmt.Errorf("JSON output is not supported with GIFs.")
	} else {
		return pathIsImage(inputBytes, flattenToJSONable)
	}
}
