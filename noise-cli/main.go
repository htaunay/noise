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
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"os/user"
	"strings"

	// 3rd party
	"github.com/spf13/cobra"
	//"github.com/spf13/viper"

	// internal
	"github.com/htaunay/noise"
)

func main() {
	Execute()
}

// Parameter vars
var imageSize uint
var octaves uint
var frequency float64
var lacunarity float64
var persistence float64
var xOffset float64
var yOffset float64
var channels uint
var outputFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "noise-cli",
	Short: "A command line interface for the noise package",
	// Long: TODO
	Run: func(cmd *cobra.Command, args []string) {
		// empty
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	opts := noise.NoiseOptions{
		Size:        imageSize,
		Octaves:     octaves,
		Frequency:   frequency,
		Lacunarity:  lacunarity,
		Persistence: persistence,
		XOffset:     xOffset,
		YOffset:     yOffset,
		Channels:    channels,
	}

	validationErrors := validateParams(opts)
	if len(validationErrors) > 0 {
		fmt.Println("The following given parameters are invalid:")
		for _, v := range validationErrors {
			fmt.Println("\t* ", v)
		}
		os.Exit(1)
	}

	matrix := noise.Build(opts)
	img := matrix2img(matrix)

	// Get user info
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	// Solve HOME dir
	outputFile = strings.Replace(outputFile, "~", usr.HomeDir, 1)
	// Add png, if missing
	if !strings.HasSuffix(outputFile, ".png") {
		outputFile += ".png"
	}

	// Write to disk
	file, err := os.Create(outputFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	png.Encode(file, img)
}

// Get called automatically on startup
func init() {

	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().UintVarP(
		&imageSize,
		"size",
		"s",
		640,
		"Size, in pixels, of the square-shaped generated image",
	)

	RootCmd.PersistentFlags().UintVarP(
		&octaves,
		"octaves",
		"o",
		2,
		`Number of times noise functions with varying amplitude and frequencies
		will be added. When set to 1, only one layer is calculated`,
	)

	RootCmd.PersistentFlags().Float64VarP(
		&frequency,
		"frequency",
		"q",
		16.0,
		"Period at which the data will be sampled",
	)

	RootCmd.PersistentFlags().Float64VarP(
		&lacunarity,
		"lacunarity",
		"l",
		2.5,
		"Multiplier that determines how quickly the frequency increases for each successive octave",
	)

	RootCmd.PersistentFlags().Float64VarP(
		&persistence,
		"persistence",
		"p",
		0.75,
		"Multiplier for how quickly the amplitudes diminish for each successive octave",
	)

	RootCmd.PersistentFlags().Float64VarP(
		&xOffset,
		"xoffset",
		"x",
		0,
		"Ratio of horizontal offset (default 0.0)",
	)

	RootCmd.PersistentFlags().Float64VarP(
		&yOffset,
		"yoffset",
		"y",
		0,
		"Ratio of vertical offset (default 0.0)",
	)

	RootCmd.PersistentFlags().UintVarP(
		&channels,
		"channels",
		"j",
		1,
		"Number of channels to try to break the noise computation into and run in parallel",
	)

	RootCmd.PersistentFlags().StringVarP(
		&outputFile,
		"file",
		"f",
		"noise.png",
		"Output file path",
	)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// empty
}

func validateParams(opts noise.NoiseOptions) []string {

	var errors []string

	if opts.Size < 4 {
		errors = append(errors, "Image size is less than the minimum of 4")
	}

	if opts.Octaves < 1 {
		errors = append(errors, "Octaves count is less than the minimum of 1")
	}

	if opts.Frequency < 0.5 {
		errors = append(errors, "Frequency is less than the minimum of 0.5")
	}

	if opts.Lacunarity < 1 {
		errors = append(errors, "Lacunarity is less than the minimum of 1")
	}

	if opts.Persistence < 0 {
		errors = append(errors, "Persistence is less than the minimum of 0")
	}

	if opts.XOffset < 0 {
		errors = append(errors, "X-axis offset is less than the minimum of 0")
	}

	if opts.YOffset < 0 {
		errors = append(errors, "Y-axis offset is less than the minimum of 0")
	}

	if opts.Channels < 1 {
		errors = append(errors, "Channels is less than the minimum of 1")
	}

	return errors
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
