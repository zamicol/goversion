// Package goversion uses go modules' method of versioning for the purpose of tracking a current project's version.
package goversion

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
)

// DefaultHash is the default function for hashing.
var DefaultHash = Hash1

// Hash is type for hashing.
// Taken from cmd/go/internal/dirhash/hash.go
type Hash func(files []string, open func(string) (io.ReadCloser, error)) (string, error)

// Hash1 is the go mod hashing function.
// Taken from cmd/go/internal/dirhash/hash.go
func Hash1(files []string, open func(string) (io.ReadCloser, error)) (string, error) {
	h := sha256.New()
	files = append([]string(nil), files...)
	sort.Strings(files)
	for _, file := range files {
		if strings.Contains(file, "\n") {
			return "", errors.New("filenames with newlines are not supported")
		}
		r, err := open(file)
		if err != nil {
			return "", err
		}
		hf := sha256.New()
		_, err = io.Copy(hf, r)
		r.Close()
		if err != nil {
			return "", err
		}
		fmt.Fprintf(h, "%x  %s\n", hf.Sum(nil), file)
	}
	return "h1:" + base64.StdEncoding.EncodeToString(h.Sum(nil)), nil
}

// HashZip hashes a zip.
// Taken from cmd/go/internal/dirhash/hash.go
func HashZip(zipfile string, hash Hash) (string, error) {
	z, err := zip.OpenReader(zipfile)
	if err != nil {
		return "", err
	}
	defer z.Close()
	var files []string
	zfiles := make(map[string]*zip.File)
	for _, file := range z.File {
		files = append(files, file.Name)
		zfiles[file.Name] = file
	}
	zipOpen := func(name string) (io.ReadCloser, error) {
		f := zfiles[name]
		if f == nil {
			return nil, fmt.Errorf("file %q not found in zip", name) // should never happen
		}
		return f.Open()
	}
	return hash(files, zipOpen)
}

var (
	slashSlash = []byte("//")
	moduleStr  = []byte("module")
)

// ModulePath returns the module path from the gomod file text.
// Must be given mod file contents.
// Taken from https://github.com/golang/go/blob/master/src/cmd/go/internal/modfile/read.go
// See https://github.com/golang/go/blob/master/src/cmd/go/internal/modfetch/codehost/git.go#L466
// For function Go uses to read mod.go file.
// If it cannot find a module path, it returns an empty string.
// It is tolerant of unrelated problems in the go.mod file.
func ModulePath(mod []byte) string {
	for len(mod) > 0 {
		line := mod
		mod = nil
		if i := bytes.IndexByte(line, '\n'); i >= 0 {
			line, mod = line[:i], line[i+1:]
		}
		if i := bytes.Index(line, slashSlash); i >= 0 {
			line = line[:i]
		}
		line = bytes.TrimSpace(line)
		if !bytes.HasPrefix(line, moduleStr) {
			continue
		}
		line = line[len(moduleStr):]
		n := len(line)
		line = bytes.TrimSpace(line)
		if len(line) == n || len(line) == 0 {
			continue
		}

		if line[0] == '"' || line[0] == '`' {
			p, err := strconv.Unquote(string(line))
			if err != nil {
				return "" // malformed quoted string or multiline module path
			}
			return p
		}

		return string(line)
	}
	return "" // missing module path
}
