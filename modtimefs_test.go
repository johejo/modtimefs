package modtimefs_test

import (
	"embed"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/johejo/modtimefs"
)

//go:embed testdata/*
var testdata embed.FS

func Example() {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/testdata/hello.txt", nil)

	fsys := modtimefs.New(testdata, time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC))
	http.FileServer(http.FS(fsys)).ServeHTTP(rec, req)
	resp := rec.Result()

	// Output: Sun, 01 Jan 2023 00:00:00 GMT
	fmt.Println(resp.Header.Get("Last-Modified"))
}
