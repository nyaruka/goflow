package utils

import "strconv"

// VariableResolver defines the interface used by Excellent objects that can be indexed into
type VariableResolver interface {
	Resolve(key string) interface{}
	Default() interface{}
}

// ResolveVariable will resolve the passed in string variable given in dot notation and return
// the value as defined by the VariableResolver passed in.
//
// Example syntaxes:
//      foo.bar.0  - 0th element of bar slice within foo, could also be "0" key in bar map within foo
//      foo.bar[0] - same as above
func ResolveVariable(env Environment, variable interface{}, key string) interface{} {
	var err error

	// self referencing
	if key == "" {
		return variable
	}

	// strip leading '.'
	if key[0] == '.' {
		key = key[1:]
	}

	rest := key
	for rest != "" && variable != nil {
		key, rest = popNextVariable(rest)

		resolver, isResolver := variable.(VariableResolver)

		// look it up in our resolver
		if isResolver {
			variable = resolver.Resolve(key)
			err, isErr := variable.(error)
			if isErr {
				return err
			}

			continue
		}

		// we are a slice
		if IsSlice(variable) {
			idx, err := strconv.Atoi(key)
			if err != nil {
				return err
			}

			variable, err = LookupIndex(variable, idx)
			if err != nil {
				return err
			}
			continue
		}

		// we are a map
		if IsMap(variable) {
			variable, err = LookupKey(variable, key)
			if err != nil {
				return err
			}
			continue
		}
	}

	return variable
}

// popNextVariable pops the next variable off our string:
//     "foo.bar.baz" -> "foo", "bar.baz"
//     "foo[0].bar" -> "foo", "[0].baz"
//     "foo.0.bar" -> "foo", "0.baz"
//     "[0].bar" -> "0", "bar"
func popNextVariable(key string) (string, string) {
	var keyStart = 0
	var keyEnd = -1
	var restStart = -1

	for i, c := range key {
		if i == 0 && c == '[' {
			keyStart++
		} else if c == '[' {
			keyEnd = i
			restStart = i
			break
		} else if c == ']' {
			keyEnd = i
			restStart = i + 1
			break
		} else if c == '.' {
			keyEnd = i
			restStart = i + 1
			break
		}
	}

	if keyEnd == -1 {
		return key, ""
	}

	return key[keyStart:keyEnd], key[restStart:]
}
