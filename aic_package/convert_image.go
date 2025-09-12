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
	"bytes"
	"fmt"
	"image"
	"strings"

	imgManip "github.com/Ares1605/ascii-image-converter-wasm/image_manipulation"
)

// This function decodes the passed image and returns an ascii art string, optionaly saving it as a .txt and/or .png file
func pathIsImage(pipedInputBytes []byte) (string, error) {

	var (
		imData image.Image
		err    error
	)

	imData, _, err = image.Decode(bytes.NewReader(pipedInputBytes))
	if err != nil {
		return "", fmt.Errorf("Can't decode input: %v", err)
	}

	imgSet, err := imgManip.ConvertToAsciiPixels(imData, dimensions, width, height, flipX, flipY, braille, dither)
	if err != nil {
		return "", err
	}

	var asciiSet [][]imgManip.AsciiChar

	if braille {
		asciiSet, err = imgManip.ConvertToBrailleChars(imgSet, negative, colored, grayscale, colorBg, fontColor, threshold)
	} else {
		asciiSet, err = imgManip.ConvertToAsciiChars(imgSet, negative, colored, grayscale, complex, colorBg, customMap, fontColor)
	}
	if err != nil {
		return "", err
	}

	ascii := flattenAscii(asciiSet, colored || grayscale, false)
	result := strings.Join(ascii, "\n")

	return result, nil
}
