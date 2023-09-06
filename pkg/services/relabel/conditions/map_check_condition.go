package conditions

type mapCheckCondition struct {
	keys   []string
	values []string
	match  map[string]string
}

func (m *mapCheckCondition) isMatch(check map[string]string) bool {
	if m.keys != nil {
		for _, matchKey := range m.keys {
			if _, ok := check[matchKey]; !ok {
				return false
			}
		}
	}

	if m.values != nil {
		for _, matchValue := range m.values {
			for _, checkValue := range check {
				if matchValue == checkValue {
					goto found
				}
			}

			return false

		found:
		}
	}

	if m.match != nil {
		for matchKey, matchValue := range m.match {
			if checkValue, ok := check[matchKey]; ok {
				if matchValue != checkValue {
					return false
				}
			} else {
				return false
			}
		}
	}

	return true
}
