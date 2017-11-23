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

import (
	"fmt"
	"testing"
	"time"
)

func PerformanceTest(t *testing.T) {

	opts := NoiseOptions{
		Size:        4096,
		Octaves:     2,
		Frequency:   16.0,
		Lacunarity:  2.5,
		Persistence: 0.75,
		Channels:    1,
	}

	var elapsed time.Duration
	elapsed = timeTest(1, opts)
	fmt.Printf("Time to run with %d channel(s) = %s\n", 1, elapsed)
	elapsed = timeTest(2, opts)
	fmt.Printf("Time to run with %d channel(s) = %s\n", 2, elapsed)
	elapsed = timeTest(4, opts)
	fmt.Printf("Time to run with %d channel(s) = %s\n", 4, elapsed)
	elapsed = timeTest(8, opts)
	fmt.Printf("Time to run with %d channel(s) = %s\n", 8, elapsed)
	elapsed = timeTest(16, opts)
	fmt.Printf("Time to run with %d channel(s) = %s\n", 16, elapsed)
	elapsed = timeTest(32, opts)
	fmt.Printf("Time to run with %d channel(s) = %s\n", 32, elapsed)
	elapsed = timeTest(64, opts)
	fmt.Printf("Time to run with %d channel(s) = %s\n", 64, elapsed)
}

func timeTest(channels int, opts NoiseOptions) time.Duration {

	opts.Channels = uint(channels)
	start := time.Now()
	Build(opts)
	end := time.Now()
	elapsed := end.Sub(start)

	return elapsed
}
