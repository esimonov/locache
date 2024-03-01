package locache_test

import (
	"go/build"
	"math/rand"
	"path"
	"reflect"
	"strconv"
	"sync"
	"syscall"
	"testing"
	"testing/quick"
	"time"

	"github.com/esimonov/locache"
)

var testLocNames []string

func init() {
	zoneinfoLoc := path.Join(build.Default.GOROOT, "/lib/time/zoneinfo.zip")

	testLocNames = loadIANATzNamesFromZip(zoneinfoLoc)

	if len(testLocNames) == 0 {
		panic("empty location names slice")
	}
}

func TestLoadLocation_SafeForConcurrentAccess(t *testing.T) {
	wg := new(sync.WaitGroup)

	for range 100 {
		for _, locName := range testLocNames {
			wg.Add(1)

			go func() {
				defer wg.Done()

				if _, err := locache.LoadLocation(locName); err != nil {
					panic(err)
				}
			}()
		}
	}

	wg.Wait()
}

func TestLoadLocation_ReturnArgumentsAreEquivalentToStandardLibraryImplementation_Property(t *testing.T) {
	f := func(n int64) bool {
		locName := "unknown" + strconv.Itoa(int(n))

		n %= int64(len(testLocNames))

		if n%2 == 0 {
			locName = testLocNames[(n^(n>>63))-(n>>63)]
		}

		loc1, err1 := time.LoadLocation(locName)

		loc2, err2 := locache.LoadLocation(locName)

		return reflect.DeepEqual(loc1, loc2) && (err1 == nil) == (err2 == nil)
	}

	cfg := &quick.Config{
		MaxCount: 10000,
		Rand:     rand.New(rand.NewSource(time.Now().UnixNano())),
	}

	if err := quick.Check(f, cfg); err != nil {
		t.Fatal(err)
	}
}

func loadIANATzNamesFromZip(zipfile string) []string {
	fd, err := open(zipfile)
	if err != nil {
		panic(err)
	}
	defer syscall.Close(int(fd))

	const ztailsize = 22

	buf := make([]byte, ztailsize)

	preadn(fd, buf, -ztailsize)

	getn(4, buf)

	n := getn(2, buf[10:])
	size := getn(4, buf[12:])
	off := getn(4, buf[16:])

	buf = make([]byte, size)

	preadn(fd, buf, off)

	ianaNames := make([]string, n)

	for i := range n {
		getn(4, buf)

		namelen := getn(2, buf[28:])
		xlen := getn(2, buf[30:])
		fclen := getn(2, buf[32:])

		ianaNames[i] = string(buf[46 : 46+namelen])

		buf = buf[46+namelen+xlen+fclen:]
	}
	return ianaNames
}

func open(name string) (uintptr, error) {
	fd, err := syscall.Open(name, syscall.O_RDONLY, 0)
	if err != nil {
		return 0, err
	}
	return uintptr(fd), nil
}

func preadn(fd uintptr, buf []byte, off int) {
	whence := 0
	if off < 0 {
		whence = 2
	}

	syscall.Seek(int(fd), int64(off), whence)
	syscall.Read(int(fd), buf)
}

func getn(n int, b []byte) (result int) {
	if len(b) < n {
		return 0
	}

	for i := range n {
		result |= int(b[i]) << (i * 8)
	}
	return result
}
