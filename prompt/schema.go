package prompt

type Template interface {
	Format(values map[string]any) (string, error)
}
