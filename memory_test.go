package cache

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func Test_MemoryCacher(t *testing.T) {
	Convey("Memory Store", t, func() {
		testStore(Options{
			Store:    MemoryStore,
			Interval: 4,
		})
	})
}
