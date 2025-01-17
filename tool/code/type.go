package code

type Runner interface {
	// CheckRuntime check program language runtime
	CheckRuntime() error
	// Run execute code
	// for example: go run code 1 2 3
	Run(code string, args []string) (string, error)
}

type RunnerFactory func() Runner
