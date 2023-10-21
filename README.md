# locache

A very simple in-memory cache for time.Locations, written in pure Go.

## Why?

It is [known](https://pkg.go.dev/time#LoadLocation) that `time.LoadLocation` is reading the locations information from disk. Furthermore, previously retrieved locations are not being cached: the same location is looked up as many times as `time.LoadLocation` is called. If your application needs to open a bunch of locations very limited number of times, this behaviour would be desirable. On the other hand, when it is constantly looking up the same locations (e.g, in a HTTP handler), this may seem wasteful. `locache` allows to perform disk lookups only once per location, and then reuse the results.

The extent of this mitigation grows linearly with the number of repeated lookups:

```
goos: darwin
goarch: arm64
pkg: github.com/esimonov/locache
PASS
benchmark                                                iter         time/iter    bytes alloc            allocs
---------                                                ----         ---------    -----------            ------
Benchmark_LoadLocation/Native_10RepeatedLookups-10       6628      176.48 μs/op     51136 B/op     220 allocs/op
Benchmark_LoadLocation/Native_100RepeatedLookups-10       721     1717.07 μs/op    511216 B/op    2200 allocs/op
Benchmark_LoadLocation/Native_1000RepeatedLookups-10       74    16707.26 μs/op   5112024 B/op   22000 allocs/op
Benchmark_LoadLocation/Locache_10RepeatedLookups-10     55234       21.50 μs/op      5328 B/op      23 allocs/op
Benchmark_LoadLocation/Locache_100RepeatedLookups-10    54465       23.67 μs/op      5328 B/op      23 allocs/op
Benchmark_LoadLocation/Locache_1000RepeatedLookups-10   31730       37.09 μs/op      5328 B/op      23 allocs/op
ok      github.com/esimonov/locache     16.544s
```

## Usage

`locache.LoadLocation` has the same signature as `time.LoadLocation`, so they can be used interchangeably.

```go
package main

import (
	"fmt"
	"github.com/esimonov/locache"
)

func main() {
	location, err := locache.LoadLocation("America/Los_Angeles")
	if err != nil {
		panic(err)
	}

	fmt.Println(time.Date(2018, 8, 30, 12, 0, 0, 0, time.UTC).In(location))
}
```

would produce the same result as:

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	location, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		panic(err)
	}

	fmt.Println(time.Date(2018, 8, 30, 12, 0, 0, 0, time.UTC).In(location))
}
```
