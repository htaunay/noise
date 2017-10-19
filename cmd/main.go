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

	// 3rd party
	"github.com/spf13/cobra"
	//"github.com/spf13/viper"

	// internal
	"github.com/htaunay/noise"
)

func main() {
	Execute()
}

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "mapgenerator",
	Short: "A command line random map generator",
	Long: `A command line random map generator, based on the Simplex Noise
    algorithm. The output is in PNG, and only 2D resolution is supported.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ops := noise.NoiseOptions{
		Size:        640,
		Octaves:     2,
		Frequency:   16.0,
		Lacunarity:  2.5,
		Persistence: 0.75,
	}

	matrix := noise.Build(ops)
	img := matrix2img(matrix)

	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	file, err := os.Create(usr.HomeDir + "/Desktop/perlin.png")
	if err != nil {
		panic(err)
	}

	defer file.Close()
	png.Encode(file, img)
}

func init() {
	// empty
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// empty
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
