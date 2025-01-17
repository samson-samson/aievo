package arxiv

type Options struct {
	Topk      int
	UserAgent string
}

type Option func(*Options)

func WithTopk(topk int) Option {
	return func(o *Options) {
		o.Topk = topk
	}
}

func WithUserAgent(userAgent string) Option {
	return func(o *Options) {
		o.UserAgent = userAgent
	}
}
