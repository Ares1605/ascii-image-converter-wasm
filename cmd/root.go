/*
Copyright © 2021 Zoraiz Hassan <hzoraiz8@gmail.com>
Copyright © 2025 Ares Stavropoulos <aresstav04@gmail.com>

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

package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"encoding/json"

	"github.com/Ares1605/ascii-image-converter-wasm/aic_package"
	image_conversions "github.com/Ares1605/ascii-image-converter-wasm/image_manipulation"

	"github.com/spf13/cobra"
)

var (
	// Flags
	complex       bool
	dimensions    []int
	width         int
	height        int
	negative      bool
	formatsTrue   bool
	colored       bool
	colorBg       bool
	grayscale     bool
	customMap     string
	flipX         bool
	flipY         bool
	jsonOutput    bool
	hundredsColor bool
	fontColor     []int
	braille       bool
	threshold     int
	dither        bool

	// Root commands
	rootCmd = &cobra.Command{
		Use:     "[piped input] | ascii-image-converter -",
		Short:   "Converts images and gifs into ascii art",
		Version: "1.13.1",
		Long:    "This tool converts images into ascii art and prints them on the terminal.\nFurther configuration can be managed with flags.",

		// Not RunE since help text is getting larger and seeing it for every error impacts user experience
		Run: func(cmd *cobra.Command, args []string) {

			if checkInputAndFlags(args) {
				return
			}

			flags := aic_package.Flags{
				Complex:             complex,
				Dimensions:          dimensions,
				Width:               width,
				Height:              height,
				Negative:            negative,
				Colored:             colored,
				CharBackgroundColor: colorBg,
				Grayscale:           grayscale,
				CustomMap:           customMap,
				FlipX:               flipX,
				FlipY:               flipY,
				JsonOutput:          jsonOutput,
				FontColor:           [3]int{fontColor[0], fontColor[1], fontColor[2]},
				Braille:             braille,
				Threshold:           threshold,
				Dither:              dither,
				// By default, color level is set to true (24-bit) color
				ColorLevel:          image_conversions.Millions,
			}
			if hundredsColor {
				flags.ColorLevel = image_conversions.Hundreds
			}

			// Check file/data type of piped input
			if !aic_package.IsInputFromPipe() {
				fmt.Printf("there is no input being piped to stdin\n")
				return
			}

			inputBytes, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				fmt.Printf("unable to read piped input: %v\n", err)
				return
			}

			if err = printAscii(inputBytes, flags); err != nil {
				fmt.Printf("%v\n", err)
			}
		},
	}
)

func printAscii(inputBytes []byte, flags aic_package.Flags) error {
	if flags.JsonOutput {
		if asciiArt, err := aic_package.ConvertJSON(inputBytes, flags); err == nil {
			marshalled, err := json.Marshal(asciiArt)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			}
			fmt.Printf("%s", marshalled)
		} else {
			fmt.Printf("Error: %v\n", err)
		}
	} else {
		if asciiArt, err := aic_package.Convert(inputBytes, flags); err == nil {
			fmt.Printf("%s", asciiArt)
		} else {
			fmt.Printf("Error: %v\n", err)
		}
	}
	fmt.Println()
	return nil
}

// Cobra configuration from here on

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().SortFlags = false
	rootCmd.Flags().SortFlags = false

	rootCmd.PersistentFlags().BoolVarP(&colored, "color", "C", false, "Display ascii art with original colors\nIf 24-bit colors aren't supported, uses 8-bit\n(Inverts with --negative flag)\n(Overrides --grayscale and --font-color flags)\n")
	rootCmd.PersistentFlags().BoolVar(&colorBg, "color-bg", false, "If some color flag is passed, use that color\non character background instead of foreground\n(Inverts with --negative flag)\n(Only applicable for terminal display)\n")
	rootCmd.PersistentFlags().IntSliceVarP(&dimensions, "dimensions", "d", nil, "Set width and height for ascii art in CHARACTER length\ne.g. -d 60,30 (defaults to terminal height)\n(Overrides --width and --height flags)\n")
	rootCmd.PersistentFlags().IntVarP(&width, "width", "W", 0, "Set width for ascii art in CHARACTER length\nHeight is kept to aspect ratio\ne.g. -W 60\n")
	rootCmd.PersistentFlags().IntVarP(&height, "height", "H", 0, "Set height for ascii art in CHARACTER length\nWidth is kept to aspect ratio\ne.g. -H 60\n")
	rootCmd.PersistentFlags().StringVarP(&customMap, "map", "m", "", "Give custom ascii characters to map against\nOrdered from darkest to lightest\ne.g. -m \" .-+#@\" (Quotation marks excluded from map)\n(Overrides --complex flag)\n")
	rootCmd.PersistentFlags().BoolVarP(&braille, "braille", "b", false, "Use braille characters instead of ascii\nTerminal must support braille patterns properly\n(Overrides --complex and --map flags)\n")
	rootCmd.PersistentFlags().IntVar(&threshold, "threshold", 0, "Threshold for braille art\nValue between 0-255 is accepted\ne.g. --threshold 170\n(Defaults to 128)\n")
	rootCmd.PersistentFlags().BoolVar(&dither, "dither", false, "Apply dithering on image for braille\nart conversion\n(Only applicable with --braille flag)\n(Negates --threshold flag)\n")
	rootCmd.PersistentFlags().BoolVarP(&grayscale, "grayscale", "g", false, "Display grayscale ascii art\n(Inverts with --negative flag)\n(Overrides --font-color flag)\n")
	rootCmd.PersistentFlags().BoolVarP(&complex, "complex", "c", false, "Display ascii characters in a larger range\nMay result in higher quality\n")
	rootCmd.PersistentFlags().BoolVarP(&negative, "negative", "n", false, "Display ascii art in negative colors\n")
	rootCmd.PersistentFlags().BoolVarP(&flipX, "flipX", "x", false, "Flip ascii art horizontally\n")
	rootCmd.PersistentFlags().BoolVarP(&flipY, "flipY", "y", false, "Flip ascii art vertically\n")
	rootCmd.PersistentFlags().BoolVarP(&jsonOutput, "json", "J", false, "Output ASCII image with JSON.\nFor programmable iteration where ANSI escape codes are not supported.\n")
	rootCmd.PersistentFlags().BoolVar(&hundredsColor, "256-color", false, "If some color flag is passed, sets the color output to 256 (8-bit) color, as opposed to true (24-bit) color.\nWeb APIs virtually exclusively support true (24-bit) color, however this color level exists to support mundane color, or environments incompatible with true (24-bit) color.\n")
	rootCmd.PersistentFlags().IntSliceVar(&fontColor, "font-color", nil, "Set font color for terminal\nPass an RGB value\ne.g. --font-color 0,0,0\n(Defaults to 255,255,255)\n")
	rootCmd.PersistentFlags().BoolVar(&formatsTrue, "formats", false, "Display supported input formats\n")

	rootCmd.PersistentFlags().BoolP("help", "h", false, "Help for "+rootCmd.Name()+"\n")
	rootCmd.PersistentFlags().BoolP("version", "v", false, "Version for "+rootCmd.Name())

	rootCmd.SetVersionTemplate("{{printf \"v%s\" .Version}}\n")

	defaultUsageTemplate := rootCmd.UsageTemplate()
	rootCmd.SetUsageTemplate(defaultUsageTemplate + "\nCopyright © 2021 Zoraiz Hassan <hzoraiz8@gmail.com>\n" +
		"Copyright © 2025 Ares Stavropoulos <aresstav04@gmail.com>\n" +
		"Distributed under the Apache License Version 2.0 (Apache-2.0)\n" +
		"For further details, visit https://github.com/Ares1065/ascii-image-converter-wasm\n")
}
