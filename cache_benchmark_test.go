package locache

import (
	"fmt"
	"testing"
	"time"
)

func Benchmark_LoadLocation(b *testing.B) {
	for _, tc := range [2]string{"time", "locache"} {
		for _, numLookups := range [3]int{10, 100, 1000} {
			b.Run(fmt.Sprintf("%s_%d_RepeatedLookups", tc, numLookups), func(b *testing.B) {
				for range b.N {
					b.StopTimer()

					sut := time.LoadLocation

					if tc == "locache" {
						sut = newCache().LoadLocation
					}

					b.StartTimer()

					for range numLookups {
						if _, err := sut("Europe/Kyiv"); err != nil {
							b.Fatal(err)
						}
					}
				}
			})
		}
	}
}
