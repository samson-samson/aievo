package file

type Options struct {
	Workspace string
}

type Option func(*Options)

func WithWorkspace(workspace string) Option {
	return func(o *Options) {
		o.Workspace = workspace
	}
}
