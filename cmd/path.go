package cmd

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"golang.org/x/mod/modfile"
)

func home() string {
	dir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	home := path.Join(dir, ".wmsctl")
	if _, err := os.Stat(home); os.IsNotExist(err) {
		if err := os.MkdirAll(home, 0o700); err != nil {
			panic(err)
		}
	}
	return home
}

func withHomeDir(dir string) string {
	home := path.Join(home(), dir)
	if _, err := os.Stat(home); os.IsNotExist(err) {
		if err := os.MkdirAll(home, 0o700); err != nil {
			panic(err)
		}
	}
	return home
}

func copyFile(src, dst string, replaces []string) error {
	var err error
	srcinfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	buf, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	var old string
	for i, next := range replaces {
		if i%2 == 0 {
			old = next
			continue
		}
		buf = bytes.ReplaceAll(buf, []byte(old), []byte(next))
	}
	return os.WriteFile(dst, buf, srcinfo.Mode())
}

func copyDir(src, dst string, replaces, ignores []string) error {
	var err error
	var fds []os.DirEntry
	var srcinfo os.FileInfo

	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}

	if err = os.MkdirAll(dst, srcinfo.Mode()); err != nil {
		return err
	}

	if fds, err = os.ReadDir(src); err != nil {
		return err
	}
	for _, fd := range fds {
		if hasSets(fd.Name(), ignores) {
			continue
		}
		srcfp := path.Join(src, fd.Name())
		dstfp := path.Join(dst, fd.Name())
		var e error
		if fd.IsDir() {
			e = copyDir(srcfp, dstfp, replaces, ignores)
		} else {
			e = copyFile(srcfp, dstfp, replaces)
		}
		if e != nil {
			return e
		}
	}
	return nil
}

func hasSets(name string, sets []string) bool {
	for _, ig := range sets {
		if ig == name {
			return true
		}
	}
	return false
}

func modulePath(filename string) (string, error) {
	modBytes, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return modfile.ModulePath(modBytes), nil
}

func gpath(home, url string) string {
	start := strings.LastIndex(url, "/")
	end := strings.LastIndex(url, ".git")
	if end == -1 {
		end = len(url)
	}
	branch := "@main"
	return path.Join(home, url[start+1:end]+branch)
}

func tree(path string, dir string) {
	_ = filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if err == nil && info != nil && !info.IsDir() {
			fmt.Printf("%s (%v bytes)\n", strings.Replace(path, dir+"/", "", -1), info.Size())
		}
		return nil
	})
}
