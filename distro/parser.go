package distro

import (
	"fmt"
	"strings"
)

func ParseDistro(data string) (*Distro, error) {
	distro := NewDistro()

	var label strings.Builder
	var pathBuilder strings.Builder
	var hashBuilder strings.Builder
	isPath := true
	path := ""
	for _, char := range data {
		if distro.Label == "" {
			if char == '\n' {
				distro.Label = label.String()
				fmt.Println(distro.Label)
				continue
			}
			label.WriteRune(char)
		} else {
			if isPath {
				if char == '@' {
					isPath = false
					path = pathBuilder.String()
					pathBuilder.Reset()
				} else {
					pathBuilder.WriteRune(char)
				}
			} else {
				if char == '\n' {
					if isPath {
						return nil, fmt.Errorf("Malformed path %s. Must be in the form of path@hash.", pathBuilder.String())
					}
					isPath = true
					distro.PathToHash[path] = hashBuilder.String()
					hashBuilder.Reset()
				} else {
					hashBuilder.WriteRune(char)
				}
			}
		}
	}
	if distro.Label == "" {
		distro.Label = label.String()
	} else {
		if isPath {
			path = pathBuilder.String()
			if path != "" {
				return nil, fmt.Errorf("Malformed path %s. Must be in the form of path@hash.", pathBuilder.String())
			}
		}
		hash := hashBuilder.String()
		if hash == "" {
			return nil, fmt.Errorf("Malformed path %s. Must be in the form of path@hash.", pathBuilder.String())
		}
		distro.PathToHash[path] = hash
	}
	return distro, nil
}
