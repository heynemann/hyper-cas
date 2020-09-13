package route

import (
	"fmt"
	"strings"

	"golang.org/x/net/publicsuffix"
)

type DistroExtractor interface {
	ExtractDistroAndPath(host, path string, header func(string) []byte) (string, string, error)
}

type PathExtractor struct {
}

func NewPathExtractor() *PathExtractor {
	return &PathExtractor{}
}

func (e *PathExtractor) ExtractDistroAndPath(host, path string, header func(string) []byte) (string, string, error) {
	if path == "" {
		return "", "", fmt.Errorf("Invalid path. A path is required and can't be empty.")
	}
	var sb strings.Builder
	distro := ""
	for _, char := range path[1:] {
		if distro == "" && char == '/' {
			distro = sb.String()
			sb.Reset()
		}
		sb.WriteRune(char)
	}
	if distro == "" {
		distro := sb.String()
		path := "/"
		if distro == "" {
			return "", "", fmt.Errorf("Invalid path. Must contain at least distro (%s).", path)
		}
		return distro, path, nil
	}
	return distro, sb.String(), nil
}

type SubdomainExtractor struct{}

func NewSubdomainExtractor() *SubdomainExtractor {
	return &SubdomainExtractor{}
}

func (e *SubdomainExtractor) ExtractDistroAndPath(host, path string, header func(string) []byte) (string, string, error) {
	suffix, _ := publicsuffix.PublicSuffix(host)
	sub := host[:len(host)-len(suffix)-1]
	values := strings.Split(sub, ".")
	if path == "" {
		path = "/"
	}
	distro := strings.Join(values[:len(values)-1], ".")
	if distro == "" {
		return "", "", fmt.Errorf("Could not extract distribution from subdomain (Host: %s).", host)
	}
	return distro, path, nil
}
