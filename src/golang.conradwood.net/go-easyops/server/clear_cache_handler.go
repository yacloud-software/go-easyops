package server

import (
	"fmt"
	"golang.conradwood.net/go-easyops/cache"
	"net/http"
	"strings"
)

func clearCacheHandler(w http.ResponseWriter, req *http.Request) {
	cacheName := strings.TrimPrefix(req.URL.Path, "/internal/clearcache")
	cacheName = strings.TrimPrefix(cacheName, "/")
	caches, err := cache.Clear(cacheName)
	if err != nil {
		s := fmt.Sprintf("Cache \"%s\" could not be cleared: %s\n", cacheName, err)
		w.Write([]byte(s))
		return
	}
	w.Write([]byte("<html><body>"))
	for _, c := range caches {
		s := fmt.Sprintf("Cleared Cache \"%s\"</br>\n", c.Name())
		w.Write([]byte(s))
	}
	w.Write([]byte("</body></html>"))
}
