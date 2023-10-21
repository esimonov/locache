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

	for i := 0; i < 100; i++ {
		for _, locName := range testLocNames {
			wg.Add(1)

			go func(name string) {
				defer wg.Done()

				if _, err := locache.LoadLocation(name); err != nil {
					panic(err)
				}
			}(locName)
		}
	}

	wg.Wait()
}

func TestLoadLocation_ReturnArgumentsAreEquivalentToNativeImplementation_Property(t *testing.T) {
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
	defer closefd(fd)

	const (
		zecheader = 0x06054b50
		zcheader  = 0x02014b50
		ztailsize = 22
	)

	buf := make([]byte, ztailsize)

	preadn(fd, buf, -ztailsize)

	if get4(buf) != zecheader {
		panic("corrupt zip file " + zipfile)
	}

	n := get2(buf[10:])
	size := get4(buf[12:])
	off := get4(buf[16:])

	buf = make([]byte, size)

	preadn(fd, buf, off)

	ianaNames := make([]string, 0)

	for i := 0; i < n; i++ {
		if get4(buf) != zcheader {
			break
		}

		namelen := get2(buf[28:])
		xlen := get2(buf[30:])
		fclen := get2(buf[32:])
		zname := buf[46 : 46+namelen]
		buf = buf[46+namelen+xlen+fclen:]
		ianaNames = append(ianaNames, string(zname))
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

func closefd(fd uintptr) {
	syscall.Close(int(fd))
}

func preadn(fd uintptr, buf []byte, off int) {
	whence := 0
	if off < 0 {
		whence = 2
	}

	syscall.Seek(int(fd), int64(off), whence)
	syscall.Read(int(fd), buf)
}

func get4(b []byte) int {
	if len(b) < 4 {
		return 0
	}
	return int(b[0]) | int(b[1])<<8 | int(b[2])<<16 | int(b[3])<<24
}

func get2(b []byte) int {
	if len(b) < 2 {
		return 0
	}
	return int(b[0]) | int(b[1])<<8
}
