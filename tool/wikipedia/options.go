package wikipedia

const (
	_defaultTopK         = 2
	_defaultDocMaxChars  = 2000
	_defaultLanguageCode = "en"
	_defaultUserAgent    = "aievo/1.0"
)

type Options struct {
	TopK         int
	DocMaxChars  int
	LanguageCode string
	UserAgent    string
}

type Option func(*Options)

func NewDefaultOptions() *Options {
	return &Options{
		TopK:         _defaultTopK,
		DocMaxChars:  _defaultDocMaxChars,
		LanguageCode: _defaultLanguageCode,
		UserAgent:    _defaultUserAgent,
	}
}

func WithTopK(topK int) Option {
	return func(o *Options) {
		o.TopK = topK
	}
}

func WithDocMaxChars(docMaxChars int) Option {
	return func(o *Options) {
		o.DocMaxChars = docMaxChars
	}
}

func WithLanguageCode(languageCode string) Option {
	return func(o *Options) {
		o.LanguageCode = languageCode
	}
}

func WithUserAgent(userAgent string) Option {
	return func(o *Options) {
		o.UserAgent = userAgent
	}
}
