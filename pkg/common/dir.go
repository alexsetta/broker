package common

type Dir struct {
	Base   string
	Files  string
	Config string
}

func NewDir(base string) *Dir {
	if base == "" {
		base = "../.."
	}
	return &Dir{
		Base:   base,
		Files:  base + "/files/",
		Config: base + "/config/",
	}
}
