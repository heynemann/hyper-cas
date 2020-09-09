package distro

import "encoding/json"

func ParseDistro(data []byte) (*Distro, error) {
	var distro Distro
	err := json.Unmarshal(data, &distro)
	return &distro, err
}
