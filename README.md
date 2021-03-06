#Cache
[![wercker status](https://app.wercker.com/status/3ddcdbd94dda6114613d69a4360f3f1c/s "wercker status")](https://app.wercker.com/project/bykey/3ddcdbd94dda6114613d69a4360f3f1c)
[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://godoc.org/github.com/neko-contrib/cache)
[![GoCover](http://gocover.io/_badge/github.com/neko-contrib/cache)](http://gocover.io/github.com/neko-contrib/cache)

[Neko](https://github.com/rocwong/neko) handler for cache management.

## Usage
~~~go
package main

import (
  "time"
  "github.com/rocwong/neko"
  nc "github.com/neko-contrib/cache"
)

func main() {
  app := neko.Classic()
  app.Use(nc.Generate(nc.Options{}))

  m.GET("/", func(ctx *neko.Context) {
    cache := ctx.MustGet(nc.MemoryStore).(nc.Cache)
    cache.Set("foo", "bar", 10 * time.Second)
  })

  m.GET("/get", func(ctx *neko.Context) {
    cache := ctx.MustGet(nc.MemoryStore).(nc.Cache)
    v, found := cache.Get("foo")
    ctx.Text(v.(string))
  })

  app.Run(":3000")
}
~~~

## Options
~~~go
cache.Options {
  // Store cache store. Default is 'MemoryStore'
  Store string
  // Config stores configuration.
  Config string
  // Interval GC interval time in seconds. Default is 60.
  Interval int
}
~~~

## Stores

#### Memory
~~~go
app.Use(cache.Generate(cache.Options{}))
~~~
