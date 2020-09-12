package utils

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
)

func Zip(contents []byte) ([]byte, error) {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	if _, err := gz.Write(contents); err != nil {
		return nil, err
	}
	if err := gz.Close(); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func Unzip(contents []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewReader(contents))
	if err != nil {
		return nil, err
	}
	result, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return result, nil
}
