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
				for i := 0; i < b.N; i++ {
					b.StopTimer()

					sut := time.LoadLocation

					if tc == "locache" {
						sut = newCache().LoadLocation
					}

					b.StartTimer()

					for j := 0; j < numLookups; j++ {
						if _, err := sut("Europe/Kyiv"); err != nil {
							b.Fatal(err)
						}
					}
				}
			})
		}
	}
}
