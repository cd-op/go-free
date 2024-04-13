// Package variables provides utility functions to work with key-value maps such
// as the ones present in the metadata headers of documents and templates.
package variables

// Coalesce merges key-value pairs from all the maps passed as arguments onto a
// single map. The values of keys in earlier maps have priority over later maps.
//
// This function does not operate recursively on nested maps. In other words,
// only the keys at the root level are checked and/or taken.
func Coalesce(maps ...map[string]any) map[string]any {
	ret := map[string]any{}

	for i := len(maps) - 1; i >= 0; i-- {
		for k, v := range maps[i] {
			ret[k] = v
		}
	}

	return ret
}
