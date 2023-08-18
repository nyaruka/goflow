package envs

import "strings"

// CollateEquals returns true if the given strings are equal in the given environment's collation
func CollateEquals(env Environment, s, t string) bool {
	return CollateTransform(env, s) == CollateTransform(env, t)
}

// CollateTransform transforms the given string into it's form to be used for collation.
func CollateTransform(env Environment, s string) string {
	return strings.ToLower(s)
}
