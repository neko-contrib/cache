package cache

import (
	"github.com/rocwong/neko"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func performRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func Test_Cacher(t *testing.T) {
	Convey("Use cache middleware", t, func() {
		m := neko.New()
		m.Use(Generate())
		m.GET("/", func(ctx *neko.Context) { ctx.Text("neko cache") })

		w := performRequest(m, "GET", "/")
		So(w.Body.String(), ShouldEqual, "neko cache")
	})
}

func testStore(opt Options) {
	m := neko.New()
	m.Use(Generate(opt))
	m.GET("/set", func(ctx *neko.Context) {
		cache := ctx.MustGet(MemoryStore).(Cache)
		cache.Set("myvalue", "memory cache")
		cache.Set("myexpire", "bar", 2*time.Second)
		cache.Set("gc", "bar", 3*time.Second)
		cache.Set("flush", "memory cache")

		cache.Set("int", 0, 0)
		cache.Set("int32", int32(0), 0)
		cache.Set("int64", int64(0), 0)
		cache.Set("uint", uint(0), 0)
		cache.Set("uint32", uint32(0), 0)
		cache.Set("uint64", uint64(0), 0)
	})

	Convey("Basic Operation", func() {
		performRequest(m, "GET", "/set")
		m.GET("/get", func(ctx *neko.Context) {
			cache := ctx.MustGet(MemoryStore).(Cache)
			v, _ := cache.Get("myvalue")
			So(v.(string), ShouldEqual, "memory cache")

			v, _ = cache.Get("nocache")
			So(v, ShouldBeNil)

			So(cache.IsExist("myvalue"), ShouldBeTrue)

			cache.Delete("myvalue")
			So(cache.IsExist("myvalue"), ShouldBeFalse)
		})

		m.GET("/expire", func(ctx *neko.Context) {
			cache := ctx.MustGet(MemoryStore).(Cache)
			v, _ := cache.Get("myexpire")
			So(v, ShouldBeNil)
			So(cache.IsExist("gc"), ShouldBeTrue)

		})
		m.GET("/flush", func(ctx *neko.Context) {
			cache := ctx.MustGet(MemoryStore).(Cache)
			cache.Flush()
			So(cache.IsExist("gc"), ShouldBeFalse)
			So(cache.IsExist("flush"), ShouldBeFalse)
		})

		performRequest(m, "GET", "/get")

		time.Sleep(2 * time.Second)
		performRequest(m, "GET", "/expire")

		time.Sleep(2 * time.Second)
		performRequest(m, "GET", "/flush")
	})

	Convey("Increment Operation", func() {
		performRequest(m, "GET", "/set")
		m.GET("/increment", func(ctx *neko.Context) {
			cache := ctx.MustGet(MemoryStore).(Cache)

			cache.Increment("int")
			v, _ := cache.Get("int")
			So(v.(int), ShouldEqual, 1)

			cache.Increment("int", 20)
			v, _ = cache.Get("int")
			So(v.(int), ShouldEqual, 21)

			So(cache.Increment("myvalue").Error(), ShouldEqual, "the item is not an integer")
			So(cache.Increment("noval").Error(), ShouldEqual, "Item noval not found")

		})
		performRequest(m, "GET", "/increment")
	})

	Convey("Decrement Operation", func() {
		performRequest(m, "GET", "/set")
		m.GET("/decrement", func(ctx *neko.Context) {

			cache := ctx.MustGet(MemoryStore).(Cache)

			cache.Decrement("int")
			v, _ := cache.Get("int")
			So(v.(int), ShouldEqual, -1)

			cache.Decrement("int", 4)
			v, _ = cache.Get("int")
			So(v.(int), ShouldEqual, -5)

			So(cache.Decrement("myvalue").Error(), ShouldEqual, "the item is not an integer")
			So(cache.Decrement("noval").Error(), ShouldEqual, "Item noval not found")

			So(cache.Decrement("int32"), ShouldBeNil)
			So(cache.Decrement("int64"), ShouldBeNil)
			So(cache.Decrement("uint"), ShouldNotBeNil)
			So(cache.Decrement("uint32"), ShouldNotBeNil)
			So(cache.Decrement("uint64"), ShouldNotBeNil)

			So(cache.Increment("int32"), ShouldBeNil)
			So(cache.Increment("int64"), ShouldBeNil)
			So(cache.Increment("uint"), ShouldBeNil)
			So(cache.Increment("uint32"), ShouldBeNil)
			So(cache.Increment("uint64"), ShouldBeNil)

			So(cache.Decrement("uint"), ShouldBeNil)
			So(cache.Decrement("uint32"), ShouldBeNil)
			So(cache.Decrement("uint64"), ShouldBeNil)

		})
		performRequest(m, "GET", "/decrement")
	})

}
