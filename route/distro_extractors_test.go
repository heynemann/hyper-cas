package route

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathExtractor(t *testing.T) {
	host := "test.something.com"
	path := "/distro/something.txt"
	headers := func(key string) []byte { return nil }

	extractor := NewPathExtractor()
	distro, path, err := extractor.ExtractDistroAndPath(host, path, headers)

	assert.Nil(t, err)
	assert.Equal(t, "distro", distro)
	assert.Equal(t, "/something.txt", path)
}

func TestPathExtractorWithEmptyPath(t *testing.T) {
	host := "test.something.com"
	path := "/distro"
	headers := func(key string) []byte { return nil }

	extractor := NewPathExtractor()
	distro, path, err := extractor.ExtractDistroAndPath(host, path, headers)

	assert.Nil(t, err)
	assert.Equal(t, "distro", distro)
	assert.Equal(t, "/", path)
}

func TestPathExtractorWithEmptyPath2(t *testing.T) {
	host := "test.something.com"
	path := "/distro/"
	headers := func(key string) []byte { return nil }

	extractor := NewPathExtractor()
	distro, path, err := extractor.ExtractDistroAndPath(host, path, headers)

	assert.Nil(t, err)
	assert.Equal(t, "distro", distro)
	assert.Equal(t, "/", path)
}

func TestPathExtractorFailWithoutPath(t *testing.T) {
	host := "test.something.com"
	path := "/"
	headers := func(key string) []byte { return nil }

	extractor := NewPathExtractor()
	distro, path, err := extractor.ExtractDistroAndPath(host, path, headers)

	assert.NotNil(t, err)
	assert.EqualError(t, err, "Invalid path. Must contain at least distro (/).")
	assert.Equal(t, "", distro)
	assert.Equal(t, "", path)
}

func TestPathExtractorFailWithoutPath2(t *testing.T) {
	host := "test.something.com"
	path := ""
	headers := func(key string) []byte { return nil }

	extractor := NewPathExtractor()
	distro, path, err := extractor.ExtractDistroAndPath(host, path, headers)

	assert.NotNil(t, err)
	assert.EqualError(t, err, "Invalid path. A path is required and can't be empty.")
	assert.Equal(t, "", distro)
	assert.Equal(t, "", path)
}

func TestSubdomainExtractor(t *testing.T) {
	testData := []struct{ host, path, expectedDistro, expectedPath string }{
		{"distro--master.something.com", "/qwe/something.txt", "distro--master", "/qwe/something.txt"},
		{"distro--master.test1.something.com", "/qwe/something.txt", "distro--master.test1", "/qwe/something.txt"},
		{"distro--master.test1.something.com", "/", "distro--master.test1", "/"},
		{"distro--master.test1.something.com", "", "distro--master.test1", "/"},
	}
	headers := func(key string) []byte { return nil }

	for _, test := range testData {
		extractor := NewSubdomainExtractor()
		distro, path, err := extractor.ExtractDistroAndPath(test.host, test.path, headers)

		assert.Nil(t, err)
		assert.Equal(t, test.expectedDistro, distro)
		assert.Equal(t, test.expectedPath, path)
	}
}

func TestSubdomainExtractorWithEmptySubdomain(t *testing.T) {
	host := "something.com"
	path := "/test.txt"
	headers := func(key string) []byte { return nil }

	extractor := NewSubdomainExtractor()
	distro, path, err := extractor.ExtractDistroAndPath(host, path, headers)

	assert.NotNil(t, err)
	assert.EqualError(t, err, "Could not extract distribution from subdomain (Host: something.com).")
	assert.Equal(t, "", distro)
	assert.Equal(t, "", path)

}
