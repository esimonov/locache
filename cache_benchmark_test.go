package locache

import (
	"fmt"
	"testing"
	"time"
)

func Benchmark_LoadLocation(b *testing.B) {
	b.ResetTimer()

	for i, tc := range []string{"Native", "Locache"} {
		for _, numLookups := range []int{10, 100, 1000} {
			b.Run(fmt.Sprintf("%s_%dRepeatedLookups", tc, numLookups), func(b *testing.B) {
				for j := 0; j < b.N; j++ {
					b.StopTimer()

					sut := time.LoadLocation

					if i == 1 {
						sut = newCache().LoadLocation
					}

					b.StartTimer()

					for k := 0; k < numLookups; k++ {
						if _, err := sut("Europe/Kyiv"); err != nil {
							b.Fatal(err)
						}
					}
				}
			})
		}
	}
}
