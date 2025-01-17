package code

type Options struct {
	ProgramLangType string
}

type Option func(*Options)

func WithProgramLangType(programLangType string) Option {
	return func(o *Options) {
		o.ProgramLangType = programLangType
	}
}
