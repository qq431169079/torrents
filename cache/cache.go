package cache

import (
	"time"
	"errors"
	"strings"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/op/go-logging"
	"github.com/videofree/torrents/util"
)

const (
	DEFAULT              = time.Duration(0)
	FOREVER              = time.Duration(-1)
	CACHE_MIDDLEWARE_KEY = "gincontrib.cache"
)

var (
	pageCachePrefix = "page"
	errCacheMiss    = errors.New("cache: key not found")
	errNotStored    = errors.New("cache: not stored")
	errNotSupported = errors.New("cache: not supported")
	log             = logging.MustGetLogger("cache")
)

type CacheStore interface {
	Get(key string, value interface{}) error
	Set(key string, value interface{}, expire time.Duration) error
	Add(key string, value interface{}, expire time.Duration) error
	Replace(key string, data interface{}, expire time.Duration) error
	Delete(key string) error
	Increment(key string, data uint64) (uint64, error)
	Decrement(key string, data uint64) (uint64, error)
	Flush() error
}

type responseCache struct {
	Status int
	Header http.Header
	Data   []byte
}

type cachedWriter struct {
	gin.ResponseWriter
	status  int
	written bool
	store   CacheStore
	expire  time.Duration
	key     string
}

func cacheKey(prefix string, u string) string {
	u = strings.Trim(u, "/")
	dotted := []string{"/", "=", "?", "&"}
	for _, dottedChar := range dotted {
		u = strings.Replace(u, dottedChar, ".", -1)
	}
	return prefix + "." + util.ToFileName(u)
}

func newCachedWriter(store CacheStore, expire time.Duration, writer gin.ResponseWriter, key string) *cachedWriter {
	return &cachedWriter{writer, 0, false, store, expire, key}
}

func (w *cachedWriter) WriteHeader(code int) {
	w.status = code
	w.written = true
	w.ResponseWriter.WriteHeader(code)
}

func (w *cachedWriter) Status() int {
	return w.status
}

func (w *cachedWriter) Written() bool {
	return w.written
}

func (w *cachedWriter) Write(data []byte) (int, error) {
	ret, err := w.ResponseWriter.Write(data)
	if err == nil {
		//cache response
		store := w.store
		val := responseCache{
			w.status,
			w.Header(),
			data,
		}
		err = store.Set(w.key, val, w.expire)
		if err != nil {
			log.Error(err)
		}
	}
	return ret, err
}

// Cache Middleware
func Cache(store CacheStore, expire time.Duration) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var cache responseCache
		key := cacheKey(pageCachePrefix, ctx.Request.URL.RequestURI())
		if err := store.Get(key, &cache); err == nil {
			for k, vals := range cache.Header {
				for _, v := range vals {
					ctx.Writer.Header().Add(k, v)
				}
			}
			ctx.AbortWithStatus(cache.Status)
			ctx.Writer.Write(cache.Data)
		} else {
			// replace writer
			writer := ctx.Writer
			ctx.Writer = newCachedWriter(store, expire, ctx.Writer, key)
			ctx.Next()
			ctx.Writer = writer
		}
	}
}
