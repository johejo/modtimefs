// Package modtimefs wraps fs.FS with fake ModTime.
// embed.FS always returns a zero value of time.Time for ModTime(), so the Last-Modified is not added to the response header when embed.FS is used with http.FileServer and http.FS.
// This package can be avoid this with user specific ModTime.
// The file in original fs.FS must implement io.Seeker to use with http.FileServer.
package modtimefs

import (
	"errors"
	"io"
	"io/fs"
	"time"
)

var errNotSeeker = errors.New("file does not implement io.Seeker")

// NewFn takes a original fs.FS and a function for spoofing ModTime, and returns a wrapped fs.FS.
func NewFn(fsys fs.FS, modTimeFn func() time.Time) fs.FS {
	return modTimeFS{FS: fsys, modTimeFn: modTimeFn}
}

// NewFn takes a original fs.FS and a static ModTime, and returns a wrapped fs.FS.
func New(fsys fs.FS, modTime time.Time) fs.FS {
	return NewFn(fsys, func() time.Time { return modTime })
}

type modTimeFS struct {
	fs.FS
	modTimeFn func() time.Time
}

func (fsys modTimeFS) Open(name string) (fs.File, error) {
	f, err := fsys.FS.Open(name)
	if err != nil {
		return nil, err
	}
	sf, ok := f.(seekerFile)
	if !ok {
		return nil, errNotSeeker
	}
	return modTimeFile{seekerFile: sf, modTimeFn: fsys.modTimeFn}, nil
}

type seekerFile interface {
	fs.File
	io.Seeker
}

type modTimeFile struct {
	seekerFile
	modTimeFn func() time.Time
}

func (f modTimeFile) Stat() (fs.FileInfo, error) {
	fi, err := f.seekerFile.Stat()
	if err != nil {
		return nil, err
	}
	return modTimeFileInfo{FileInfo: fi, modTimeFn: f.modTimeFn}, nil
}

type modTimeFileInfo struct {
	fs.FileInfo
	modTimeFn func() time.Time
}

func (fi modTimeFileInfo) ModTime() time.Time {
	return fi.modTimeFn()
}
