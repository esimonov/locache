# locache

A very simple in-memory cache for time.Locations, written in pure Go.

## Why?

It is [known](https://pkg.go.dev/time#LoadLocation) that `time.LoadLocation` is reading the locations information from disk. Furthermore, previously retrieved locations are not being cached: the same location is looked up as many times as `time.LoadLocation` is called. If your application needs to open a bunch of locations very limited number of times, this behaviour would be desirable. On the other hand, when it is constantly looking up the same locations, this may seem wasteful. `locache` allows to perform disk lookups only once per location, and then reuse the results.

The mitigation becomes more noticeable as the number of repeated lookups grows:

```
goos: darwin
goarch: arm64
pkg: github.com/esimonov/locache
PASS
benchmark                                                 iter         time/iter    bytes alloc            allocs
---------                                                 ----         ---------    -----------            ------
Benchmark_LoadLocation/time_10_RepeatedLookups-10         6994      181.24 μs/op     51136 B/op     220 allocs/op
Benchmark_LoadLocation/time_100_RepeatedLookups-10         705     1638.50 μs/op    511217 B/op    2200 allocs/op
Benchmark_LoadLocation/time_1000_RepeatedLookups-10         75    16422.01 μs/op   5112026 B/op   22000 allocs/op
Benchmark_LoadLocation/locache_10_RepeatedLookups-10     58425       20.53 μs/op      5328 B/op      23 allocs/op
Benchmark_LoadLocation/locache_100_RepeatedLookups-10    54592       23.36 μs/op      5328 B/op      23 allocs/op
Benchmark_LoadLocation/locache_1000_RepeatedLookups-10   31910       35.70 μs/op      5328 B/op      23 allocs/op
ok      github.com/esimonov/locache     16.247s
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
