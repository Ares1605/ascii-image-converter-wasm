# ascii-image-converter-wasm

[![release-version](https://img.shields.io/github/v/release/Ares1605/ascii-image-converter-wasm?label=Latest%20Version)](https://github.com/Ares1605/ascii-image-converter-wasm/releases/latest)
[![license](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/Ares1605/ascii-image-converter-wasm/blob/master/LICENSE.txt)
[![language](https://img.shields.io/badge/Language-Go-blue)](https://golang.org/)
![release-downloads](https://img.shields.io/github/downloads/Ares1605/ascii-image-converter-wasm/total?color=1d872d&label=Release%20Downloads)

ascii-image-converter-wasm is a WASM compatible command-line tool that converts images into ascii art and prints them out onto the console. Available on Windows, Linux and macOS.

Now supports braille art!

Input formats currently supported:
* JPEG/JPG
* PNG
* BMP
* WEBP
* TIFF/TIF
* GIF

<p align="center">
  <img src="https://raw.githubusercontent.com/Ares1605/ascii-image-converter-wasm/master/example_gifs/all.gif">
</p>

## Table of Contents

-  [A Quick Note](#a-quick-note)
-  [Installation](#installation)
	*  [Go](#go)
	*  [Linux (binaries)](#linux)
	*  [Windows (binaries)](#windows)
-  [CLI Usage](#cli-usage)
	*  [Flags](#flags)
-  [Library Usage](#library-usage)
-  [Contributing](#contributing)
-  [Packages Used](#packages-used)
-  [License](#license)

## A Quick Note
Before reading the rest of the README, please understand the following:
* This repo was forked from [Ares1605/ascii-image-converter-wasm](https://github.com/Ares1605/ascii-image-converter-wasm), which did all the work in building an ASCII image rendern building an ASCII image renderer.
* It was forked to remove non-WASM-compatible features like config files, reading images from file path, and autosizing to terminal width. Fundamentally these features cannot be used or built in a WASM environment.

## Installation

### Go

```
go install github.com/Ares1605/ascii-image-converter-wasm@latest
```
<hr>

For physically installing the binaries, follow the steps with respect to your OS.

### Linux

Download the archive for your distribution's architecture [here](https://github.com/Ares1605/ascii-image-converter-wasm/releases/latest), extract it, and open the extracted directory.

Now, open a terminal in the same directory and execute this command:

```
sudo cp ascii-image-converter-wasm /usr/local/bin/
```
Now you can use ascii-image-converter-wasm in the terminal. Execute `ascii-image-converter-wasm -h` for more details.

### Windows

You will need to set an Environment Variable to the folder the ascii-image-converter-wasm.exe executable is placed in to be able to use it in the command prompt. Follow the instructions in case of confusion:

Download the archive for your Windows architecture [here](https://github.com/Ares1605/ascii-image-converter-wasm/releases/latest), extract it, and open the extracted folder. Now, copy the folder path from the top of the file explorer and follow these instructions:
* In Search, search for and then select: Advanced System Settings
* Click Environment Variables. In the section User Variables find the Path environment variable and select it. Click "Edit".
* In the Edit Environment Variable window, click "New" and then paste the path of the folder that you copied initially.
* Click "Ok" on all open windows.

Now, restart any open command prompt and execute `ascii-image-converter-wasm -h` for more details.

<br>

## CLI Usage

> **Note:** Decrease font size or increase terminal width (like zooming out) for maximum quality ascii art

The basic usage for converting an image into ascii art is as follows. You can also supply multiple image paths and urls as well as a GIF.

```
[piped input] | ascii-image-converter-wasm -
```
Example:
```
myImage.jpeg | ascii-image-converter-wasm -
```

### Flags

#### --color OR -C

> **Note:** Your terminal must support 24-bit or 8-bit colors for appropriate results. If 24-bit colors aren't supported, 8-bit color escape codes will be used

Display ascii art with the colors from original image.

```
[piped input] | ascii-image-converter-wasm -W <width> -C -
# Or
[piped input] | ascii-image-converter-wasm -W <width> --color -
```

<p align="center">
  <img src="https://raw.githubusercontent.com/Ares1605/ascii-image-converter-wasm/master/example_gifs/color.gif">
</p>

#### --braille OR -b

> **Note:** Braille pattern display heavily depends on which terminal or font you're using. In windows, try changing the font from command prompt properties if braille characters don't display

Use braille characters instead of ascii. For this flag, your terminal must support braille patters (UTF-8) properly. Otherwise, you may encounter problems with colored or even uncolored braille art.
```
[piped input] | ascii-image-converter-wasm -W <width> -b -
# Or
[piped input] | ascii-image-converter-wasm -W <width> --braille -
```

<p align="center">
  <img src="https://raw.githubusercontent.com/Ares1605/ascii-image-converter-wasm/master/example_gifs/braille.gif">
</p>

#### --threshold

Set threshold value to compare for braille art when converting each pixel into a dot. Value must be between 0 and 255.

Example:
```
[piped input] | ascii-image-converter-wasm -W <width> -b --threshold 170 -
```

#### --dither

Apply dithering on image to make braille art more visible. Since braille dots can only be on or off, dithering images makes them more visible in braille art.

Example:
```
[piped input] | ascii-image-converter-wasm -W <width> -b --dither -
```

<p align="center">
  <img src="https://raw.githubusercontent.com/Ares1605/ascii-image-converter-wasm/master/example_gifs/dither.gif">
</p>

#### --color-bg

If any of the coloring flags is passed, this flag will transfer its color to each character's background. instead of foreground.
```
[piped input] | ascii-image-converter-wasm -W <width> -C --color-bg -
```

#### --dimensions OR -d

> **Note:** Don't immediately append another flag with -d

Set the width and height for ascii art in CHARACTER lengths.
```
[piped input] | ascii-image-converter-wasm -d <width>,<height> -
# Or
[piped input] | ascii-image-converter-wasm --dimensions <width>,<height> -
```
Example:
```
[piped input] | ascii-image-converter-wasm -d 60,30 -
```
<p align="center">
  <img src="https://raw.githubusercontent.com/Ares1605/ascii-image-converter-wasm/master/example_gifs/dimensions.gif">
</p>

#### --width OR -W

> **Note:** Don't immediately append another flag with -W

Set width of ascii art. Height is calculated according to aspect ratio.
```
[piped input] | ascii-image-converter-wasm -W <width> -
# Or
[piped input] | ascii-image-converter-wasm --width <width> -
```
Example:
```
[piped input] | ascii-image-converter-wasm -W 60 -
```

#### --height OR -H

> **Note:** Don't immediately append another flag with -H

Set height of ascii art. Width is calculated according to aspect ratio.
```
[piped input] | ascii-image-converter-wasm -H <height> -
# Or
[piped input] | ascii-image-converter-wasm --height <height> -
```
Example:
```
[piped input] | ascii-image-converter-wasm -H 60 -
```

#### --map OR -m

> **Note:** Don't immediately append another flag with -m

Pass a string of your own ascii characters to map against. Passed characters must start from darkest character and end with lightest. There is no limit to number of characters.

Empty spaces can be passed if string is passed inside quotation marks. You can use both single or double quote for quotation marks. For repeating quotation mark inside string, append it with \ (such as  \\").

```
[piped input] | ascii-image-converter-wasm -W <width> -m "<string-of-characters>" -
# Or
[piped input] | ascii-image-converter-wasm -W <width> --map "<string-of-characters>" -
```
Following example contains 7 depths of lighting.
```
[piped input] | ascii-image-converter-wasm -W <width> -m " .-=+#@" -
```

<p align="center">
  <img src="https://raw.githubusercontent.com/Ares1605/ascii-image-converter-wasm/master/example_gifs/map.gif">
</p>

#### --grayscale OR -g

Display ascii art in grayscale colors. This is the same as --color flag, except each character will be encoded with a grayscale RGB value.

```
[piped input] | ascii-image-converter-wasm -W <width> -g -
# Or
[piped input] | ascii-image-converter-wasm -W <width> --grayscale -
```

#### --negative OR -n

Display ascii art in negative colors. Works with both uncolored and colored text from --color flag.

```
[piped input] | ascii-image-converter-wasm -W <width> -n -
# Or
[piped input] | ascii-image-converter-wasm -W <width> -negative -
```

<p align="center">
  <img src="https://raw.githubusercontent.com/Ares1605/ascii-image-converter-wasm/master/example_gifs/negative.gif">
</p>

#### --complex OR -c

Print the image with a wider array of ascii characters for more detailed lighting density. Sometimes improves accuracy.
```
[piped input] | ascii-image-converter-wasm -W <width> -c -
# Or
[piped input] | ascii-image-converter-wasm -W <width> --complex -
```

#### --flipX OR -x

Flip the ascii art horizontally on the terminal.

```
[piped input] | ascii-image-converter-wasm -W <width> --flipX -
# Or
[piped input] | ascii-image-converter-wasm -W <width> -x -
```

#### --flipY OR -y
Flip the ascii art vertically on the terminal.

```
[piped input] | ascii-image-converter-wasm -W <width> --flipY -
# Or
[piped input] | ascii-image-converter-wasm -W <width> -y -
```



#### --font-color

This flag takes an RGB value that sets the font color to the displayed ascii art in terminal.

```
[piped input] | ascii-image-converter-wasm -W <width> --font-color 0,0,0 # For black font color
```

#### --formats

Display supported input formats.

```
ascii-image-converter-wasm --formats
```

<br>

## Library Usage

First, install the library with:
```
go get -u github.com/Ares1605/ascii-image-converter-wasm/aic_package
```

For an image:

```go
package main

import (
	"fmt"

	"github.com/Ares1605/ascii-image-converter-wasm/aic_package"
)

func main() {
	// If file is in current directory. This can also be a URL to an image or gif.
	filePath := "myImage.jpeg"

	flags := aic_package.DefaultFlags()

	// This part is optional.
	// You can directly pass default flags variable to aic_package.Convert() if you wish.
	// There are more flags, but these are the ones shown for demonstration
	flags.Dimensions = []int{50, 25}
	flags.Colored = true
	flags.CustomMap = " .-=+#@"
	flags.FontFilePath = "./RobotoMono-Regular.ttf" // If file is in current directory
	flags.SaveBackgroundColor = [4]int{50, 50, 50, 100}

	// Note: For environments where a terminal isn't available (such as web servers), you MUST
	// specify atleast one of flags.Width, flags.Height or flags.Dimensions

	// Conversion for an image
	asciiArt, err := aic_package.Convert(filePath, flags)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%v\n", asciiArt)
}
```
<br>

> **Note:** GIF conversion is not advised as the function may run infinitely, depending on the GIF.

For a GIF:

```go
package main

import (
	"fmt"

	"github.com/Ares1605/ascii-image-converter-wasm/aic_package"
)

func main() {
	filePath = "myGif.gif"

	flags := aic_package.DefaultFlags()

	_, err := aic_package.Convert(filePath, flags)
	if err != nil {
		fmt.Println(err)
	}
}

```

<br>

## Contributing

You can fork the project and implement any changes you want for a pull request. However, for major changes, please open an issue first to discuss what you would like to implement.
PS. This is the forked project!

## Packages Used

[github.com/spf13/cobra](https://github.com/spf13/cobra)

[github.com/fogleman/gg](https://github.com/fogleman/gg)

[github.com/mitchellh/go-homedir](https://github.com/mitchellh/go-homedir)

[github.com/disintegration/imaging](https://github.com/disintegration/imaging)

[github.com/gookit/color](https://github.com/gookit/color)

[github.com/makeworld-the-better-one/dither](https://github.com/makeworld-the-better-one/dither)

## License

[Apache-2.0](https://github.com/TheZoraiz/ascii-image-converter/blob/master/LICENSE.txt)
