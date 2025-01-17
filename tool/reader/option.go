package reader

type Options struct {
	ReaderType string
}

type Option func(*Options)

func WithReaderType(readerType string) Option {
	return func(o *Options) {
		o.ReaderType = readerType
	}
}
