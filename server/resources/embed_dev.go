//go:build dev

package resources

import (
	"io/fs"
	"os"
	"path"
	"runtime"
)

var Content fs.FS = os.DirFS("server/resources")

func init() {
	_, file, _, _ := runtime.Caller(0)
	dir := path.Dir(file)
	Content = os.DirFS(dir)
}
