noise
=====

A 2D texture generator based on Perlin noise written in Go.

It comes with both CLI and GUI built in, but can also be used directly as a
library. In this README you will find examples of how to use noise with each
of these approaches.

## Installation

```bash
git clone git@github.com:htaunay/noise.git
cd $GOPATH/github.com/htaunay/noise

# installs necessary dependencies as well
make install
```

## Parameters

The noise package's interface is basically a single function that receives

#### Image Size

Side of square texture length, in pixels

#### Frequency

Period at which data is sampled. Think as this param as defining how many
variations (i.e tiles) of the noise will appear in a side of the texture,
e.g. if the **frequency** is set to two, you would have 2x2 samples in the
output.

This parameters accepts floats, and therefore if you give it a non-integer,
the trailing sampled data in each dimension will be capped accordingly.

#### Octaves

Number of times that different noise functions - varying between each other
based on the **lacunarity** and **persistence** values - will be added toguether.

If set to the minimum of 1, there will be a single layer, and therefore only
the **frequency** and **offsets** will infleunce the generation of the texture.

#### Lacunarity

Multiplier that determines how quickly the frequency increases for each
successive octave.

#### Persistence

Determines how much influence should each successive octave have, quantitatively
over the previous one.

Influence starts at 1.0 for the first layer, add is multiplied by the
persistence on each extra octave.

#### Offsets

Horizontal and vertical ratios for which the sample generation will be shifted.

The amount is based of the textures size, e.g. an offset of 1.0 means that a
given axis will start exactly after the last pixel of the original sample.

The x-axis shifts to the right, while the y-axis shifts downward.

## Package

```go
import "github.com/htaunay/noise"

opts := noise.NoiseOptions{

	// example values
	Size:        1024,
	Octaves:     2,
	Frequency:   16,
	Lacunarity:  2.5,
	Persistence: 0.75,
	XOffset:     0,
	YOffset:     0,
	Channels:    1,
}

// returns [][]uint8 with values varying from 0-255
matrix := noise.Build(ops)
```

## CLI

If installation completed correctly, the CLI bin should be located in the
Go's bin folder.

```go
// In case you haven't added it to your $PATH
cd $GOPATH/bin

// Check out the options available
noise-cli --help
```

## GUI

If installation completed correctly, the GUI bin should be located in the
Go's bin folder.

```go
// In case you haven't added it to your $PATH
cd $GOPATH/bin

// Should open a window
noise-gui
```

### Commands

#### X and Y Offsets

Controlled by the arrow keys (←,↑,→,↓)

#### Frequency

**Q**-key increases frequency, and while holding **Shift** it decreases

#### Octaves

**O**-key increases the count of octaves, and while holding **Shift** it decreases

#### Lacunarity

**L**-key increases lacunarity, and while holding **Shift** it decreases

#### Persistence

**P**-key increases persistence, and while holding **Shift** it decreases

#### Filter

**F**-key toggles the default filter

## References

* [Understanding Perlin Noise](http://flafla2.github.io/2014/08/09/perlinnoise.html)
* [libnoise](http://libnoise.sourceforge.net/index.html)
* [Noise, being a pseudorandom artists](http://catlikecoding.com/unity/tutorials/noise/)
