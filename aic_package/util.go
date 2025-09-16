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
	"os"
	"strings"
	"os/exec"
	"runtime"

	gookitColor "github.com/gookit/color"
	imgManip "github.com/Ares1605/ascii-image-converter-wasm/image_manipulation"
)

// flattenToAscii flattens a two-dimensional grid of ascii characters into a string
// of ascii, with ANSI color codes
func flattenToAscii(asciiSet [][]imgManip.AsciiChar, colored bool) string {
	var ascii []string

	for _, line := range asciiSet {
		var tempAscii string

		for _, char := range line {
			if colored {
				tempAscii += char.OriginalColor
			} else if fontColor != [3]int{255, 255, 255} {
				tempAscii += char.SetColor
			} else {
				tempAscii += char.Simple
			}
		}

		ascii = append(ascii, tempAscii)
	}

	return strings.Join(ascii, "\n")
}

type ColoredChar struct {
	Char          string `json:"char"`
	RGBColor      *gookitColor.RGBColor `json:"rgb"`
}

// flattenToJSONable flattens the asciiSet by simplifying the set to only what's required in understanding
// each character and it's respective color
func flattenToJSONable(asciiSet [][]imgManip.AsciiChar, colored bool) [][]ColoredChar {
	simplified := make([][]ColoredChar, len(asciiSet))

	for i, line := range asciiSet {
		simplifiedLine := make([]ColoredChar, len(asciiSet[i]))

		for i, char := range line {
			if colored {
				simplifiedLine[i] = ColoredChar{
					Char: char.Simple,
					RGBColor: &char.OriginalColorRGB,
				}
			} else if fontColor != [3]int{255, 255, 255} {
				simplifiedLine[i] = ColoredChar{
					Char: char.Simple,
					RGBColor: &char.SetColorRGB,
				}
			} else {
				simplifiedLine[i] = ColoredChar{
					Char: char.Simple,
					// If no color is requsted, we set RGBColor to nil, and
					// trust the end-user handles this appropriately as "we don't know,
					// and we don't care".
					RGBColor: nil,
				}
			}
		}

		simplified[i] = simplifiedLine
	}

	return simplified
}

// Following is for clearing screen when showing gif
var clear map[string]func()

func init() {
	clear = make(map[string]func())
	clear["linux"] = func() {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["darwin"] = clear["linux"]
}

func clearScreen() {
	value, ok := clear[runtime.GOOS]
	if ok {
		value()
	} else {
		fmt.Println("Error: your platform is unsupported, terminal can't be cleared")
		os.Exit(0)
	}
}

func IsInputFromPipe() bool {
	fileInfo, _ := os.Stdin.Stat()
	return fileInfo.Mode()&os.ModeCharDevice == 0
}
