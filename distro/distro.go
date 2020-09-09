package distro

type DistroItem struct {
	path string
	hash string
}

type Distro struct {
	Paths  []string
	Hashes []string
}

func NewDistro() *Distro {
	return &Distro{
		Paths:  []string{},
		Hashes: []string{},
	}
}

func (d *Distro) AppendPath(path string, hash string) {
	d.Paths = append(d.Paths, path)
	d.Hashes = append(d.Hashes, hash)
}
