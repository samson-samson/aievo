package reader

type Reader interface {
	Read(url string) (string, error)
}

type Factory func() Reader

type ReadParam struct {
	Url string
}
