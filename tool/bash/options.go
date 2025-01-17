package bash

type Options struct {
	AskHumanInput   bool
	PreExecCommands []string
}

type Option func(*Options)

func NewDefaultOptions() *Options {
	return &Options{
		AskHumanInput: false,
	}
}

func WithAskHumanInput(askHumanInput bool) Option {
	return func(o *Options) {
		o.AskHumanInput = askHumanInput
	}
}

func WithPreExecCommands(preExecCommands []string) Option {
	return func(o *Options) {
		o.PreExecCommands = preExecCommands
	}
}
