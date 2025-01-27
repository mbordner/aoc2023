package files

import (
	"fmt"
	"github.com/mbordner/memfs"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test_File_Functions(t *testing.T) {
	origFS := FS
	FS = &memFS{fs: memfs.New()}
	defer func() {
		FS = origFS

	}()

	data := `line 1
line 2
line 3
line 4
`

	err := WriteContent("/test/test", []byte(data))
	assert.Nil(t, err)

	lines, err := GetLines("/test/test")
	assert.Nil(t, err)
	assert.Equal(t, 5, len(lines))

	err = WriteContent("/test/test", []byte(data))
	assert.Nil(t, err)
	assert.Equal(t, 5, len(lines))

	err = WriteContent("/test/test2", []byte(data))
	assert.Nil(t, err)

	err = WriteContent("/test/test3", []byte(data))
	assert.Nil(t, err)

	var tmpZFName string
	tmpZFName, err = CreateTempFile("tmp")
	assert.Nil(t, err)
	assert.NotNil(t, tmpZFName)

	err = Zip(`/test`, tmpZFName)
	assert.Nil(t, err)

	var tmpDirName string
	tmpDirName, err = CreateTempDir(`/tmp`, `dir`)
	assert.Nil(t, err)
	assert.NotNil(t, tmpDirName)

	err = Unzip(tmpZFName, tmpDirName)
	assert.Nil(t, err)

	lines, err = GetLines(fmt.Sprintf("%s/test", tmpDirName))
	assert.Nil(t, err)
	assert.Nil(t, err)
	assert.Equal(t, 5, len(lines))

}

// memFS implements FileSystem using an in memory filesystem
type memFS struct{ fs *memfs.FS }

func (mfs *memFS) Open(name string) (File, error)   { return mfs.fs.Open(name) }
func (mfs *memFS) Create(name string) (File, error) { return mfs.fs.Create(name) }
func (mfs *memFS) OpenFile(name string, flag int, perm os.FileMode) (File, error) {
	return mfs.fs.OpenFile(name, flag, perm)
}
func (mfs *memFS) Stat(name string) (os.FileInfo, error) { return mfs.fs.Stat(name) }
func (mfs *memFS) Remove(name string) error              { return mfs.fs.Remove(name) }
func (mfs *memFS) CreateTemp(dir, pattern string) (File, error) {
	return mfs.fs.CreateTemp(dir, pattern)
}
func (mfs *memFS) MkdirAll(path string, perms os.FileMode) error { return mfs.fs.MkdirAll(path, perms) }
func (mfs *memFS) RemoveAll(path string) error                   { return mfs.fs.RemoveAll(path) }
func (mfs *memFS) ReadDir(name string) ([]os.DirEntry, error)    { return mfs.fs.ReadDir(name) }
func (mfs *memFS) MkdirTemp(dir, pattern string) (name string, err error) {
	return mfs.fs.MkdirTemp(dir, pattern)
}
func (mfs *memFS) TempDir() string { return mfs.fs.TempDir() }
