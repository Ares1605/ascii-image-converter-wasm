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
	"image/gif"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"

	imgManip "github.com/Ares1605/ascii-image-converter-wasm/image_manipulation"
)

type GifFrame struct {
	asciiCharSet [][]imgManip.AsciiChar
	delay        int
}

/*
This function grabs each image frame from passed gif and turns it into ascii art. If SaveGifPath flag is passed,
it'll turn each ascii art into an image instance of the same dimensions as the original gif and save them
as an ascii art gif.

Multi-threading has been implemented in multiple places due to long execution time
*/
func pathIsGif(inputBytes []byte) error {

	var (
		originalGif *gif.GIF
		err         error
	)

	originalGif, err = gif.DecodeAll(bytes.NewReader(inputBytes))
	if err != nil {
		return fmt.Errorf("Can't decode input: %v", err)
	}

	var (
		asciiArtSet    = make([]string, len(originalGif.Image))
		gifFramesSlice = make([]GifFrame, len(originalGif.Image))

		counter             = 0
		concurrentProcesses = 0
		wg                  sync.WaitGroup
		hostCpuCount        = runtime.NumCPU()
	)

	fmt.Printf("Generating ascii art... 0%%\r")

	// Get first frame of gif and its dimensions
	firstGifFrame := originalGif.Image[0].SubImage(originalGif.Image[0].Rect)
	firstGifFrameWidth := firstGifFrame.Bounds().Dx()
	firstGifFrameHeight := firstGifFrame.Bounds().Dy()

	// Multi-threaded loop to decrease execution time
	for i, frame := range originalGif.Image {

		wg.Add(1)
		concurrentProcesses++

		go func(i int, frame *image.Paletted) {

			frameImage := frame.SubImage(frame.Rect)

			// If a frame is found that is smaller than the first frame, then this gif contains smaller subimages that are
			// positioned inside the original gif. This behavior isn't supported by this app
			if firstGifFrameWidth != frameImage.Bounds().Dx() || firstGifFrameHeight != frameImage.Bounds().Dy() {
				fmt.Printf("Error: Input contains subimages smaller than default width and height\n\nProcess aborted because ascii-image-converter doesn't support subimage placement and transparency in GIFs\n\n")
				os.Exit(0)
			}

			var imgSet [][]imgManip.AsciiPixel

			imgSet, err = imgManip.ConvertToAsciiPixels(frameImage, dimensions, width, height, flipX, flipY, braille, dither)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(0)
			}

			var asciiCharSet [][]imgManip.AsciiChar
			if braille {
				asciiCharSet, err = imgManip.ConvertToBrailleChars(imgSet, negative, colored, grayscale, colorBg, fontColor, threshold, colorLevel)
			} else {
				asciiCharSet, err = imgManip.ConvertToAsciiChars(imgSet, negative, colored, grayscale, complex, colorBg, customMap, fontColor, colorLevel)
			}
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(0)
			}

			gifFramesSlice[i].asciiCharSet = asciiCharSet
			gifFramesSlice[i].delay = originalGif.Delay[i]

			asciiArtSet[i] = flattenToAscii(asciiCharSet, colored || grayscale)

			counter++
			percentage := int((float64(counter) / float64(len(originalGif.Image))) * 100)
			fmt.Printf("Generating ascii art... " + strconv.Itoa(percentage) + "%%\r")

			wg.Done()

		}(i, frame)

		// Limit concurrent processes according to host's CPU count to avoid overwhelming memory
		if concurrentProcesses == hostCpuCount {
			wg.Wait()
			concurrentProcesses = 0
		}
	}

	wg.Wait()
	fmt.Printf("                              \r")

	// Display the gif
	loopCount := 0
	for {
		for i, asciiFrame := range asciiArtSet {
			clearScreen()
			fmt.Println(asciiFrame)
			time.Sleep(time.Duration((time.Second * time.Duration(originalGif.Delay[i])) / 100))
		}

		// If gif is infinite loop
		if originalGif.LoopCount == 0 {
			continue
		}

		loopCount++
		if loopCount == originalGif.LoopCount {
			break
		}
	}

	return nil
}
