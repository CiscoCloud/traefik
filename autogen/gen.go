// Code generated by go-bindata.
// sources:
// templates/consul_catalog.tmpl
// templates/docker.tmpl
// templates/kubernetes.tmpl
// templates/kv.tmpl
// templates/marathon.tmpl
// templates/notFound.tmpl
// DO NOT EDIT!

package autogen

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

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _templatesConsul_catalogTmpl = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x8c\x54\x4d\x6f\xe3\x20\x10\xbd\xe7\x57\x8c\x50\x0e\xbb\x52\x4b\xf7\x1c\x29\x87\x46\xda\x55\x4f\x55\xb5\xbb\xb7\x28\x07\x6c\x4f\x62\x54\x17\x2c\xc0\xdb\x56\x96\xff\xfb\x82\xf9\x08\xb4\x6e\x94\x53\xc0\xf3\x78\xf3\xe6\x3d\xc8\xbe\x62\xf5\x33\x8a\x46\x1f\x56\xe3\xa8\x98\x38\x21\xac\xb9\x68\xf0\xed\x06\xd6\x42\x36\x08\x9b\x2d\xd0\x47\xbb\xd0\xd3\xb4\x02\x18\x47\x7e\x04\x81\xf0\xed\x84\xe6\xde\x18\xc5\xab\xc1\x20\x10\x14\xac\xea\x90\xf8\x23\xf4\x0f\xaa\x7f\xbc\x46\xfa\x97\x9d\x34\x10\xa3\x06\x24\xdf\x81\x1c\x59\xa7\x91\xcc\x2c\x00\xa9\x2d\x0d\x8b\xdb\x71\xb4\x94\x3b\xbf\xf1\x3c\xd3\x44\xb5\x65\x42\xa5\x69\x5e\x7c\x64\x2f\x18\xb4\x79\xa5\xd3\x74\x98\x39\x01\x06\xd5\xc1\x16\xc8\x8c\xce\xd4\xf5\x4a\x1a\x59\xcb\x6e\x59\x5f\x6b\x4c\x6f\x65\x6d\xee\xee\xf2\x2e\xf7\x4d\xa3\x50\xeb\xa8\x64\x33\x8e\xe5\xd9\x27\xa9\xcc\x34\x91\xd0\xd8\x56\x5f\x91\x9f\x5a\xe3\xec\x2a\x9b\x87\xf9\xa8\xaf\x2f\x4b\x88\xae\x38\xa2\x57\x6e\x5a\x08\x6c\xe9\x33\x40\xa0\xdf\x9e\x5b\x65\x67\x2c\x7f\xd8\xc5\x75\xfc\x4d\xa1\xc6\x8e\x31\xc6\xb5\xf6\xfb\x39\xdf\x50\x73\xce\xc6\x72\xcd\x55\x3d\x70\xb3\x53\xc8\x9e\x51\x7d\x3d\x56\xc0\x55\x1e\x47\x80\x26\x4c\x9a\x2b\xce\x54\x52\xce\xa5\xa5\x6b\x10\x95\xd9\xf8\x4b\x72\x9f\x32\xbe\xf5\x2e\x18\x2e\x85\x8f\xfa\x13\x2d\x59\x9d\x6d\xf0\xb3\x74\x92\x35\x3b\xd6\x31\x51\x5f\x9a\xc4\xa1\xaa\x80\xba\x38\x47\x4e\x77\xcd\x14\x39\xb1\x9f\xe1\x05\x4d\x2b\x9b\xa0\xbf\xa4\xcb\xd5\xa7\x29\xf6\x47\x25\x85\x29\x5f\x69\x19\xe8\x19\x41\xe3\xca\x8a\x28\x83\x75\xbd\x83\x40\xd7\xfa\xac\xb5\x84\x39\x05\x3d\xd3\xfa\x41\x6a\xf3\x80\xac\xb1\x9e\xb9\x5b\x57\x9a\x16\x9b\xd0\x12\xf9\xc1\xb7\xf9\xe9\xcf\xfa\x7a\xc5\xa5\xe2\xe6\xfd\x22\x55\xc0\x7c\x20\xf9\x11\xdd\x5f\xa3\x30\xea\xfd\x49\x72\x61\xf4\xe7\x18\x13\xcd\x8c\xea\x67\xd4\xc5\x18\x33\xb6\xf0\x78\xb2\x93\x56\xe7\x3e\x3a\x6d\xfb\xfc\xcc\x3a\x2f\x1c\x04\x17\x24\xb5\xd6\xdd\xe4\x8f\xf0\xb0\xca\xdf\xe6\x35\x09\x51\x25\x9d\x52\xff\x73\xdb\x5a\x5b\x17\x43\x04\x50\x43\x87\xe9\x9f\xee\x57\xa0\xfb\xed\x3e\x3a\x15\xe9\xe2\xfc\x0f\x00\x00\xff\xff\x09\xd7\xe6\xea\xdd\x05\x00\x00")

func templatesConsul_catalogTmplBytes() ([]byte, error) {
	return bindataRead(
		_templatesConsul_catalogTmpl,
		"templates/consul_catalog.tmpl",
	)
}

func templatesConsul_catalogTmpl() (*asset, error) {
	bytes, err := templatesConsul_catalogTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "templates/consul_catalog.tmpl", size: 1501, mode: os.FileMode(436), modTime: time.Unix(1466539430, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _templatesDockerTmpl = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x8c\x52\x41\x4f\x32\x31\x10\xbd\xf3\x2b\x26\x0d\x47\x28\xdf\x99\x84\xc3\xa7\xd1\xe0\xc5\x6c\xbc\x78\x20\x1c\xea\xee\x88\x8d\x6b\x4b\xa6\x45\x25\xb5\xff\xdd\x69\x77\xbb\xb0\x44\x13\xb9\xb4\xf3\x66\xde\xa3\xef\xcd\x6e\x9e\x54\xfd\x8a\xa6\x71\xdb\x10\x48\x99\x1d\x82\xbc\xb6\xc6\x2b\x6d\x90\x5c\x8c\x13\xe0\xdf\x30\x23\xfb\xcb\x3c\x84\x1d\xfa\xab\xae\x00\x19\xa3\x74\x48\xef\x4c\xe8\x4f\xee\xcb\x7b\xf5\x86\xf0\x05\x84\xfb\x56\xd5\x08\x62\x21\x40\x88\x73\x40\x32\x30\x17\x31\x6e\xf3\x7f\x1c\xa8\x85\x15\x88\x2c\x5c\x91\xf5\xb6\xb6\x6d\x52\x5e\x2e\x16\x19\xbb\xab\xfe\x37\x0d\xa1\x73\x19\xec\xc6\x2c\xf9\x54\x89\x2c\xf0\x81\x7a\xf7\xe2\x59\x23\xf7\x1e\xbb\x8a\xbb\x93\x10\xf8\x91\x7c\x4e\x36\xcf\xc4\xce\x46\x5e\xa7\x05\x9a\xc1\xb4\x1e\x6c\xc3\x72\x05\xf2\xb6\x0c\xe7\x10\x4e\x5c\x29\xca\x95\x5d\x0e\x7c\x7e\x05\x8b\x9e\x34\x92\x84\x36\x0d\x7e\x8e\x74\xff\x65\xad\x3e\xc4\x64\xf7\xc7\x3c\x4f\x8c\xce\xdb\x5e\x39\xb7\xb6\xce\xaf\x51\x35\xac\xdc\x3b\xac\xc6\xe8\x39\x29\x71\x48\x5b\xd2\xfe\x38\x4c\x97\xfa\x62\x0e\x8d\xa7\x63\x65\xb5\xf1\x8e\x47\x37\x25\x17\x66\xdc\x9c\x75\x2e\x48\x90\xf6\x94\x82\x9f\x71\xd1\xc7\xdb\x6d\xf1\x0f\x31\x49\xb2\x07\x8f\xdc\xcf\xe7\xfc\xb7\x30\xb3\x1c\x1d\x5a\x1c\xbe\x8a\xb2\x90\x87\x04\x8e\x33\x2a\x2b\xfe\x0e\x00\x00\xff\xff\x8a\x53\xa5\xc9\xcd\x02\x00\x00")

func templatesDockerTmplBytes() ([]byte, error) {
	return bindataRead(
		_templatesDockerTmpl,
		"templates/docker.tmpl",
	)
}

func templatesDockerTmpl() (*asset, error) {
	bytes, err := templatesDockerTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "templates/docker.tmpl", size: 717, mode: os.FileMode(436), modTime: time.Unix(1466539430, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _templatesKubernetesTmpl = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x7c\x50\xb1\x4e\xc5\x30\x0c\xdc\xdf\x57\x44\x4f\x8c\x28\x1f\x80\xd4\x85\x01\x75\x40\x08\x15\x21\x86\xaa\x43\xa0\xa6\x20\x4a\x8b\x92\x14\x86\x28\xff\x8e\xeb\xc4\x4e\xdb\xe1\x75\x69\x7c\xe7\xb3\xcf\xd7\xbe\x9a\xb7\x2f\x98\x7a\xd7\x85\x60\xcd\x34\x80\xba\xca\xc8\x83\xf9\x86\x6b\xa9\xd4\x4d\xa5\xf4\x6d\xee\x8d\xf1\xa4\xf0\x13\x85\x03\xfb\x0b\x36\x0b\x52\xb1\xf6\xb3\x56\x3f\x11\xc4\x32\x59\xa9\xcf\x21\x6c\xb7\xc5\x78\xd6\x49\x9d\xa8\x32\x16\x99\x8e\xb4\x8b\x1d\x55\xa5\x0a\xa9\x9f\x9b\x7b\x24\x89\xfb\x83\xcf\xe1\xc3\x23\x5d\xd8\x17\x82\xc4\x2e\xee\xc1\x37\xff\x4f\xed\xbb\x9d\x27\xbf\x3f\x9e\xa1\x7c\x0c\x97\x74\xfe\x1d\xb7\xd3\xc0\xa2\x26\xb3\x5b\x61\xb6\xcb\xd1\x25\xc3\xdc\xc0\x29\x26\xdb\x3f\xc6\xb9\x7a\x76\xbe\x06\xd3\x63\x6a\x64\x5e\x3a\x1f\x77\xe4\x31\x75\x3b\x2f\x1e\xb2\x4f\x7a\x53\xe6\x22\x6e\x56\x48\x32\xbf\x68\x56\x93\x3c\x51\x32\x55\x32\xb7\xcb\x08\xf9\x06\x22\x75\x83\x00\x87\x7e\x0c\xf5\x3f\x00\x00\xff\xff\xa1\x49\x18\x50\x50\x02\x00\x00")

func templatesKubernetesTmplBytes() ([]byte, error) {
	return bindataRead(
		_templatesKubernetesTmpl,
		"templates/kubernetes.tmpl",
	)
}

func templatesKubernetesTmpl() (*asset, error) {
	bytes, err := templatesKubernetesTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "templates/kubernetes.tmpl", size: 592, mode: os.FileMode(436), modTime: time.Unix(1466539430, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _templatesKvTmpl = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x84\x54\xc1\x6e\xdb\x30\x0c\xbd\xe7\x2b\x04\x61\xc7\xc1\xde\x79\x40\x0e\xeb\xb0\xad\x87\x1e\x0a\xec\x18\xe4\xa0\x38\x4c\x23\xd4\x91\x0c\x5a\x5e\x53\x18\xf9\xf7\x51\x26\x25\xcb\x49\x8b\xe4\x56\xbe\xf7\xc8\xf7\x68\xb1\xe3\xf8\xe5\x80\xde\x05\x70\xfb\x5e\x7d\x5f\xab\x27\xdb\x07\x55\x3d\x23\x1c\xec\x59\xe9\x3a\x63\xb5\x56\x97\xcb\x8a\xd8\x3b\xd3\xbc\x26\xf2\x35\x3b\x61\xb5\x26\xee\x6a\x93\xfe\xdc\x8e\x23\x1a\xf7\x02\x2a\x8b\x17\xad\x62\xa7\x8a\x2b\x3d\xe0\x3f\xc0\xd9\x48\xa6\xe8\x5a\x20\xf6\x11\xb9\x8d\xc5\x66\xb0\xe1\x01\xc1\xbc\x02\x46\xc9\x1f\x08\x4a\x6b\x55\x11\x5b\xc0\x1d\x83\x24\xd2\x70\xee\x10\xfa\xde\x7a\xa7\xa7\x59\x6f\x36\x1c\xd5\x55\x13\x02\xb2\xe9\x4a\x8f\xe3\x93\x29\x3c\x5c\x2e\xba\x5a\xd2\xb7\x2b\x45\xbf\xb9\xb1\x5a\x2b\x7d\x63\x8c\x64\x34\x6d\xd2\x4f\xb6\x5b\x6f\xf6\x0f\xa6\x35\xae\xb9\x31\x1d\xa1\x9d\x40\xd1\xf2\x09\xc2\xd1\xef\x4b\xbb\xa5\xf8\xae\xd9\x92\xcc\x56\xb9\xa1\xd8\x5c\xf6\x5a\x9a\x3c\x99\xf3\x4f\xef\xdc\x8f\x53\xb8\xb2\x48\x40\x43\x40\x74\x67\x4e\x7e\x70\x81\xdd\x25\xc1\xaf\x73\x40\xd3\x04\x8f\xbf\x07\xd7\x7c\x2e\x85\x44\x3b\x10\xad\xcc\x37\xcf\xbd\x2d\x2e\x7a\xdf\x0d\x2f\x22\xce\xcd\x56\x29\xf7\x22\x1a\xf5\xe0\xef\x57\x7a\xe6\xdd\x7c\x3c\x72\xde\x51\xb1\x2b\x79\xd9\xf2\x3c\xef\xfa\x12\x5e\xc6\xe8\xe1\x6b\x36\x39\x60\xcb\xd3\x8b\x9d\x51\x4d\xc7\xb9\x11\x7f\x03\xfb\x72\xe4\x10\xc2\x98\x28\x5c\xd6\x85\xa9\x6c\x6e\x93\xaf\x77\xbe\xc0\x5c\x92\xf4\xc5\xfd\x4f\x57\x27\x9e\x12\x06\x2e\xe0\xfb\xb3\xb7\x2e\x4c\x47\xf9\xb7\x6b\x6d\x88\xd3\xe3\xe4\x09\xeb\x26\x4c\x8b\x62\x9e\x18\x03\xe6\xce\x39\x63\x3a\xe7\xeb\x9c\x52\xcf\x59\x3b\xd3\xf7\x8f\xbe\x0f\x8f\x60\xf6\x74\x26\x39\x73\xc0\x01\x58\xb1\x64\xa4\xf9\x1d\x5a\x8f\x36\xbc\xcf\x8a\x6f\x42\x17\x20\x11\xcb\x5c\x6b\xb5\xc9\xeb\x29\xea\xc2\x54\xd1\x6a\xfc\x4a\x5f\x65\x27\x53\x9e\x6d\x5a\x10\xfa\x21\x40\xf1\x9f\x93\x66\x71\xa9\xd6\xb9\x41\x24\x4a\x7f\xc6\x0a\xe4\xf3\x8d\x55\xcc\xfd\xe0\xa5\xc4\x1f\x0e\x2d\xdc\xac\x31\x16\xf3\x0e\x0b\xb7\xf9\x49\xfc\x0f\x00\x00\xff\xff\x49\x00\xf6\xb8\xee\x05\x00\x00")

func templatesKvTmplBytes() ([]byte, error) {
	return bindataRead(
		_templatesKvTmpl,
		"templates/kv.tmpl",
	)
}

func templatesKvTmpl() (*asset, error) {
	bytes, err := templatesKvTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "templates/kv.tmpl", size: 1518, mode: os.FileMode(436), modTime: time.Unix(1466539430, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _templatesMarathonTmpl = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x94\x52\x4f\x6b\x3a\x31\x10\xbd\xfb\x29\x86\xe5\x77\xfc\x19\xef\x82\x87\x96\xb6\x58\xe8\x41\x4a\xa1\x07\xf1\x10\xd7\x51\x83\x61\x13\x26\xb1\x45\xd2\x7c\xf7\xe6\xbf\x4a\x4b\xc1\xbd\xcc\x24\xf3\xf2\x78\xef\xcd\x3a\xf7\x8f\x6b\x6d\x60\x3a\x03\x76\xa7\xb5\x14\x3d\xb7\x42\x0d\xc6\xfb\xd1\x72\xcd\xfb\x03\x0e\x1b\xb3\x72\x8e\xf8\xb0\x43\x60\x6f\xdc\x1c\xe2\x08\xc2\xd7\xc6\xac\x34\xce\xed\xd0\xde\xe7\x1e\x18\x24\x5e\xef\x99\x41\xfa\x40\x32\xa5\x8e\x9d\x63\xcf\x0f\xf0\x05\x84\x5a\xf2\x1e\xa1\x63\x1d\x74\xe3\xce\xfb\x55\x62\x3d\x92\x84\x19\x74\x89\x6b\x41\xca\xaa\x5e\xc9\x33\xd9\x74\x32\x09\xef\xe7\xca\xd8\xd0\x67\x8c\x22\x7b\x9e\x77\x89\xe3\x13\xc5\x6e\x6f\x03\x4d\x42\xbc\xe7\x53\xc3\x8c\x9c\x0b\x02\x43\x1d\x2d\xb7\xa4\x06\x7b\xe5\x30\x62\x7e\xcb\x22\xd8\x6d\x60\x56\xbb\x1f\x56\x26\x17\x56\x4a\x28\xd1\xcc\x65\x3e\x4f\xe5\x6d\xcb\x29\x8b\xd6\xdc\x98\x68\x6b\x8e\x7c\x83\x54\xa5\x2f\xae\x6f\x59\x12\xa2\x49\x28\x12\xf6\xd4\x40\xf5\x9c\xc7\x38\x58\x3a\x2d\x94\x18\xac\x09\x88\x65\x35\x16\x80\x8f\x17\x13\x56\x96\x18\x82\x8e\x0a\xfe\x87\x43\x89\x25\xaf\xe1\x16\xb7\x8c\xd4\xd1\xa2\xc9\x65\xbc\x0f\x7a\xff\x0c\x06\x80\x8e\x12\xdb\x92\x6b\x20\xaf\xf1\xb2\xa4\xe1\x9c\xd8\x46\xc5\x2f\x7c\x8d\xb2\x8a\xbd\x49\x92\x4c\x2f\x57\x89\xab\x05\x70\x45\x57\x9c\x37\x1d\x69\x98\xf7\xdf\x54\xe4\xff\xe4\xdc\xd5\xfa\x1d\x00\x00\xff\xff\x4d\x59\x86\xf0\x35\x03\x00\x00")

func templatesMarathonTmplBytes() ([]byte, error) {
	return bindataRead(
		_templatesMarathonTmpl,
		"templates/marathon.tmpl",
	)
}

func templatesMarathonTmpl() (*asset, error) {
	bytes, err := templatesMarathonTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "templates/marathon.tmpl", size: 821, mode: os.FileMode(436), modTime: time.Unix(1466539430, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _templatesNotfoundTmpl = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xb2\x51\x74\xf1\x77\x0e\x89\x0c\x70\x55\xc8\x28\xc9\xcd\xb1\xe3\xb2\x81\x51\xa9\x89\x29\x76\x5c\x0a\x40\x60\x53\x92\x59\x92\x93\x6a\x17\x52\x94\x98\x9a\x96\x99\x6d\xa3\x0f\xe1\x72\xd9\xe8\x43\x94\xd8\x24\xe5\xa7\x54\x42\x54\xfa\x67\x00\x81\x42\x6e\x62\x9e\x8e\x42\x49\x46\x66\xb1\x02\x10\x25\x25\xa6\xe8\xe9\xe9\x01\x15\x43\x54\x01\x35\x81\x8c\x07\x04\x00\x00\xff\xff\xb8\x78\x48\x56\x75\x00\x00\x00")

func templatesNotfoundTmplBytes() ([]byte, error) {
	return bindataRead(
		_templatesNotfoundTmpl,
		"templates/notFound.tmpl",
	)
}

func templatesNotfoundTmpl() (*asset, error) {
	bytes, err := templatesNotfoundTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "templates/notFound.tmpl", size: 117, mode: os.FileMode(436), modTime: time.Unix(1466539430, 0)}
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
	"templates/consul_catalog.tmpl": templatesConsul_catalogTmpl,
	"templates/docker.tmpl": templatesDockerTmpl,
	"templates/kubernetes.tmpl": templatesKubernetesTmpl,
	"templates/kv.tmpl": templatesKvTmpl,
	"templates/marathon.tmpl": templatesMarathonTmpl,
	"templates/notFound.tmpl": templatesNotfoundTmpl,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
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
	"templates": &bintree{nil, map[string]*bintree{
		"consul_catalog.tmpl": &bintree{templatesConsul_catalogTmpl, map[string]*bintree{}},
		"docker.tmpl": &bintree{templatesDockerTmpl, map[string]*bintree{}},
		"kubernetes.tmpl": &bintree{templatesKubernetesTmpl, map[string]*bintree{}},
		"kv.tmpl": &bintree{templatesKvTmpl, map[string]*bintree{}},
		"marathon.tmpl": &bintree{templatesMarathonTmpl, map[string]*bintree{}},
		"notFound.tmpl": &bintree{templatesNotfoundTmpl, map[string]*bintree{}},
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
