package distro

type Distro struct {
	Label      string
	PathToHash map[string]string
}

func NewDistro() *Distro {
	return &Distro{PathToHash: map[string]string{}}
}
