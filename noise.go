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

package noise

import "math"

type NoiseOptions struct {
	Size        uint
	Octaves     uint
	Frequency   float64
	Lacunarity  float64
	Persistence float64
	XOffset     float64
	YOffset     float64
	Channels    uint
}

func Build(opts NoiseOptions) [][]uint8 {

	// allocate space
	matrix := make([][]uint8, opts.Size)
	for i := 0; i < int(opts.Size); i++ {
		matrix[i] = make([]uint8, opts.Size)
	}

	ch := make(chan bool, opts.Channels)
	chunkSize := int(opts.Size / opts.Channels)
	for i := 0; i < int(opts.Size); i += chunkSize {
		go populate(matrix, i, opts, ch)
	}

	<-ch
	return matrix
}

func populate(m [][]uint8, iStart int, opts NoiseOptions, ch chan bool) {

	channelOffset := int(opts.Size / opts.Channels)
	iEnd := min(iStart+channelOffset, int(opts.Size))
	for i := iStart; i < iEnd; i++ {

		y := float64(i)/float64(opts.Size) + opts.YOffset/opts.Frequency
		for j := 0; j < int(opts.Size); j++ {

			x := float64(j)/float64(opts.Size) + opts.XOffset/opts.Frequency
			//sample := noise(x,y,opts.frequency) * 0.5 + 0.5
			sample := sum(x, y, opts)*0.5 + 0.5
			scale := uint8(sample * 255.0)
			m[i][j] = scale
		}
	}

	ch <- true
}

// Auxiliary vec2 type
type vec2 struct {
	x, y float64
}

func (v vec2) normalized() vec2 {
	h := math.Sqrt(v.x*v.x + v.y*v.y)
	return vec2{x: v.x / h, y: v.y / h}
}

// Hash mask
const hm_size = 255

var hm = [...]uint8{
	151, 160, 137, 91, 90, 15, 131, 13, 201, 95, 96, 53, 194, 233, 7, 225,
	140, 36, 103, 30, 69, 142, 8, 99, 37, 240, 21, 10, 23, 190, 6, 148,
	247, 120, 234, 75, 0, 26, 197, 62, 94, 252, 219, 203, 117, 35, 11, 32,
	57, 177, 33, 88, 237, 149, 56, 87, 174, 20, 125, 136, 171, 168, 68, 175,
	74, 165, 71, 134, 139, 48, 27, 166, 77, 146, 158, 231, 83, 111, 229, 122,
	60, 211, 133, 230, 220, 105, 92, 41, 55, 46, 245, 40, 244, 102, 143, 54,
	65, 25, 63, 161, 1, 216, 80, 73, 209, 76, 132, 187, 208, 89, 18, 169,
	200, 196, 135, 130, 116, 188, 159, 86, 164, 100, 109, 198, 173, 186, 3, 64,
	52, 217, 226, 250, 124, 123, 5, 202, 38, 147, 118, 126, 255, 82, 85, 212,
	207, 206, 59, 227, 47, 16, 58, 17, 182, 189, 28, 42, 223, 183, 170, 213,
	119, 248, 152, 2, 44, 154, 163, 70, 221, 153, 101, 155, 167, 43, 172, 9,
	129, 22, 39, 253, 19, 98, 108, 110, 79, 113, 224, 232, 178, 185, 112, 104,
	218, 246, 97, 228, 251, 34, 242, 193, 238, 210, 144, 12, 191, 179, 162, 241,
	81, 51, 145, 235, 249, 14, 239, 107, 49, 192, 214, 31, 181, 199, 106, 157,
	184, 84, 204, 176, 115, 121, 50, 45, 127, 4, 150, 254, 138, 236, 205, 93,
	222, 114, 67, 29, 24, 72, 243, 141, 128, 195, 78, 66, 215, 61, 156, 180,
}

// Adjacent gradients
const grad_size = 7

var grad = [...]vec2{
	vec2{x: 1.0, y: 0.0},
	vec2{x: -1.0, y: 0.0},
	vec2{x: 0.0, y: 1.0},
	vec2{x: 0.0, y: -1.0},
	vec2{x: 1.0, y: 1.0}.normalized(),
	vec2{x: -1.0, y: 1.0}.normalized(),
	vec2{x: 1.0, y: -1.0}.normalized(),
	vec2{x: -1.0, y: -1.0}.normalized(),
}

func noise(x, y, frequency float64) float64 {

	fx := x * frequency
	fy := y * frequency
	ix0 := int(fx)
	iy0 := int(fy)
	tx0 := fx - float64(ix0)
	ty0 := fy - float64(iy0)
	tx1 := tx0 - 1.0
	ty1 := ty0 - 1.0
	ix0 &= hm_size
	iy0 &= hm_size

	ix1 := (ix0 + 1) & hm_size
	iy1 := (iy0 + 1) & hm_size

	h0 := int(hm[ix0])
	h1 := int(hm[ix1])
	hm00 := (h0 + iy0) & hm_size
	hm10 := (h1 + iy0) & hm_size
	hm01 := (h0 + iy1) & hm_size
	hm11 := (h1 + iy1) & hm_size
	g00 := grad[hm[hm00]&grad_size]
	g10 := grad[hm[hm10]&grad_size]
	g01 := grad[hm[hm01]&grad_size]
	g11 := grad[hm[hm11]&grad_size]

	v00 := dot(g00, vec2{x: tx0, y: ty0})
	v10 := dot(g10, vec2{x: tx1, y: ty0})
	v01 := dot(g01, vec2{x: tx0, y: ty1})
	v11 := dot(g11, vec2{x: tx1, y: ty1})

	tx := smooth(tx0)
	ty := smooth(ty0)

	l1 := lerp(v00, v10, tx)
	l2 := lerp(v01, v11, tx)
	noise := lerp(l1, l2, ty) * math.Sqrt(2.0)

	return noise
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func smooth(t float64) float64 {
	return t * t * t * (t*(t*6.0-15.0) + 10.0)
}

func lerp(min, max, pos float64) float64 {

	diff := max - min
	value := min + float64(diff)*pos
	return value
}

func dot(v1, v2 vec2) float64 {
	return v1.x*v2.x + v1.y*v2.y
}

func sum(x float64, y float64, opts NoiseOptions) float64 {

	sum := noise(x, y, opts.Frequency)

	localFrequency := opts.Frequency
	amplitude := 1.0
	breadth := 1.0

	for i := 1; i < int(opts.Octaves); i++ {
		localFrequency *= opts.Lacunarity
		amplitude *= opts.Persistence
		breadth += amplitude
		sum += noise(x, y, localFrequency) * amplitude
	}

	return sum / breadth
}
