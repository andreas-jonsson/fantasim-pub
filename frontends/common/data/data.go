// Code generated by vfsgen; DO NOT EDIT.

package data

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	pathpkg "path"
	"time"
)

// FS statically implements the virtual filesystem provided to vfsgen.
var FS = func() http.FileSystem {
	fs := vfsgen_FS{
		"/": &vfsgen_DirInfo{
			name:    "/",
			modTime: time.Date(2018, 2, 17, 12, 21, 22, 798722915, time.UTC),
		},
		"/16x16_sm_ascii.png": &vfsgen_FileInfo{
			name:    "16x16_sm_ascii.png",
			modTime: time.Date(2018, 2, 17, 12, 21, 22, 798722915, time.UTC),
			content: []byte("\x89\x50\x4e\x47\x0d\x0a\x1a\x0a\x00\x00\x00\x0d\x49\x48\x44\x52\x00\x00\x01\x00\x00\x00\x01\x00\x01\x03\x00\x00\x00\x66\xbc\x3a\x25\x00\x00\x00\x06\x50\x4c\x54\x45\x00\x00\x00\xff\xff\xff\xa5\xd9\x9f\xdd\x00\x00\x00\x01\x62\x4b\x47\x44\x00\x88\x05\x1d\x48\x00\x00\x00\x09\x70\x48\x59\x73\x00\x00\x0b\x13\x00\x00\x0b\x13\x01\x00\x9a\x9c\x18\x00\x00\x00\x07\x74\x49\x4d\x45\x07\xe2\x02\x10\x07\x29\x25\xb7\x8c\x49\xff\x00\x00\x0a\x1c\x49\x44\x41\x54\x68\xde\xbd\x99\xcf\x8b\xe4\xc6\x15\xc7\x4b\x11\xac\x2e\x43\x17\xc6\x97\x4e\xdc\x3b\x05\x3e\xe7\xb0\x66\x03\x51\x93\x1a\x29\x10\x43\xfe\x02\x43\x8e\x03\x86\xdd\x4b\x0e\x63\x7c\x70\x1b\x8b\x69\x1d\x7d\x58\x7c\x8e\xa1\xf1\x1f\xe0\xbf\x60\xc1\x73\x28\x7c\x70\x7c\x58\x72\xde\x83\x6d\x44\x02\x3e\x19\xa3\x63\x1d\x8a\xee\xbc\xef\x7b\x55\x92\xba\xa7\x67\x76\x37\xd9\x1d\x89\xfe\x25\x7d\x5a\xef\xd5\xab\xf7\xab\x24\xa5\x0e\xb6\xdd\x4e\x5e\x93\xad\x0e\x75\x50\x2a\x6b\x73\x97\xb5\x23\x50\x74\x2a\x14\x1e\x47\xd4\xf9\xbd\xf5\x76\x51\xe6\xae\xe8\x72\x37\x00\x41\xf7\xd9\x56\x87\xf7\x9e\x10\xf0\xe7\xdf\xac\xb7\x55\x53\x74\xd8\x71\x52\xf7\xbd\x56\xe1\xe4\x41\xf6\xdd\x6c\xb1\xdc\x10\xf0\xb7\xdf\x3d\xfe\x79\xbd\xd5\x7d\xd1\xe1\xe5\x8d\xf1\x5d\xa1\xac\x00\xc6\x0f\x80\xf1\x7f\xfa\xc1\x78\xfc\x1b\x80\xc6\x15\x48\x84\xee\x07\xa0\x0e\x4d\x55\x07\x00\x8b\xf2\x3f\xff\x34\xbf\x10\x10\x74\x98\x5d\x44\x1d\xb0\x37\x15\x04\x09\x30\xff\x3b\x29\xe9\x66\x8b\xa6\x8a\x40\x1d\x0e\x80\xf3\xa2\x1b\x80\xbf\xfe\x7b\x75\x6a\x3c\x44\x4c\x81\xac\x2d\x08\x60\x11\x00\x74\x0f\x25\xab\x26\x29\x09\x40\xbb\x99\x65\x25\xff\xf2\x6d\x53\xc9\x10\xef\x3f\x4d\xc3\x34\x64\x13\xed\x4e\x1b\x1e\x26\x74\x80\x1d\x93\x25\x71\x15\x8d\xdf\x6d\xd5\xb0\xa1\x60\x6a\xcc\x44\xb2\xa4\xbc\x30\x0b\x55\xc9\xa6\x1e\x27\xeb\x10\x58\xf0\xfb\x8b\xa7\xfb\x25\x36\xb9\x10\x34\xce\x5a\x88\xda\x77\x97\x36\x0b\xb9\xb3\xd6\xb0\x26\xb9\xc3\x3e\xfe\xf1\xc1\x89\xaa\x5d\x1e\x8a\xce\xda\x3a\x94\x0b\xf8\x11\x06\x9b\xd1\x10\x01\xe4\x8e\x8e\x91\x49\x68\x92\xad\x7d\xb8\x79\x40\x87\x74\xaf\x7b\xba\x02\x79\x8c\x52\x27\x0f\xf0\x47\x3e\x24\xc0\x1f\x3d\xfe\x63\x3c\x01\xab\x53\x55\xaa\x45\xa9\xfb\xaa\x51\x85\x37\x1d\x74\x78\xb8\x11\x1d\x00\xac\x43\xbd\x2d\xc9\x3e\x98\x42\x95\x07\xd2\x82\x80\x7a\x53\x2e\x60\xf2\x77\x9e\xbf\xfd\x3d\xbc\xa3\x54\x70\x82\x01\xa0\x41\x46\xe0\xed\xef\xdf\x79\xbe\x07\x88\x08\xba\xf8\x06\x4e\x03\x15\x8d\x87\x88\xda\xe3\x34\xd9\x26\x29\x99\x6f\xcc\x2f\xeb\x2d\x06\x49\x07\x49\xc9\xda\x2f\x4a\xf8\xc7\x30\xcc\x7c\xa3\x4e\x20\x22\x77\x74\x30\x0e\xd3\x5a\x9a\xbe\x64\xa8\x9c\x75\x90\x39\x4d\x86\x2a\x17\x24\x3c\x99\x3a\xe7\x61\x42\xc4\x68\xea\x07\x27\x71\xba\xc7\xc9\x62\x43\xb7\xea\x4d\x6c\x8b\xf2\xde\xf9\x8d\x27\x73\xb7\x28\x01\x2c\xe6\xa6\x55\x5d\xde\x15\x6e\xe2\x34\x73\x05\xc5\xaa\x06\xd1\x5d\xd9\xda\x65\x3d\x3c\xe2\xb7\xff\x18\x00\x3b\x00\x75\xa0\x28\x70\x79\x5f\xb4\x19\x05\xd0\x00\x34\x03\x60\x6d\xe5\x97\x1c\x2e\x59\xc7\x4e\x23\x9b\x8f\xc0\xa2\xb4\xea\x7e\x4f\x80\x03\x50\x74\x30\x1a\x47\x4b\x3f\x00\xb5\xcf\x3b\xb3\xc9\x5b\x00\x12\x4a\xec\x40\xdd\x30\x0a\x13\x0a\x67\xae\x08\x4f\x00\x5f\xa5\x70\x83\x1d\x94\xd5\x4f\x97\x5e\x00\x88\xc8\x7a\x00\xba\xe5\x2b\x28\x05\x25\x4d\x53\xf6\x02\x40\xc9\xbc\x67\x0f\x55\xac\x83\x52\x71\x98\xf1\x0a\x18\x26\x8e\x92\x76\x23\x60\xbc\x6d\x64\xc2\xc5\x50\x90\x4e\x7e\xa2\xa2\x08\x98\x7a\x5e\x9a\x0d\xf4\x16\x53\xe7\x2d\xce\xcc\x5f\x6a\xb2\xde\xec\x46\x6e\x48\xb1\xd8\x1a\x6f\x28\x0a\xf0\x2e\x3b\x3b\x2e\x0d\xd9\xc0\xdc\x35\xf9\x34\x4c\x8b\x0c\x91\x76\x0c\x15\xd3\xa4\x5b\x0a\x7f\x72\x78\xa7\xf0\x81\xef\x8a\x3f\x2d\xa7\x91\x4e\xfb\xc2\xdd\x00\x40\x44\xe1\x34\xcd\xef\xfc\x38\x00\xf5\x74\xab\x7d\xd6\xe7\x41\x9d\x5d\xd1\x94\x50\xe9\x20\xd9\xbe\xf6\xca\x8a\x2f\x40\x41\xf8\x53\xe1\x19\xa8\xbd\x0e\x26\x70\x9e\x23\xc0\xc0\x93\x5c\x04\x3a\x16\x81\x4b\x8b\x18\x15\x25\x91\x8a\x22\xc2\x1d\x05\x38\xb1\xb3\x92\x70\xfc\x23\x00\x6a\x85\x0c\x13\xe1\x1f\x0d\x15\xd5\xb3\x62\x28\xdd\x8a\xa1\x48\x87\x64\x6a\x55\x8e\xef\x24\x42\x89\xa9\xe1\x4c\xb7\x6e\x07\x89\xfd\x7f\xdb\xc4\x6e\xb7\x6c\xf0\xe7\xb5\xc7\x6b\x1d\xd6\x94\x38\x29\x3e\xbc\xba\x6c\xec\x25\x4f\x0b\x7d\x3f\xbf\x47\x19\x36\xe0\xb5\xe6\xcf\x11\xa8\x1a\xf8\x08\x01\xf6\x60\x63\xf3\x36\x96\x84\x57\x30\xdb\xe3\x9f\x87\x33\x71\xd6\x19\xf0\x00\xea\x09\x10\x67\x29\xf9\x47\xdd\x33\x70\x45\x07\x1f\x3d\x82\x27\xe0\x44\xdd\xd5\x9d\x65\xd3\x13\xc0\xdf\xa0\xd1\x21\xd0\x8c\xc0\xd9\x15\x03\x8f\xb7\x7b\x22\x9a\x51\x84\xb5\x67\x13\x60\x4f\x49\xcb\x4a\xd2\x3b\x80\x73\x35\x19\xe1\x08\x34\xfc\xbb\x89\x40\x32\xd4\xa5\x12\x43\x51\xbb\x42\xc7\x40\x90\xa1\x90\x7c\x92\xa9\x2f\x95\x61\xc0\xf8\x04\x90\x55\x5f\x38\x59\x6f\x76\xe3\xe2\x24\x0a\x42\xf3\x71\xb4\xe8\xe9\xe6\x2a\xa6\x74\x19\xa2\xd8\x7d\x0a\xa0\x2f\x28\x24\x70\x64\x7b\xf8\x4c\x3e\xab\x86\x1c\xa5\xd1\x8a\x52\x68\xcf\x7d\xd8\xe8\x26\xf2\xc9\xd1\xe1\x35\xa2\xbb\x3f\x5d\x8d\x80\x4a\x00\x2a\x79\x46\xa5\x99\xd2\x78\xcf\x99\x82\xb3\x13\x65\x06\x01\xce\xae\x62\xfa\x41\x21\x88\x80\x87\x27\x18\x2e\xcf\x55\x83\xb8\xe4\xf4\x83\x52\x22\x00\xcf\x7b\xaf\x58\x07\x0c\x16\x83\xa3\xdc\x80\x62\xd4\xc7\xba\xc7\xb3\xcf\x80\xee\xd1\x52\xa0\x4a\x68\x04\x6f\x02\x48\x8b\x46\x00\x54\x9a\xaa\xc9\x5d\x65\x35\x0a\xa2\x00\x34\xc5\x34\xb9\xb8\x34\xbc\xf1\x74\x25\x18\x59\x40\xba\x51\x00\xaa\x69\x2c\x3a\x6a\x43\x9d\xbc\x74\x1d\x0c\xcc\x01\xac\xb7\xa9\xfc\xde\x18\xb9\x77\xe1\x14\x45\x2c\xf1\x78\x43\x17\xd3\x20\x2f\xd9\x46\xe9\x94\x9d\x50\x5b\x13\x00\x07\xbe\x0e\xb8\x38\x0c\x4e\xf5\x69\x48\x1c\x80\x51\x44\x47\x33\x80\x42\x24\xe9\x87\x7a\x38\xfa\xcb\x25\x00\x29\x4f\x34\x6b\x63\x29\x22\xb3\x53\x89\x66\xa0\x91\xfc\x4b\x42\x33\x37\x06\xe7\xba\x1d\x33\x4c\xee\xce\xae\x62\x71\x41\xb8\x71\x7e\x08\x23\x50\x77\x13\xa0\x4e\x00\xe7\xa5\x94\x1f\x26\xc0\x34\x81\x98\x08\xd8\xd5\x21\x10\x5d\x57\x25\x1d\x06\x2f\x87\x88\x2d\x07\x0e\xe9\x7f\xc9\x09\xd2\xf8\x53\x7e\x8f\xf9\x81\x94\xdc\x22\xf4\x0c\xa5\x64\x24\x10\x9c\x32\x21\x01\xbc\x42\x39\x28\xe6\xb0\xae\xbf\x93\xdc\x70\xad\x18\x65\x7b\x07\x32\x0e\x1f\x35\x5d\xe3\xec\xff\x23\x07\xe0\xf6\x8e\x1e\x00\xbc\x52\x99\xa1\x2c\x92\x0d\xc4\x0e\x7b\x19\x46\x00\x0a\xc1\x1e\x29\x88\x2d\xb9\x0f\xa0\x28\x93\xa7\x3c\xa1\x14\xb2\x9f\x40\xce\xae\x10\x7c\x99\x67\xf9\x3d\xb2\x6d\x04\x24\x81\x20\x3f\xf0\x42\xa0\x17\xe0\x43\x4f\x8b\x20\xf9\xbf\x24\x10\xe4\x07\x44\x26\x9a\x03\x88\xc8\x7a\x02\xc4\x17\xd5\x90\x1f\x28\x88\xc9\x2b\xe3\x7a\x0b\x5d\x69\x0d\xfd\x23\x80\x7f\x43\x2d\xd3\xa6\x05\x19\x1a\x06\x3a\x89\x69\xce\x7d\xca\x0f\x3c\xcc\x68\x28\x02\x28\x3f\x5c\x22\x89\x64\x3e\xe5\x07\x08\x4d\xa6\x96\x04\x72\x79\xcc\x27\x5e\x4b\xe7\x90\x36\x32\x03\xe5\xa7\xd3\x95\xe6\x8b\xa7\x6f\x78\x3f\x5d\x71\x3f\x49\xa6\xe5\xe1\x61\x95\x0c\xbd\xf1\x4d\xde\x4f\x57\x74\x44\x62\xc9\xda\x53\x0a\x15\x00\xf2\x4d\xde\x69\xdd\xdc\x28\x3b\xc9\x71\x55\x73\x38\x1e\x12\x81\xe9\x1a\x9b\xd5\xb4\x6b\x87\x5d\x5a\x2f\x35\x6d\x56\xd3\x2e\xc0\xa4\xa0\xab\xd4\x5d\x0c\xf5\x27\x77\xf9\x08\xd4\xdc\x4d\x62\x97\xb6\x22\xb9\x9c\x88\x60\x77\x99\x24\x21\x2f\xf5\x67\x0f\xb0\xfb\x2d\x82\xb2\x6a\x22\x62\x5a\x2f\x21\xec\x9a\x0e\x46\x2e\xbf\x95\x3d\xfe\x0a\x32\xd8\xe8\xcd\x46\xb2\x03\xef\x29\x3b\x4e\x80\x62\x6f\x75\x59\xbc\x96\xb5\x26\xfc\x72\x1a\xe8\x9d\x4c\x8b\x78\x81\xe6\xd9\x4c\x40\xed\xe9\x2c\xaf\xda\x43\xf2\x82\xa2\xb5\x71\x61\x54\x78\xf1\xee\xd8\x95\x26\x2f\x90\xf8\xc2\xa7\x66\x8b\x10\x40\x0d\xbf\x5a\x3e\xdb\x6f\x62\x01\x69\x5b\x49\xdb\x86\x5a\xb1\x74\xe3\xb8\x87\x1e\x44\xc9\x95\x94\x5d\xd5\x61\x19\x6b\xc3\xd4\xe8\xd4\x41\xf4\xd2\x19\xf2\x7d\x83\xc3\xb6\x95\xe6\x88\x3a\x18\xaa\x3e\xa4\x43\x6b\xc2\x71\x80\x44\x28\x06\x6a\xd2\xe1\xfa\x7c\x6a\xe4\x9c\x0b\x16\xf1\xb0\x9d\x02\x26\x76\x81\x86\x94\x3c\x63\x1d\x56\x0f\x37\xcb\x67\x07\x2d\x25\x0d\xb3\xa2\x2b\x30\x00\x77\x5b\x4e\xfc\x59\x7c\x1c\xee\x46\xab\xcd\x0b\x1d\xed\x10\x46\x1f\x10\x67\xc1\xca\x89\x80\x46\xbb\x63\xcb\xbf\x49\x58\xbf\x8a\x7f\x14\xf3\x62\x7e\xfc\x4c\x2f\xbb\x7e\xaa\x9f\x1e\xaf\x13\x9d\xec\xf7\xfb\xfb\x7c\x47\x46\x96\x3d\xe5\xbc\x9c\x93\xe3\x7f\xf1\x91\x8b\x21\x42\xeb\xe3\x6e\xde\x41\xff\x08\xd8\x12\x8b\xca\x45\xd9\x1d\x78\xa3\xa4\x0d\xa5\xe6\xcd\xbc\x29\xba\x59\xb9\xb8\x80\x1f\x3a\x59\xb8\x54\x16\x40\x1c\xa0\x47\xfd\x5f\x5c\xcc\x70\xd7\x09\x89\x82\xed\x67\xe5\x0a\x59\xa8\xdb\xb3\x9e\x1a\x40\x55\x92\x90\x98\xa5\x79\x82\xae\xd0\x85\x50\xbb\xb8\x5d\x53\xa6\xa7\xdd\x7d\x44\x6a\xc6\x42\x10\xdd\xec\x07\xca\xf5\x5d\xfe\xd3\x67\xae\x78\x5e\x3c\x2b\x3a\x06\xea\x09\xc0\x2e\xcb\x63\xd2\x57\x7a\x53\x74\x2c\x62\x6c\x04\xce\x02\x37\x3a\x73\x2e\x67\xbf\x37\x7f\x28\x3a\x56\x72\x04\xac\x00\x9c\x1b\xaa\x45\x55\xc5\x61\xa6\x9c\x80\xe5\x46\x4c\xcc\x78\x2f\xed\x49\x34\x14\xfa\x86\x21\x17\x0c\x6b\xbf\xb9\x9d\x9f\x44\x53\xdf\xe6\x0f\x2f\xb7\xbd\xf5\xd6\x07\x1f\x7c\xf2\x49\x4a\x5a\x2a\x36\x3d\x2a\x2d\xea\x94\x7a\xff\xfd\x2f\xbf\xfc\xf1\xc7\x5b\x80\x77\xdf\x7d\xf4\xe8\xd7\x5f\x71\x7a\xc7\xc0\xce\x61\x40\xbb\xc0\xef\x00\x3e\xff\xfc\xeb\xaf\xbf\xf9\xe6\x16\x60\xd4\xe1\x06\x60\xd4\xe1\x06\x20\xe9\xb0\x23\x24\xd0\x09\x69\xd4\x90\xb0\xf1\x9d\x90\xa4\xc3\x8d\x40\xd2\xe1\x46\x20\xe9\x70\x23\xb0\x6f\x07\x28\x16\x78\xa1\xb9\x0b\x50\x38\x77\x07\x76\x38\x02\xec\xdb\xe1\x08\xb0\x6f\x87\x23\xc0\xa8\x83\xb8\xde\xa4\x55\x95\x50\x1e\x75\xb8\x01\x90\x66\x46\x5e\xb1\x8c\xa9\xf1\xf5\xca\xc0\x0e\xf1\x53\xef\x42\xbd\xdb\x59\x7a\xe1\xf3\x4e\x80\x1d\xb6\xf8\xee\x6c\x6d\xd3\x1a\x65\x68\x4e\xdf\x3c\x90\x7a\xe0\xa4\x64\x1d\xd5\x1b\x94\xbc\x6b\xc0\x8d\xd3\x94\xbe\xbd\x18\x48\x21\x2f\x39\x5a\x52\x58\x7c\xcc\x81\x8f\xdd\xcb\x00\xf8\x41\x27\xc8\x58\x8a\x7f\xdd\x3d\x20\x8f\x65\xea\x1d\x37\x77\x64\x1a\x0e\x97\xdd\xeb\x00\x12\xf3\x62\x00\xa7\x61\x2c\x28\xc9\xc1\xbb\x83\xd1\xef\x1a\x48\x49\x3c\x85\xff\xab\x03\xaf\xe3\xa6\x75\x7a\x10\x90\x56\x04\xbc\x2a\x18\x2b\xd8\xf8\xa4\x20\xe1\xbc\x37\xb1\x06\x9a\x2f\xb8\xb1\xda\x56\xf4\x42\x93\x90\x92\xd8\x7a\xbb\xe0\x62\xa9\xea\x9f\x52\x6b\x56\x13\x80\x1e\x50\x92\xdc\xa7\x5f\x9d\xc9\xdd\xc5\x65\xb0\x2b\xab\x66\xa5\x6e\xf9\xde\x64\x8b\xd5\x22\xae\xf4\xe9\x57\xcb\x0d\x93\xa5\x17\xa0\x70\x16\x0b\xc9\x76\x76\xc1\x2d\x6b\x60\x40\x09\x60\x0f\x01\x74\xb3\x00\x8c\x9b\x00\x24\xa2\x4f\x22\xd0\x59\x5f\x03\x8c\xb2\xe8\x61\x5b\x28\x79\xef\x5c\x00\x16\xb1\xe4\x75\xc9\xac\xac\x2c\x01\x3e\x6f\x65\x69\x8b\x61\x4a\x6f\x4a\xc3\xac\xfd\x25\x01\x64\x9a\x9e\x1a\xe0\x16\x86\xc2\xda\xb7\x0a\x72\xb7\x0f\x86\xa2\xf5\xe4\xc9\x9c\x2c\xd1\x95\xb4\x94\x87\x99\xab\x46\x9e\x33\xb2\xa1\xf0\x14\xf3\xb6\x47\x92\xff\xd7\x0d\x85\x83\x0b\x06\x3c\x2c\x2b\x8e\x79\x43\xd6\xe2\x1e\x7b\x2e\x8f\xec\xda\x78\x2f\x4d\xc9\x23\xd9\x08\x14\xe8\xdd\xfe\xc5\x0d\xd6\x77\xe2\x2a\xd9\x56\xf7\x11\x05\x80\x47\x33\x00\xb2\xb6\xde\x4a\xc3\x98\x6f\x67\x17\x78\x22\x24\xb7\xac\x70\xbb\xb1\x95\xdc\xf4\x5e\x90\x96\x13\x00\x3f\x23\xc4\x15\x46\x80\x1c\x8c\x16\x43\x45\x87\x07\x3d\xb3\x8b\xa2\x8d\x40\x14\x21\xc0\x05\x77\x94\x0e\xd3\x1e\xb5\x18\x94\x64\x60\xc6\xb7\x4e\x73\xf7\x31\x80\xa4\x43\x1c\x26\x03\x86\x07\x98\xb5\x9f\xb5\xa2\xdd\xa4\x77\x8b\xb7\x13\x9e\x88\x45\xeb\x29\x90\x1e\xee\x7c\x4c\xc3\x5c\xa7\xe8\x38\x0e\xe4\xee\x43\x3f\xde\xce\xbb\x36\x55\xd4\x0b\xbb\xac\xbf\x71\x59\x49\x2a\xb6\xe4\x01\xdd\x2b\x4f\xf6\x7f\x01\x29\x6a\x88\xa7\x4b\x50\x42\x64\x00\x00\x00\x00\x49\x45\x4e\x44\xae\x42\x60\x82"),
		},
		"/font_8x16.png": &vfsgen_FileInfo{
			name:    "font_8x16.png",
			modTime: time.Date(2018, 2, 7, 10, 21, 22, 191273714, time.UTC),
			content: []byte("\x89\x50\x4e\x47\x0d\x0a\x1a\x0a\x00\x00\x00\x0d\x49\x48\x44\x52\x00\x00\x01\x00\x00\x00\x00\x80\x01\x03\x00\x00\x00\xde\x7b\x25\x49\x00\x00\x00\x06\x50\x4c\x54\x45\x00\x00\x00\xff\xff\xff\xa5\xd9\x9f\xdd\x00\x00\x00\x09\x70\x48\x59\x73\x00\x00\x0b\x13\x00\x00\x0b\x13\x01\x00\x9a\x9c\x18\x00\x00\x00\x07\x74\x49\x4d\x45\x07\xe1\x0b\x02\x0c\x04\x29\xe3\xd1\x3f\x11\x00\x00\x06\x0d\x49\x44\x41\x54\x58\xc3\xe5\x97\xcf\x6b\x24\x37\x16\xc7\xc5\x2c\xd4\x49\x78\x9c\xdb\x63\x6d\xec\xcb\xfc\x01\x22\x81\x9a\x62\x57\xd8\x10\xf6\x5f\xd8\xbb\x98\x09\xda\x1c\x44\xe2\x53\xa5\x60\x45\xcd\xec\x65\xfe\x84\xc0\x90\x73\x2e\xf9\x1f\x02\x15\x17\x88\x3d\x08\x1f\x43\x43\x75\x06\x9f\xda\x97\x25\x34\x04\x9c\x3e\x18\xd5\x7e\x9f\xaa\xdb\x3f\x66\x7a\x3c\x64\x92\x40\x20\xea\x76\x59\xd5\xfd\xa9\xa7\xf7\x4b\xef\xa9\x85\x58\x8f\x11\xaf\x7b\xc7\x06\x78\xfe\x00\x17\x7f\xfb\x9b\xb6\xdd\x00\x07\xfa\xf8\x99\xe8\x0a\xb2\xcf\x82\x20\xa2\x0d\xf0\x1f\x3c\x89\x3b\x00\x3b\xf6\xa3\xa7\x74\xbe\xa3\xed\x60\x84\xd6\xd7\xc0\xd7\x83\xdb\xd5\x9a\x81\xbf\xda\xe3\x67\xb4\x3c\x68\xed\x50\x41\x6e\x06\x76\x13\x4b\x48\x95\x6e\xc5\xa8\xfb\x0f\xad\x7a\x3a\xac\x8e\xc8\x0e\x8e\x97\x20\x25\x1e\x55\x13\xe0\x2f\x46\xba\xb0\x2f\x1b\x00\x3a\x25\xb2\xff\xce\x3a\x48\xd3\xd9\x0a\xea\x7e\xdb\xa7\x74\x81\xe7\x3f\xfe\x36\x6a\xf5\xf4\x82\x25\xec\x65\x20\xa5\x6e\xf4\x00\x5e\x5e\x40\x42\x9b\x01\x82\x04\xd6\x61\xcf\xa5\x76\x92\xe0\x2b\x5e\xc2\x57\xe0\xb1\x44\x6c\x4f\xfe\x35\xc0\x0a\xb1\x57\x25\x4d\x2d\xa9\xf4\x28\x65\xa0\xda\x65\x33\x75\x1f\x69\x79\x41\xec\x87\x3d\x99\x88\x34\x5b\x91\x76\xd9\x51\xbb\x62\x32\xb3\xa1\xf3\x05\xc1\x93\xd0\x81\x97\xd8\xea\xea\x6e\xba\xf1\xf7\xc4\xe2\x9d\xc1\x7a\x3f\x80\xee\x7f\xd6\x32\xa0\xee\x79\xdc\x42\xe9\x4a\xc9\x35\xa2\xc9\x7b\x99\xaa\xe4\xfd\xc6\x1e\x6d\x5d\x10\x4e\x6d\xec\xb6\x55\x08\xfb\x9d\x09\x21\x08\x51\x08\x83\xab\x7e\xe4\x4e\x4f\x9d\x51\x32\x7f\xff\xa0\x6f\x8a\x42\x77\x5d\x11\x02\xa2\x21\x14\x03\x22\x75\xa1\x12\x4a\x5a\x96\x51\xf4\x24\x0b\xb7\x06\x10\x0f\xc9\x26\x38\x2f\x6b\x00\x39\xc9\xe4\x80\x28\xc4\xab\x2b\xe9\x5b\xd6\x5c\x52\x06\x0a\x9a\x03\x18\x5b\x91\x04\x0d\xa4\x8a\x54\x04\x0a\x85\x10\x46\x14\x6b\x40\xc5\x8d\x04\xd5\x93\x29\x64\x11\x14\x03\xaa\x95\x59\xf5\xf4\xc2\xc4\xac\x03\xee\x4c\x8f\x80\x4f\x00\x67\x3e\x65\x37\xba\x10\x22\x67\x3e\xa6\xd4\x59\x0a\x41\x06\x00\x32\x5b\x31\x2d\xe1\x5f\xd4\x82\x1d\x05\xe0\xb9\x6e\x93\x3f\xf0\x5e\xf9\x46\x28\xf8\x81\xee\xc4\x42\xbd\x67\xb0\x7e\xdf\xb1\xfb\x2e\xa0\x62\xe8\x4a\xaf\x52\xd2\x41\x1f\x2c\x96\x7d\xf0\x57\x78\x8d\xa1\xc7\x18\xb5\xd0\x0e\x19\x51\x59\xeb\xac\xb5\x81\xa4\x35\x17\x8b\x60\xf1\x1a\x32\xd0\xab\xe7\x12\xd1\x0e\xce\x9e\xda\x27\x4f\x4e\x33\x30\x5e\x66\xe0\x25\x00\xdb\xbf\x50\x1d\xe7\x48\x08\xb6\xb3\x9f\x7e\xda\x01\x70\x66\x4c\x0c\x18\x02\xa0\xad\x54\xe7\x0c\xbc\x0a\xbe\xb3\x4d\xd3\x25\x92\x8d\x19\x70\x17\xb0\x97\x01\x90\x26\x75\x92\x81\x94\x25\xbc\x0a\x0c\xf4\x67\xc1\x04\x27\x01\x20\x6f\x94\xaa\x26\x09\x00\x0c\x92\x94\xa2\x33\x7d\x00\x60\x0b\x06\x34\x19\xb5\xcf\xc0\x3c\xb0\x92\x19\xb0\x4f\x18\xf8\x1e\x06\x05\x3b\x5a\xfa\x4e\xed\x30\xd0\x05\x36\xd3\x58\x06\x2c\x03\xaf\x18\xd0\xb6\xa7\x5e\x15\x0c\xf8\xc0\x8e\x5a\xfe\x2d\xe8\x66\x91\xe0\xa8\xa5\x5f\x78\xed\xc9\xf6\x7a\xd4\x0f\xf4\x36\xef\xca\x77\xb9\x7f\xe7\xb5\x4d\xff\xeb\x87\xfa\xd5\x00\x89\x73\xb1\x8f\xc4\x38\xa7\xe2\xbc\xba\x9b\x47\x3b\x74\x52\xf3\x7f\x03\xcb\x1c\x76\x40\x61\x5e\x4b\x7f\x6c\xad\xf9\x06\xf8\x04\x57\xb1\x05\x60\x61\x4d\x03\xcf\x99\xda\x55\x3b\x96\x16\x73\x3f\xaf\xe7\xfe\x2a\x22\x9f\x42\x62\x00\xab\x4a\x17\x5c\x58\xc6\x9a\x0a\x47\x23\xd2\x25\xd6\x41\x01\xb0\x21\x9e\xc0\x69\x9c\x94\xb6\x8b\xc9\x44\x4b\x45\x43\x03\x03\xd6\x30\xa0\xb9\xc2\x90\x40\x52\x46\x00\xdd\x6d\xc0\x54\x00\x06\x0a\xea\x0d\xc0\xad\x01\xa9\xa2\x1d\x74\x30\xd7\x40\x88\x21\x03\x76\x0d\x84\x32\x6a\x68\xc3\x6b\xa0\x28\x89\xda\xfb\x9a\xb3\x40\x17\x0b\x3d\x58\xef\x31\xdf\xaf\x91\x0f\x6d\x82\xa3\xee\x64\x81\xcd\x13\xb3\xc9\x87\xe2\xae\xc7\xe3\x1b\xc0\x6b\x89\xd3\x70\x45\xe5\xb1\x3c\x58\x7f\xb2\x7a\xcf\xf4\xa8\xb6\x94\x23\xb9\x2b\x4c\x85\xa8\x18\x41\xd8\x56\x8a\xa3\x6d\x94\xc1\xa6\xa5\x4a\x5c\xed\x08\x1d\xa9\x8a\xca\x89\x2a\x28\xab\x95\xa8\x8c\x38\xc2\xb4\xe1\x32\x4f\xae\xb7\x7b\xd0\x1d\xdf\x82\xc5\x45\x58\xe2\x54\x71\x98\x46\xc2\x86\x69\x3f\xb1\xe8\x00\xa7\xeb\x46\xb3\x56\x01\xe7\x81\x38\x4d\x43\xe8\x0d\x76\x90\xe8\xa2\x6f\x30\xd8\x85\x55\x55\x39\x67\x3f\x8b\x98\xc6\x88\x96\xd3\x2d\xe9\x09\x03\x41\x62\x70\x0f\x82\xef\xb1\x37\xff\x9e\x30\x9d\x00\x33\xda\x16\x00\xfa\x97\xf7\x26\xa5\x94\x01\xbf\x17\x6f\x00\xfa\x1c\x3a\xc4\x2e\x62\xd8\xae\xeb\xf2\xa1\xc2\xb4\xd7\x40\x0f\x09\xd0\x7c\x02\xf4\x04\x60\x89\xd9\x35\xd0\x1a\x02\xa0\x63\x60\x40\x6e\x74\xb0\xf3\x6b\x80\x16\x0c\xc8\xda\xd7\x18\x05\xf4\xd0\x5a\x87\x90\xbe\x38\xc3\xb4\xae\x5b\xef\xe9\x8a\x7e\xa2\x4d\xd8\xdf\x28\x15\xeb\xcf\x67\x5b\xcf\x12\xb7\xd2\xe1\xe4\xad\x61\x6d\x7e\x59\x16\xd4\x37\xd3\x0f\xfe\xf9\x03\xac\x29\x85\x28\x4b\xbc\xb9\xe9\xa1\xfd\x89\xb9\xe6\x9d\xdf\xf1\x11\xe7\x1f\xdf\x7c\x71\x17\x50\x28\xb9\x35\x42\xac\x18\xa0\x2d\x12\x90\x3c\x66\x1e\x32\x70\x7a\x4a\x5b\x24\xe4\xb1\x38\xaa\x72\x06\x6c\xd3\xa1\xa9\x7c\x9c\x5f\xe6\xa2\x11\x23\x95\xb3\x2c\x61\x55\x8a\xd5\x65\x99\x2e\x4b\xe4\x3e\x6a\x7f\xb4\xa9\xf5\x2a\x87\xdb\xb9\x8d\x04\x2a\xca\xa2\x60\x09\x9e\x81\x57\xd0\xa5\x2b\x94\xa2\x59\xc9\x12\x56\xab\xcb\x04\x09\x97\x29\xad\x56\x22\x32\x70\xc6\x26\x14\xc6\xea\x6b\x09\x25\xde\x65\x2e\x82\x19\xe0\xe6\xd6\x15\x67\x67\x7a\xad\xc3\x5b\x80\xaf\xbe\xd4\xd7\x56\xdc\x00\xb5\xf6\x35\x03\x88\x6a\x71\x74\xe3\x87\x1b\xe0\x76\x84\x0b\xb1\x45\xc2\xad\x71\x58\x88\xfb\x25\xdc\x8e\xe6\xdb\x81\x37\x24\xe4\xb3\x10\x7f\x22\xf2\x2b\x47\x81\xf8\x8f\x7f\x0f\x2c\x1f\x8e\x7f\x20\xe0\xb0\x7c\x7c\xfc\xf3\xf8\x78\xfc\x79\x2c\x47\x51\x1e\x1e\x8a\x72\x7c\x5d\x82\x52\x53\xa7\x29\xa7\xa7\x6f\x24\x1c\x8e\x23\xbf\x1f\x1f\xe3\x71\x96\x80\x71\x7c\x78\x78\x3c\x8e\x2b\x7c\xbc\x7c\x38\xb5\x26\x96\xc0\x0a\xf0\x45\x50\xfe\x83\x1e\xf4\x27\x03\xde\xe7\x0c\x4b\xf9\x24\xf9\x70\x96\x2b\x54\x93\x36\xcd\x06\xdd\x6a\x7f\x7d\x70\x10\xdc\xb4\xa5\x9b\xe5\x7e\x17\xb8\x7b\x8a\xb6\x42\x89\xf8\x8b\xc2\xf6\x50\x72\x6f\x03\xa8\x09\x48\x41\xd8\x9a\x5c\x40\xd1\x34\xf9\x50\xb0\xb7\x39\x17\x18\xf0\x75\xec\x9c\x69\xed\x1c\x45\x56\xb6\x2d\x7e\x91\x90\xe4\xae\xce\x65\x4d\xba\xff\x02\x98\xcf\x3a\xa7\x66\x96\xd0\xab\x8f\x86\xc1\x07\xd1\xe2\xa8\x82\xb2\x96\x81\x15\x80\x19\x24\x10\x03\xc9\xd9\x61\xe0\x25\x58\x42\xcb\x96\xfc\x2f\xef\xd8\x59\xd8\x48\x00\xf0\x93\xc9\x3f\x45\x69\x86\xca\x89\x72\x73\x0d\x98\x99\x27\x0d\x20\xeb\x00\x2b\x68\x96\x8f\x3f\x13\x30\x07\x10\x66\x86\xcb\x93\x45\x4d\x0c\x6c\x32\x80\x75\x07\x58\x5b\x91\x4e\x0c\xb5\xd5\x8f\x5a\x74\xfb\x41\x8c\xf8\x45\x3d\xd5\xf8\xfd\x5b\xbe\xed\x7e\x71\x28\x7e\x43\xe0\xff\x48\xde\xdd\x1e\xe2\x83\xa6\x91\x00\x00\x00\x00\x49\x45\x4e\x44\xae\x42\x60\x82"),
		},
		"/tiles.json": &vfsgen_CompressedFileInfo{
			name:             "tiles.json",
			modTime:          time.Date(2018, 2, 21, 12, 1, 49, 44791818, time.UTC),
			uncompressedSize: 719,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x9c\x91\xd1\x6a\x84\x30\x10\x45\xdf\xfd\x8a\x90\x67\x29\xcd\xba\xa6\xd6\x9f\x91\x59\x8d\x1a\x88\x89\x24\x23\xbb\x50\xfc\xf7\x32\xb1\xb4\xec\xd0\x2d\x65\x5f\xef\x49\xee\xc9\x4c\x3e\x0a\x21\x84\x90\x83\x19\x61\x73\x28\x5b\x71\x04\x39\x4c\x61\x8b\xbd\x91\xad\x90\x63\xf0\xd8\x35\x37\xa5\x5f\x56\x3f\xc9\xf2\xe7\x08\x5a\x67\xba\xab\x1d\x70\x96\xad\x68\x38\x98\x8d\x9d\x66\x2a\x55\x3a\x93\xbd\x2c\x0e\x1d\xa4\xde\xda\x47\x32\xa5\x6f\x4a\x77\x69\xe9\xf2\xa9\x3f\x8d\x4a\xff\x5b\x49\x34\x3d\x52\x66\x68\xb0\x53\x17\x8b\xcf\x0b\xbf\x54\x99\x2d\xb0\xae\xd6\x4f\x77\xc2\x0c\x7c\xf0\xa4\xd4\x55\x79\x9f\x5f\x01\x4d\x94\xad\x38\x37\xbf\x81\x13\x09\xde\x18\xc1\x68\xa8\x8a\xc7\xab\xcd\x86\x13\x2f\x5a\xc1\x2d\x54\x53\xb3\xbc\x87\x1e\x37\x5a\x4d\xad\xf9\x0d\x07\x9e\x46\xe3\xf9\x65\x4b\xb4\x8c\x8a\x9b\x13\x1e\xc3\x55\xfc\xc2\xe8\x42\xc8\xc3\xf1\xa9\x07\x03\x4b\xf0\x44\xce\xfc\x55\x11\x90\xba\xea\x77\x2e\x87\x18\x8d\xa3\x57\xbd\x7e\x83\xfd\xf8\xeb\x62\xff\x0c\x00\x00\xff\xff\x42\x8c\x5a\x31\xcf\x02\x00\x00"),
		},
		"/tileset_1bit.png": &vfsgen_FileInfo{
			name:    "tileset_1bit.png",
			modTime: time.Date(2018, 2, 7, 10, 21, 22, 191273714, time.UTC),
			content: []byte("\x89\x50\x4e\x47\x0d\x0a\x1a\x0a\x00\x00\x00\x0d\x49\x48\x44\x52\x00\x00\x00\x80\x00\x00\x00\x80\x01\x03\x00\x00\x00\xf9\xf0\xf3\x88\x00\x00\x00\x06\x50\x4c\x54\x45\x00\x00\x00\xff\xff\xff\xa5\xd9\x9f\xdd\x00\x00\x00\x09\x70\x48\x59\x73\x00\x00\x0b\x13\x00\x00\x0b\x13\x01\x00\x9a\x9c\x18\x00\x00\x00\x07\x74\x49\x4d\x45\x07\xe1\x0a\x04\x0c\x0c\x30\x57\x68\x6b\xb5\x00\x00\x00\x1d\x69\x54\x58\x74\x43\x6f\x6d\x6d\x65\x6e\x74\x00\x00\x00\x00\x00\x43\x72\x65\x61\x74\x65\x64\x20\x77\x69\x74\x68\x20\x47\x49\x4d\x50\x64\x2e\x65\x07\x00\x00\x06\x1c\x49\x44\x41\x54\x48\xc7\x4d\x95\x6d\x6c\x53\x55\x18\xc7\x9f\xdb\xdb\xae\x2f\xdb\x68\xbb\x71\xd9\x20\x63\xeb\xba\x4d\xc9\xc0\xd0\x81\x0b\x95\xa8\xeb\x62\x88\x63\x44\xd2\x2d\x3d\x50\x93\xea\x16\x35\xc1\x64\x24\xcc\x48\xc4\x0f\x84\x9e\x10\x70\x97\xa1\xc2\x07\x9a\xb0\x4f\x56\xe1\x43\x39\x1b\x69\x43\x86\x43\x1c\xf4\x46\x19\xbe\x30\xcc\x3e\xb8\xb8\x44\x42\xaf\x41\x01\x89\xc4\xaa\xbc\x34\xb6\xdb\xf1\x39\xf7\x0e\xe4\xdc\xdd\x2d\xfb\xdf\xdf\xf9\x3f\xcf\x73\x5e\x21\xbe\xc0\xb9\xf9\xc6\x39\x80\xdc\x0f\x54\x02\x6a\xbc\x32\x60\x73\x78\x20\x69\x03\xca\x46\x49\x22\xcb\x79\x16\x68\x55\xd0\x20\x12\xdb\x19\x4b\x3f\x14\x42\xdd\x24\x0a\x9c\xb3\x51\xfc\x8c\x3f\x1e\xf0\x15\x50\xa0\x90\xd8\x4e\x21\x8b\x96\x0e\x50\xf6\x19\x02\x1b\x4d\x90\x2c\x8f\x73\x07\x74\xbd\xb6\x48\x30\x26\x04\xb0\x84\xaa\x9e\xf0\x88\x73\xc9\x7a\x8c\x2f\x46\x01\x2a\x3c\x2c\x70\xb2\x68\x08\x8f\xf2\xb0\xc0\x4b\x5f\x3e\xce\xc3\x48\x0c\xd6\x3e\x78\xec\x21\x12\xf3\x59\xab\x0b\x98\xfa\xff\x79\x78\xac\xb2\x66\x84\xcd\xf2\x04\xe1\x22\x2c\xd8\x75\x88\x2f\x50\x7c\x19\x33\xcb\x77\xe7\xc1\x6c\x1c\xff\x01\x8a\x1d\xea\x17\x05\x6b\x0b\x80\x44\x2d\xa1\x17\xe3\xf8\x57\xef\xe8\xa0\x95\xbb\xb0\x3f\x75\x97\xe2\x05\xec\xb6\x77\xf6\x3e\x6d\x16\xfd\xa9\xc3\xf2\x0d\x7a\x6c\x4c\x33\xa6\x0b\xc3\x06\xea\xfa\xdd\x8d\xc4\x7d\x42\xc8\xfb\x5e\x05\x20\x88\x44\xfd\x24\x00\xc3\x36\x56\x1d\x17\x51\x3a\x4a\xc1\x9b\xe0\x9e\x25\x3d\xe1\x6d\x4a\x9f\x10\x42\x9e\xc0\x67\xb0\x22\x29\x10\xce\x6f\x61\x4b\xcc\xf8\x74\x50\x48\x84\x90\x32\x6b\x4b\x33\x3e\x21\x9c\x15\x68\xc5\xef\xc3\xac\x72\x57\xcb\x8e\x96\x1d\x1d\x79\xd0\xe1\x59\x05\xa3\x6c\x6b\xe6\xcd\xb7\x9a\x6f\xc9\x9a\x84\x44\x0a\x99\xb1\x38\x17\x8f\x25\x24\xf5\x43\x3b\x09\x2b\x64\xa9\x57\x11\x4f\x95\x66\x49\xc2\x46\x11\xe4\xc3\xea\xb8\x78\x5c\x5d\xf6\x10\x12\xd8\x22\x4a\x9f\x78\x5c\x33\x6e\xdd\x88\xe2\x1a\xc3\xa2\x43\x10\x32\x46\xc2\x20\x62\x00\xe5\x6f\xdb\xc2\x92\x26\x53\xf4\x18\x66\xa9\x94\xaa\x7e\x14\x5d\xe6\x2b\x0b\x58\x28\xac\xb4\x91\x4e\x12\x55\xd5\x2f\xb6\x2e\x69\x76\x76\xa1\xe0\x12\x51\xce\x71\xfe\xdd\x66\xb9\xa5\xfc\x76\x7d\x1e\x2a\x09\xd9\x44\x06\x55\xf5\x50\x57\xd3\x1b\x4a\xc4\xe3\x81\x15\x2c\xe5\xda\x9f\x14\xc2\xfa\xaf\x6a\x76\x23\xa1\x74\x62\x94\x7e\x55\x3d\x7c\xdb\x79\x79\xe9\x52\x0b\x35\x33\x4d\x0a\xe1\x95\x6f\x6b\x76\xa3\x10\x20\x36\x62\x1b\x14\xc2\x86\xbf\x95\x08\x0a\x98\xa9\x85\x65\x84\xf0\xf5\x74\xf5\x5b\x4b\x74\x68\x97\xd0\x23\x2a\x4c\x25\xcd\x88\xd2\xca\x70\x50\x53\xa6\x50\xfd\xab\x12\x36\x6a\xe9\x8c\x2d\x12\xa4\xea\x8c\xf0\x48\xb2\x0a\x8e\xcb\x43\xd6\x7a\x3f\x4f\x44\x40\x89\xe0\xa0\x46\x04\x61\xd7\x37\xbc\xd7\x3b\x8e\x99\x8a\x99\x13\xeb\x23\x7d\x35\xfd\x43\xfa\xea\x22\x81\x02\x3d\x76\xfe\xd8\xb9\xc4\x79\xa8\x45\x93\xd4\x1a\xce\x25\x9a\x1e\x49\x1f\x1f\x18\xc1\x99\xc3\x54\xa3\x48\xe8\xe3\x6c\x3c\xd5\xce\xd0\x63\x3f\x73\x65\x38\x97\xe9\xf4\xc8\xf4\x71\xc7\x88\x39\xa6\xfd\xd8\x45\x1b\x9a\x18\xca\xc0\x04\xb4\xee\x17\xd5\x72\x6e\xd7\x77\xde\xd8\xa9\xc3\x0d\x24\x30\x8c\x20\xf4\x40\x2e\x00\xd0\x07\xab\xc4\xc4\xe8\x48\x14\xa2\x53\x61\x80\x4d\xc2\xa3\x67\xbb\xc6\xb9\x33\xbf\xe1\x0a\x12\x5b\xa0\x9e\x5d\x64\x0c\x57\xb2\xbb\xd4\x6b\x12\x36\xf2\x2e\xd9\x87\x42\xc3\x42\xdb\x16\x83\x28\x63\x1a\xd7\x50\x90\x0b\x8b\x84\xfc\x13\xb5\x83\x49\x0c\x20\xb1\x1a\x60\x9f\xb9\x1b\x9c\xd3\xdd\xb3\x38\xfd\x2b\x1f\x6f\x0f\x49\x3b\x38\x4c\xc5\xc6\x34\x57\x29\x16\x87\x9f\x75\x5c\x26\x20\x4a\x39\xbe\x5c\x55\x11\xf3\x1d\x15\x44\x12\x8b\xc9\x54\xaa\xaa\xa7\x16\xfa\x71\x21\x55\x22\xd1\xd3\xa3\xdb\xb1\x4b\x33\x8e\x08\x80\xcb\xf0\xd0\x64\xbf\xbf\x62\x8d\x35\x33\x0d\xb0\x1c\xc4\xe4\x53\x89\xb1\x65\xcf\x97\x1d\x1d\xa2\xf6\x09\x24\x5c\x8c\x4a\x7e\xff\x56\x70\x7a\x76\xea\x15\x37\x8c\x28\x14\xe7\x17\xb4\xf2\x40\xe0\xba\x27\x07\xcc\x92\x72\x71\xae\xaa\x52\xbf\x32\x18\x9d\xaa\x9d\x32\x08\x3c\x30\x98\x25\xb3\x2e\xb8\xe1\x4a\xcd\x15\x23\x0a\x49\xa8\xea\xba\xc1\xa7\xa2\xbd\x53\x55\x48\x6c\x27\x04\x25\xb2\x39\xdf\x3a\xd9\xb6\xc5\xbb\x05\x4f\x28\x51\x8b\xdf\x1f\x83\x48\xcc\x20\x3a\x5f\xfd\x51\xc1\x83\x8b\x05\x21\x38\xd9\x36\xd0\x36\x00\x0c\x77\x54\x82\xf8\xfd\x10\xb2\xe7\xbb\x67\xbb\x67\x81\xc4\xda\x0d\x0f\x2c\x1c\x0e\x0e\x1f\x1c\x06\x36\xd6\x9a\x92\xa8\x51\xbe\xd9\xc8\xb6\x76\x62\x09\x31\x66\xf5\xe0\x88\x56\xa0\xc0\x56\xdf\x61\x96\x90\xaa\x3a\x1c\x5d\x00\x55\x30\x88\x99\xce\x10\xab\x0f\x17\xc9\x8c\x21\x14\x8c\x4c\xad\x3e\xbf\x7f\xcf\xc3\x20\x40\x1d\xcc\xe1\xcc\x91\x1e\x67\x98\xb1\x17\xae\xd7\x02\xd4\x9a\x44\xca\x91\xf1\xfb\x5f\xbf\x80\x8e\x01\x29\x0f\xb8\xe3\x88\x33\x4c\x48\x43\x49\x08\xb2\x79\xc0\x38\x70\x37\xc4\x79\x19\x40\x3f\xbe\x24\x62\x23\xc2\x63\xfd\xcd\x8e\x52\x87\xa0\xc4\x01\xe3\x44\x62\xcf\xc3\x26\xab\xef\x12\x1e\x74\x84\x84\x0d\x8f\x1d\xff\xac\xbd\xb7\xf6\xde\x4a\xd3\xc3\x89\x51\x5e\xfe\xcb\x67\xf5\x59\x1d\x82\x20\x9d\xce\x5e\xc6\xb6\x3e\x27\x3c\xac\x00\x62\x99\x1e\x1a\xf7\xfb\x1b\xf6\x3e\xaa\x16\x97\x69\xe7\x01\x42\x94\xc8\xa2\x60\xae\x8e\xbc\x9b\xf3\x5c\x11\xf0\x37\x7a\xb8\xf3\xee\xbc\x07\xdd\x50\x70\xe0\xc8\x31\x16\xae\xec\x9b\x0f\x3b\x7b\x12\x28\x14\xfe\x1d\x45\x22\x59\x96\x5d\xc8\xb8\x4e\x7c\xaf\xc9\xa0\x2a\x82\x10\xcb\x85\xf3\x07\x7f\xe6\x8a\x6e\xb5\x41\x10\x9f\x1e\xc8\x96\x28\x7c\x90\xcc\x15\x1d\xaa\x57\x10\x27\xa6\xb2\xf3\x09\xf2\xf1\xd9\x5c\xb1\xa0\x7a\x8b\x48\x9c\xb8\x9c\xc5\x73\x9d\x9d\xd1\x64\xf5\x08\x9e\xbc\x8c\xfd\x32\x24\x3c\x6e\xd7\xe4\x8a\xea\x61\x6f\x03\x12\x8d\x07\x3a\x4a\x40\x5d\x19\x14\xe6\xbd\x5e\x24\xde\x9c\x8a\xcf\x93\x44\x37\x7a\xa8\xa3\x28\x3c\xe9\x71\x64\xd6\x20\xcc\x3c\xee\x4c\xe4\x8a\xc3\x26\x91\xb4\x66\x2f\x51\xf8\xa4\x90\x2b\xce\x13\x83\x48\xda\x2e\xda\x12\x84\xda\x72\x45\x66\x12\x58\x49\x19\xde\x95\x0b\x9a\x6c\x6c\x13\xa0\x62\x54\xdc\x1c\xef\x16\xf3\x96\x91\xb4\x10\xbe\x1e\xb9\x02\x97\x07\xde\x68\x42\x18\x84\xb2\xc0\xd3\xbf\x89\x21\xca\xce\xa3\x50\xa9\x75\x68\x8e\x68\xa8\xbc\x2e\x88\x17\x9d\x29\xd4\x35\x56\xad\xe2\xf7\x9b\xae\xb1\xd3\x86\xe0\xd6\x9c\x0d\x4a\x8c\x36\x04\x26\xb3\xf3\x86\x60\x9f\x73\xd5\xaf\x1b\xcb\xce\x3e\x13\xa3\x92\x21\x48\x73\xca\x5c\x5b\xec\x4c\x5f\xd3\x1f\x66\x16\x20\x15\x94\xa3\xfe\xc9\x93\xb3\x8d\x48\x98\x61\xf3\xab\x02\x6d\xb1\x53\x7d\xbe\xb1\x0b\xd8\x8c\x3c\xa2\xba\x52\x3a\x43\x1a\x63\xa7\xb1\x19\xc2\x3b\x79\x77\xe1\xe4\xcf\xc1\xc9\xf1\xf6\xf1\x76\x43\xc8\x95\xdc\x85\x53\xdb\x6a\xa2\xe9\xbb\xe9\xbb\x28\x38\xb5\xac\x66\xd7\xb3\xd7\xaa\x0b\x67\xc9\x59\x82\x82\x6d\x46\x5c\x1a\xb4\xb6\x72\x90\x9d\x66\xc2\xc3\xa6\x1b\xe7\x47\x5e\xd6\xcc\x3c\xfe\x03\x4b\x9e\x54\xc4\x3c\x0d\x04\x46\x00\x00\x00\x00\x49\x45\x4e\x44\xae\x42\x60\x82"),
		},
	}
	fs["/"].(*vfsgen_DirInfo).entries = []os.FileInfo{
		fs["/16x16_sm_ascii.png"].(os.FileInfo),
		fs["/font_8x16.png"].(os.FileInfo),
		fs["/tiles.json"].(os.FileInfo),
		fs["/tileset_1bit.png"].(os.FileInfo),
	}

	return fs
}()

type vfsgen_FS map[string]interface{}

func (fs vfsgen_FS) Open(path string) (http.File, error) {
	path = pathpkg.Clean("/" + path)
	f, ok := fs[path]
	if !ok {
		return nil, &os.PathError{Op: "open", Path: path, Err: os.ErrNotExist}
	}

	switch f := f.(type) {
	case *vfsgen_CompressedFileInfo:
		gr, err := gzip.NewReader(bytes.NewReader(f.compressedContent))
		if err != nil {
			// This should never happen because we generate the gzip bytes such that they are always valid.
			panic("unexpected error reading own gzip compressed bytes: " + err.Error())
		}
		return &vfsgen_CompressedFile{
			vfsgen_CompressedFileInfo: f,
			gr: gr,
		}, nil
	case *vfsgen_FileInfo:
		return &vfsgen_File{
			vfsgen_FileInfo: f,
			Reader:          bytes.NewReader(f.content),
		}, nil
	case *vfsgen_DirInfo:
		return &vfsgen_Dir{
			vfsgen_DirInfo: f,
		}, nil
	default:
		// This should never happen because we generate only the above types.
		panic(fmt.Sprintf("unexpected type %T", f))
	}
}

// vfsgen_CompressedFileInfo is a static definition of a gzip compressed file.
type vfsgen_CompressedFileInfo struct {
	name              string
	modTime           time.Time
	compressedContent []byte
	uncompressedSize  int64
}

func (f *vfsgen_CompressedFileInfo) Readdir(count int) ([]os.FileInfo, error) {
	return nil, fmt.Errorf("cannot Readdir from file %s", f.name)
}
func (f *vfsgen_CompressedFileInfo) Stat() (os.FileInfo, error) { return f, nil }

func (f *vfsgen_CompressedFileInfo) GzipBytes() []byte {
	return f.compressedContent
}

func (f *vfsgen_CompressedFileInfo) Name() string       { return f.name }
func (f *vfsgen_CompressedFileInfo) Size() int64        { return f.uncompressedSize }
func (f *vfsgen_CompressedFileInfo) Mode() os.FileMode  { return 0444 }
func (f *vfsgen_CompressedFileInfo) ModTime() time.Time { return f.modTime }
func (f *vfsgen_CompressedFileInfo) IsDir() bool        { return false }
func (f *vfsgen_CompressedFileInfo) Sys() interface{}   { return nil }

// vfsgen_CompressedFile is an opened compressedFile instance.
type vfsgen_CompressedFile struct {
	*vfsgen_CompressedFileInfo
	gr      *gzip.Reader
	grPos   int64 // Actual gr uncompressed position.
	seekPos int64 // Seek uncompressed position.
}

func (f *vfsgen_CompressedFile) Read(p []byte) (n int, err error) {
	if f.grPos > f.seekPos {
		// Rewind to beginning.
		err = f.gr.Reset(bytes.NewReader(f.compressedContent))
		if err != nil {
			return 0, err
		}
		f.grPos = 0
	}
	if f.grPos < f.seekPos {
		// Fast-forward.
		_, err = io.CopyN(ioutil.Discard, f.gr, f.seekPos-f.grPos)
		if err != nil {
			return 0, err
		}
		f.grPos = f.seekPos
	}
	n, err = f.gr.Read(p)
	f.grPos += int64(n)
	f.seekPos = f.grPos
	return n, err
}
func (f *vfsgen_CompressedFile) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekStart:
		f.seekPos = 0 + offset
	case io.SeekCurrent:
		f.seekPos += offset
	case io.SeekEnd:
		f.seekPos = f.uncompressedSize + offset
	default:
		panic(fmt.Errorf("invalid whence value: %v", whence))
	}
	return f.seekPos, nil
}
func (f *vfsgen_CompressedFile) Close() error {
	return f.gr.Close()
}

// vfsgen_FileInfo is a static definition of an uncompressed file (because it's not worth gzip compressing).
type vfsgen_FileInfo struct {
	name    string
	modTime time.Time
	content []byte
}

func (f *vfsgen_FileInfo) Readdir(count int) ([]os.FileInfo, error) {
	return nil, fmt.Errorf("cannot Readdir from file %s", f.name)
}
func (f *vfsgen_FileInfo) Stat() (os.FileInfo, error) { return f, nil }

func (f *vfsgen_FileInfo) NotWorthGzipCompressing() {}

func (f *vfsgen_FileInfo) Name() string       { return f.name }
func (f *vfsgen_FileInfo) Size() int64        { return int64(len(f.content)) }
func (f *vfsgen_FileInfo) Mode() os.FileMode  { return 0444 }
func (f *vfsgen_FileInfo) ModTime() time.Time { return f.modTime }
func (f *vfsgen_FileInfo) IsDir() bool        { return false }
func (f *vfsgen_FileInfo) Sys() interface{}   { return nil }

// vfsgen_File is an opened file instance.
type vfsgen_File struct {
	*vfsgen_FileInfo
	*bytes.Reader
}

func (f *vfsgen_File) Close() error {
	return nil
}

// vfsgen_DirInfo is a static definition of a directory.
type vfsgen_DirInfo struct {
	name    string
	modTime time.Time
	entries []os.FileInfo
}

func (d *vfsgen_DirInfo) Read([]byte) (int, error) {
	return 0, fmt.Errorf("cannot Read from directory %s", d.name)
}
func (d *vfsgen_DirInfo) Close() error               { return nil }
func (d *vfsgen_DirInfo) Stat() (os.FileInfo, error) { return d, nil }

func (d *vfsgen_DirInfo) Name() string       { return d.name }
func (d *vfsgen_DirInfo) Size() int64        { return 0 }
func (d *vfsgen_DirInfo) Mode() os.FileMode  { return 0755 | os.ModeDir }
func (d *vfsgen_DirInfo) ModTime() time.Time { return d.modTime }
func (d *vfsgen_DirInfo) IsDir() bool        { return true }
func (d *vfsgen_DirInfo) Sys() interface{}   { return nil }

// vfsgen_Dir is an opened dir instance.
type vfsgen_Dir struct {
	*vfsgen_DirInfo
	pos int // Position within entries for Seek and Readdir.
}

func (d *vfsgen_Dir) Seek(offset int64, whence int) (int64, error) {
	if offset == 0 && whence == io.SeekStart {
		d.pos = 0
		return 0, nil
	}
	return 0, fmt.Errorf("unsupported Seek in directory %s", d.name)
}

func (d *vfsgen_Dir) Readdir(count int) ([]os.FileInfo, error) {
	if d.pos >= len(d.entries) && count > 0 {
		return nil, io.EOF
	}
	if count <= 0 || count > len(d.entries)-d.pos {
		count = len(d.entries) - d.pos
	}
	e := d.entries[d.pos : d.pos+count]
	d.pos += count
	return e, nil
}
