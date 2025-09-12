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
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	// Image format initialization
	_ "image/jpeg"
	_ "image/png"

	// Image format initialization
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
// Can be sent directly to ConvertImage() for default ascii art
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
		FontFilePath:        "",
		FontColor:           [3]int{255, 255, 255},
		Braille:             false,
		Threshold:           128,
		Dither:              false,
	}
}

/*
Convert() takes an image or gif path/url as its first argument
and a aic_package.Flags literal as the second argument, with which it alters
the returned ascii art string.
*/
func Convert(flags Flags) (string, error) {

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
	fontPath = flags.FontFilePath
	fontColor = flags.FontColor
	braille = flags.Braille
	threshold = flags.Threshold
	dither = flags.Dither

	// Declared at the start since some variables are initially used in conditional blocks
	var (
		inputBytes []byte
		err             error
	)

	// Check file/data type of piped input

	if !isInputFromPipe() {
		return "", fmt.Errorf("there is no input being piped to stdin")
	}

	inputBytes, err = ioutil.ReadAll(os.Stdin)
	if err != nil {
		return "", fmt.Errorf("unable to read piped input: %v", err)
	}

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
		return "", errors.New("File type of piped input could not be determined, input may be malformed or not be one of the supported file types")
	}

	if inputIsGif {
		return "", pathIsGif(inputBytes)
	} else {
		return pathIsImage(inputBytes)
	}
}
