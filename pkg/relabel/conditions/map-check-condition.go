package conditions

type mapCheckCondition struct {
	keys   []string
	values []string
	match  map[string]string
}

func (m *mapCheckCondition) isMatch(_ map[string]string) bool {
	//TODO implement me
	panic("implement me")
}
