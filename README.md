# modtimefs

## Description

Package modtimefs wraps fs.FS with fake ModTime.

embed.FS always returns a zero value of time.Time for ModTime(), so the Last-Modified is not added to the response header when embed.FS is used with http.FileServer and http.FS.

This package can be avoid this with user specific ModTime.

The file in original fs.FS must implement io.Seeker to use with http.FileServer.

## Example

```go
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
```
