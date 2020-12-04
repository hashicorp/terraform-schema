package schema

type coreSchemaRequiredErr struct{}

func (e coreSchemaRequiredErr) Error() string {
	return "core schema required (none provided)"
}
