// Code generated for package packed by go-bindata DO NOT EDIT. (@generated)
// sources:
// manifest/config/config.yaml
package packed

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

// Name return file name
func (fi bindataFileInfo) Name() string {
	return fi.name
}

// Size return file size
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}

// Mode return file mode
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}

// Mode return file modify time
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}

// IsDir return file whether a directory
func (fi bindataFileInfo) IsDir() bool {
	return fi.mode&os.ModeDir != 0
}

// Sys return file is sys mode
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _manifestConfigConfigYaml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x84\x52\x4d\x73\xd3\x30\x10\xbd\xfb\x57\xec\xa4\x67\x27\x71\x42\xd2\x44\xbf\x80\x3b\xdc\x19\xd5\x5e\x3b\x6a\x65\x4b\x95\xe4\xd2\x70\xee\x4c\xcb\x94\xa1\xcc\x30\x90\x0e\x0c\x37\xa0\x5c\xfa\x71\x60\x0a\xd3\x42\xfb\x67\xea\xa4\x3d\xf1\x17\x18\xd9\x4e\xdc\x94\xb4\x28\x87\x48\x6f\xf5\xd6\xfb\xde\x93\x92\x3e\x71\x00\x12\x11\xa0\xfd\x07\x18\x08\x6d\x08\x70\xe1\x53\x6e\xb7\x39\x26\x85\x32\x04\xbc\x4e\xb3\xe9\x39\x03\x63\xa4\xbd\x99\xca\x80\x9a\x39\x4e\xb3\x9e\xff\xe6\x19\x1d\x2f\x3f\xc7\xb6\x3f\x04\xb8\x92\x46\x0e\x40\x10\xf0\xff\x12\xbb\x8b\x88\x31\x5d\x15\xea\x61\x6a\xcf\x0e\x79\x97\x89\xc6\x0f\xa6\x32\x75\xc1\x77\xc1\xeb\xb7\xea\x5e\xb7\x57\xf7\xea\x7d\xe2\xb5\xda\xcb\xfd\x05\x78\xeb\x1e\xbc\x5d\xe2\xa9\x46\x95\xd0\x18\x09\x28\x21\xac\x59\x92\x6a\xfd\x5c\xa8\x80\xc0\x63\x33\x58\xf5\x3a\xcd\x76\xab\x57\x77\x00\x0c\x8b\x51\xa4\x86\x40\xd7\x0e\x6b\x0c\x27\xe0\xd9\x9d\x46\xb5\x81\xa5\xa4\xf5\x14\xd5\x90\x94\x50\x23\x3f\xe5\x78\x69\xf5\xb4\x50\x1c\xf3\x4a\x20\xdd\x28\xd4\xb3\x4a\x14\xea\x80\x1a\x2a\x95\xf0\x51\xeb\xe2\x46\xc0\x67\xe5\x20\xe0\x8e\x8c\xf4\x7a\x6e\x7e\xe1\xdf\x2d\x4d\x4e\x65\xfe\xa3\x76\xcb\x01\xf0\x45\x12\xb2\x88\x40\xad\x66\x23\x5b\x71\x0b\x9d\xc3\xc1\xb3\x50\x28\xf4\x69\xfe\x36\x2a\xfd\x52\x68\x13\x29\xd4\x0f\x78\x20\x15\x86\x6c\xb3\x6c\xa8\x59\x12\xa5\x9c\x2a\x02\x46\xa5\x56\x0d\x26\x11\x4b\xb0\xac\xc6\x74\xd3\x65\x01\x47\xd7\x17\x49\xa2\x4b\xaf\x2c\x28\x24\x26\x15\x38\x45\x39\x0b\xd1\xb5\x16\x13\xe8\x38\x00\x5c\x44\x6e\x11\x3f\x2a\x25\x54\x89\xbc\xa0\xb2\xfc\x98\x12\xfe\x1a\x9a\x78\x3d\x7f\x12\x34\x46\xf7\x76\x0a\xf3\x41\xf7\x7b\xcb\x5d\x07\x80\xfa\xd6\x51\x77\x0d\x87\xd3\xf1\xd1\x57\x68\xe6\x81\x54\x31\x33\x74\x8d\x58\xc3\xa4\x04\x8d\x90\xcc\x27\xf0\x84\xc5\x92\xe3\x53\x7b\x28\x67\xe1\xb8\x81\xbc\x1a\x4f\xa1\xb1\xc9\xb7\x1d\x2e\xa2\x88\x25\x91\x9d\x24\xbf\x62\xa5\x6a\xc1\x91\x00\x4b\x42\x01\xf9\x5a\x82\xf1\xeb\x83\x6c\xe7\x34\xdb\x3b\x19\x8f\xbe\x64\x97\xa3\xc9\xe1\xcb\xc9\xd9\xc1\x8c\x12\x32\x3e\x55\x0e\xb3\xb5\x04\xe3\xf7\xdb\x57\xe7\xa7\x77\x18\x92\x9a\x01\x81\x5a\x43\x48\xd3\xb0\x26\x30\x1f\x1b\xf6\x09\xb9\xd3\x8c\x1b\x5c\x44\xba\x06\xff\xac\x25\x28\x3a\x5d\x5d\x7e\xca\x0e\xf7\xaf\x7f\x1c\x67\x17\x5b\x7f\x7e\xbd\xba\x39\xdf\xbf\x3e\xfa\x3c\xf9\xb6\x9b\x9d\xed\x15\xe0\xd5\xcf\xdd\xc9\x87\x2d\xdb\x65\xf2\xf1\x28\xfb\xfd\xae\x4c\x6c\xc0\xb4\x11\x56\xf3\xf2\xfd\x5d\xc7\xa3\xd3\x9b\xd1\xf7\xd2\x31\x96\x54\x5e\xd8\x10\x67\x84\xfd\xe3\xec\xcd\xd7\xeb\x8b\xb7\xd9\xf6\x59\xb6\x73\x32\xb3\xa6\xa2\x15\x7e\x54\x9c\x85\xb4\xdc\x1c\xe7\x6f\x00\x00\x00\xff\xff\xc8\x9b\x12\x0a\x0e\x05\x00\x00")

func manifestConfigConfigYamlBytes() ([]byte, error) {
	return bindataRead(
		_manifestConfigConfigYaml,
		"manifest/config/config.yaml",
	)
}

func manifestConfigConfigYaml() (*asset, error) {
	bytes, err := manifestConfigConfigYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "manifest/config/config.yaml", size: 1294, mode: os.FileMode(420), modTime: time.Unix(1689920880, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"manifest/config/config.yaml": manifestConfigConfigYaml,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//
//	data/
//	  foo.txt
//	  img/
//	    a.png
//	    b.png
//
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"manifest": &bintree{nil, map[string]*bintree{
		"config": &bintree{nil, map[string]*bintree{
			"config.yaml": &bintree{manifestConfigConfigYaml, map[string]*bintree{}},
		}},
	}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}
