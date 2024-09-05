package util

func isValidIdentifierStart(ch string) bool {
	return ch == "_" || (ch >= "a" && ch <= "z") || (ch >= "A" && ch <= "Z")
}

func isValidIdentifier(ch string) bool {
	return isValidIdentifierStart(ch) || (ch >= "0" && ch <= "9")
}

func NameIsValidIdentifier(name string) bool {
	if len(name) == 0 {
		return false
	}
	// postgres 'name' type (default)
	if len(name) > 63 {
		return false
	}
	if !isValidIdentifierStart(string(name[0])) {
		return false
	}
	for i := 1; i < len(name); i++ {
		if !isValidIdentifier(string(name[i])) {
			return false
		}
	}
	return true
}
