package search

type Options struct {
	TopK   int
	Engine string
	ApiKey string
}

type Option func(*Options)

func WithTopK(topk int) Option {
	return func(o *Options) {
		o.TopK = topk
	}
}

func WithEngine(engine string) Option {
	return func(o *Options) {
		o.Engine = engine
	}
}

func WithApiKey(apiKey string) Option {
	return func(o *Options) {
		o.ApiKey = apiKey
	}
}
