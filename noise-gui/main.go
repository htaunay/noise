package main

import (
	// std packages
	"fmt"
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

const octaveStep uint = 1
const frequencyStep float64 = 0.5
const lacunarityStep float64 = 0.1
const persistenceStep float64 = 0.05

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
			if shift {
				opts.Octaves -= octaveStep
			} else {
				opts.Octaves += octaveStep
			}
			fmt.Printf("Octave from %d ", opts.Octaves)
			opts.Octaves = clampUint(opts.Octaves, 1, 8)
			fmt.Printf("to %d\n", opts.Octaves)
			change = true
		}

		// frequency control
		if win.Pressed(pixelgl.KeyQ) {
			if shift {
				opts.Frequency -= frequencyStep
			} else {
				opts.Frequency += frequencyStep
			}
			fmt.Printf("Frequency from %f ", opts.Frequency)
			opts.Frequency = clampFloat64(opts.Frequency, 0.1, 128.0)
			fmt.Printf("to %f\n", opts.Frequency)
			change = true
		}

		// lacunarity control
		if win.Pressed(pixelgl.KeyL) {
			if shift {
				opts.Lacunarity -= lacunarityStep
			} else {
				opts.Lacunarity += lacunarityStep
			}
			fmt.Printf("Lacunarity from %f ", opts.Lacunarity)
			opts.Lacunarity = clampFloat64(opts.Lacunarity, 1.0, 4.0)
			fmt.Printf("to %f\n", opts.Lacunarity)
			change = true
		}

		// persistence control
		if win.Pressed(pixelgl.KeyP) {
			if shift {
				opts.Persistence -= persistenceStep
			} else {
				opts.Persistence += persistenceStep
			}
			fmt.Printf("Persistence from %f ", opts.Persistence)
			opts.Persistence = clampFloat64(opts.Persistence, 0.0, 1.0)
			fmt.Printf("to %f\n", opts.Persistence)
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
			img.Set(j, i, color.RGBA{scale, scale, scale, 255})
		}
	}

	return img
}

func clampUint(value, min, max uint) uint {

	if value < min {
		return min
	}

	if value > max {
		return max
	}

	return value
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
