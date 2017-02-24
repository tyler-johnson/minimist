package minimist

import "strings"

// Argv is the result of parsing command-line arguments.
type Argv map[string]interface{}

// NonFlags are the rest of the args which were not parsed as flags before "--"
func (am Argv) NonFlags() []string {
	return am["_"].([]string)
}

// Unparsed are args that came after "--"
func (am Argv) Unparsed() []string {
	return am["--"].([]string)
}

// Get returns the value at deep key
func (am Argv) Get(key string) interface{} {
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
