// Copyright © 2017 Henrique Taunay <henrique@taunay.me>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	// std packages
	// "fmt"
	"image"
	"image/color"

	// 3rd party
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

	// internal
	"github.com/htaunay/noise"
)

func main() {
	pixelgl.Run(openWindow)
}

// Initial values
var imageSize uint = 640
var octaves uint = 2
var frequency float64 = 16.0
var lacunarity float64 = 2.5
var persistence float64 = 0.75
var xOffset float64 = 0
var yOffset float64 = 0
var applyFilter bool = false

const octaveStep uint = 1
const frequencyStep float64 = 0.5
const lacunarityStep float64 = 0.1
const persistenceStep float64 = 0.05
const offsetStep float64 = 0.25

func openWindow() {

	cfg := pixelgl.WindowConfig{
		Title:  "Noise Generator",
		Bounds: pixel.R(0, 0, 640, 640),
		VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	opts := noise.NoiseOptions{
		Size:        imageSize,
		Octaves:     octaves,
		Frequency:   frequency,
		Lacunarity:  lacunarity,
		Persistence: persistence,
		XOffset:     xOffset,
		YOffset:     yOffset,
		Channels:    1,
	}

	matrix := noise.Build(opts)
	img := matrix2img(matrix)

	pic := pixel.PictureDataFromImage(img)
	sprite := pixel.NewSprite(pic, pic.Bounds())
	sprite.Draw(win, pixel.IM.Moved(win.Bounds().Center()))

	// main loop
	for !win.Closed() {

		// were there any changes?
		change := false

		// negative/postive control
		shift := false
		if win.Pressed(pixelgl.KeyLeftShift) || win.Pressed(pixelgl.KeyRightShift) {
			shift = true
		}

		// octave control
		if win.JustPressed(pixelgl.KeyO) {
			// fmt.Printf("Octave from %d ", opts.Octaves)
			if shift {
				opts.Octaves -= octaveStep
			} else {
				opts.Octaves += octaveStep
			}
			opts.Octaves = clampUint(int(opts.Octaves), 1, 8)
			// fmt.Printf("to %d\n", opts.Octaves)
			change = true
		}

		// frequency control
		if win.Pressed(pixelgl.KeyQ) {
			// fmt.Printf("Frequency from %f ", opts.Frequency)
			if shift {
				opts.Frequency -= frequencyStep
			} else {
				opts.Frequency += frequencyStep
			}
			opts.Frequency = clampFloat64(opts.Frequency, 0.1, 128.0)
			// fmt.Printf("to %f\n", opts.Frequency)
			change = true
		}

		// lacunarity control
		if win.Pressed(pixelgl.KeyL) {
			// fmt.Printf("Lacunarity from %f ", opts.Lacunarity)
			if shift {
				opts.Lacunarity -= lacunarityStep
			} else {
				opts.Lacunarity += lacunarityStep
			}
			opts.Lacunarity = clampFloat64(opts.Lacunarity, 1.0, 4.0)
			// fmt.Printf("to %f\n", opts.Lacunarity)
			change = true
		}

		// persistence control
		if win.Pressed(pixelgl.KeyP) {
			// fmt.Printf("Persistence from %f ", opts.Persistence)
			if shift {
				opts.Persistence -= persistenceStep
			} else {
				opts.Persistence += persistenceStep
			}
			opts.Persistence = clampFloat64(opts.Persistence, 0.0, 1.0)
			// fmt.Printf("to %f\n", opts.Persistence)
			change = true
		}

		// offset control
		xDiff := 0.0
		if win.Pressed(pixelgl.KeyLeft) {
			xDiff -= offsetStep
		}
		if win.Pressed(pixelgl.KeyRight) {
			xDiff += offsetStep
		}
		yDiff := 0.0
		if win.Pressed(pixelgl.KeyUp) {
			yDiff -= offsetStep
		}
		if win.Pressed(pixelgl.KeyDown) {
			yDiff += offsetStep
		}

		if xDiff != 0 || yDiff != 0 {
			// fmt.Printf("Offset from %.2f/%.2f ", opts.XOffset, opts.YOffset)
			x := opts.XOffset + xDiff
			y := opts.YOffset + yDiff
			opts.XOffset = clampFloat64(x, 0, 1024*1024)
			opts.YOffset = clampFloat64(y, 0, 1024*1024)
			// fmt.Printf("to %.2f/%.2f\n", opts.XOffset, opts.YOffset)
			change = true
		}

		// apply filter control
		if win.JustPressed(pixelgl.KeyF) {

			applyFilter = !applyFilter
			// fmt.Println("Toggling filter")
			change = true
		}

		// only update image if necessary
		if change {

			m := noise.Build(opts)
			i := matrix2img(m)

			p := pixel.PictureDataFromImage(i)
			s := pixel.NewSprite(p, p.Bounds())
			s.Draw(win, pixel.IM.Moved(win.Bounds().Center()))
		}

		win.Update()
	}
}

func matrix2img(m [][]uint8) *image.RGBA {

	imgHeight := len(m)
	imgWidth := len(m[0])

	img := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))
	for i := 0; i < imgHeight; i++ {
		for j := 0; j < imgWidth; j++ {
			scale := m[i][j]
			img.Set(j, i, getColor(scale))
		}
	}

	return img
}

func getColor(scale uint8) color.RGBA {

	if applyFilter {

		if scale < 127 {
			return color.RGBA{0, 0, 255, 255}
		} else if scale < 137 {
			return color.RGBA{239, 221, 111, 255}
		} else if scale < 180 {
			return color.RGBA{44, 176, 55, 255}
		} else {
			return color.RGBA{165, 42, 42, 255}
		}
	} else {
		return color.RGBA{scale, scale, scale, 255}
	}
}

func clampUint(value int, min uint, max uint) uint {

	if value < int(min) {
		return min
	}

	if value > int(max) {
		return max
	}

	return uint(value)
}

func clampFloat64(value, min, max float64) float64 {

	if value < min {
		return min
	}

	if value > max {
		return max
	}

	return value
}
