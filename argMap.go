package minimist

import "strings"

// ArgMap is the result of parsing command-line arguments.
type ArgMap map[string]interface{}

// NonFlags are the rest of the args which were not parsed as flags before "--"
func (am ArgMap) NonFlags() []string {
	return am["_"].([]string)
}

// Unparsed are args that came after "--"
func (am ArgMap) Unparsed() []string {
	return am["--"].([]string)
}

// Get returns the value at deep key
func (am ArgMap) Get(key string) interface{} {
	keys := strings.Split(key, ".")
	head := keys[:len(keys)-1]
	last := keys[len(keys)-1]
	res := am

	for _, k := range head {
		existing := res[k]
		if ex, ok := existing.(map[string]interface{}); ok {
			res = ex
		} else {
			return nil
		}
	}

	return res[last]
}
