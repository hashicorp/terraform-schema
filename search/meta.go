package search

type Meta struct {
	Path      string
	Filenames []string

	Variables map[string]Variable
	Lists     map[string]List
}
