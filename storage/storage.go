// Copyright (c) 2020 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package storage

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"

	"github.com/iotexproject/phoenix-gem/auth"
)

type (
	// Object is a generic representation of a storage object
	Object struct {
		Meta         Metadata
		Path         string
		Content      []byte
		LastModified time.Time
	}
	// Metadata represents the meta information of the object
	// includes object name , object version , etc...
	Metadata struct {
		Name    string
		Version string
	}

	// Backend is a generic interface for storage backends
	Backend interface {
		CreateBucket(bucket string) (Object, error)
		DeleteBucket(bucket string) error
		ListObjects(bucket, prefix string) ([]Object, error)
		GetObject(bucket, path string) (Object, error)
		PutObject(bucket, path string, content []byte) error
		DeleteObject(bucket, path string) error
	}
)

// HasExtension determines whether or not an object contains a file extension
func (object Object) HasExtension(extension string) bool {
	return filepath.Ext(object.Path) == fmt.Sprintf(".%s", extension)
}

func cleanPrefix(prefix string) string {
	return strings.Trim(prefix, "/")
}

func removePrefixFromObjectPath(prefix string, path string) string {
	if prefix == "" {
		return path
	}
	path = strings.Replace(path, fmt.Sprintf("%s/", prefix), "", 1)
	return path
}

func objectPathIsInvalid(path string) bool {
	return strings.Contains(path, "/") || path == ""
}

func NewStorage(store auth.Store) (Backend, error) {
	var provider Backend
	var err error
	switch store.Name() {
	case "s3", "minio":
		scr := credentials.NewStaticCredentials(
			store.AccessKey(),
			store.AccessToken(),
			"")
		provider = NewAmazonS3BackendWithCredentials("", store.Region(), store.Endpoint(), "", scr)
	default:
		err = fmt.Errorf("storage provider `%s` not supported", store.Name())
	}
	return provider, err

}
