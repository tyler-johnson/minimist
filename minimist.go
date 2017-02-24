package minimist

import (
	"os"
	"regexp"
	"strings"
	//"strings"
)

func nextString(list []string, i int) *string {
	if i+1 < len(list) {
		return &list[i+1]
	}
	return nil
}

func sliceContains(slice []string, needle string) bool {
	if slice == nil {
		return false
	}
	for _, s := range slice {
		if s == needle {
			return true
		}
	}
	return false
}

var integerRe = regexp.MustCompile(`^-?\d+$`)
var numberRe = regexp.MustCompile(`^-?\d+(\.\d+)?(e-?\d+)?$`)

// --port=8000
var longFormEqualRe = regexp.MustCompile(`^--.+=`)
var longFormEqualValsRe = regexp.MustCompile(`^--([^=]+)=(.*)$`)

// --port 8000
var longFormRe = regexp.MustCompile(`^--.+`)
var longFormKeyRe = regexp.MustCompile(`^--(.+)`)

//longFormSpaceValsRe := regexp.MustCompile(`^--([^=])=([\s\S]*)$`)

// --no-debug
var negateRe = regexp.MustCompile(`^--no-.+`)
var negateValsRe = regexp.MustCompile(`^--no-(.+)`)

// -abc
var shortFormRe = regexp.MustCompile(`^-[^-]+`)

var lettersRe = regexp.MustCompile(`^[A-Za-z]`)

var notWordRe = regexp.MustCompile(`\W`)

var dashesRe = regexp.MustCompile(`^(-|--)`)

var trueFalseRe = regexp.MustCompile(`^(true|false)`)

type Config struct {
	Alias    map[string]string
	Defaults map[string]interface{}
}

// Parse parses os.Args excluding os.Args[0].
func Parse(conf *Config) Argv {
	return ParseArgv(os.Args[1:], conf)
}

// ParseArgv parses an argv for options.
func ParseArgv(argv []string, conf *Config) Argv {
	rest := []string{}
	var alias map[string]string
	result := map[string]interface{}{
		"_":  rest,
		"--": []string{},
	}

	if conf != nil {
		if conf.Alias != nil {
			alias = conf.Alias
		}

		for k, v := range conf.Defaults {
			result[k] = v
		}
	}

	var setArg func(key string, val interface{})
	setArg = func(key string, val interface{}) {
		keys := strings.Split(key, ".")
		head := keys[:len(keys)-1]
		last := keys[len(keys)-1]
		res := result

		for _, k := range head {
			existing := res[k]
			if ex, ok := existing.(map[string]interface{}); ok {
				res = ex
			} else {
				child := make(map[string]interface{})
				res[k] = child
				res = child
			}
		}

		res[last] = val

		if a, ok := alias[key]; ok {
			setArg(a, val)
		}
	}

	l := len(argv)
	argsAt := func(i int) string {
		if i > -1 && i < l {
			return argv[i]
		}
		return ""
	}

	i := 0
	for i < len(argv) {
		arg := argv[i]

		if arg == "--" {
			result["--"] = argv[i+1:]
			break
		}

		argAt := func(i int) string {
			if i >= 0 && i < len(arg) {
				return arg[i : i+1]
			}
			return ""
		}
		if longFormEqualRe.MatchString(arg) {
			// --long-form=value

			m := longFormEqualValsRe.FindStringSubmatch(arg)
			//fmt.Printf("--long-form= %s\n", arg)
			setArg(m[1], m[2])

		} else if negateRe.MatchString(arg) {
			//fmt.Printf("--no-flag %s\n", arg)

			m := negateValsRe.FindStringSubmatch(arg)
			setArg(m[1], false)

		} else if longFormRe.MatchString(arg) {
			// --long-form
			//fmt.Printf("--long-form %s\n", arg)

			key := longFormKeyRe.FindStringSubmatch(arg)[1]
			next := argsAt(i + 1)

			if next == "" {
				// --arg
				setArg(key, true)
			} else if next[0:1] == "-" {
				// --arg -o | --arg --other
				setArg(key, true)
			} else {
				setArg(key, next)
				i++
			}
		} else if shortFormRe.MatchString(arg) {
			// -abc a, b are boolean c is undetermined
			//fmt.Printf("-short-form %s\n", arg)

			letters := arg[1:]

			L := len(letters)
			lettersAt := func(i int) string {
				if i < L {
					return letters[i : i+1]
				}
				return ""
			}

			broken := false
			k := 0
			for k < len(letters) {
				next := arg[k+2:]
				if next == "-" {
					setArg(lettersAt(k), next)
					k++
					continue
				}
				if lettersRe.MatchString(lettersAt(k)) && numberRe.MatchString(next) {
					setArg(lettersAt(k), next)
					broken = true
					break
				}
				if k+1 < len(letters) && notWordRe.MatchString(lettersAt(k+1)) {
					setArg(lettersAt(k), next)
					broken = true
					break
				}

				setArg(lettersAt(k), true)
				k++
			}

			key := argAt(len(arg) - 1)
			if !broken && key != "-" {

				if i+1 < len(argv) {
					nextArg := argv[i+1]
					if !dashesRe.MatchString(nextArg) {
						setArg(key, nextArg)
						i++
					}
				} else {
					setArg(key, true)
				}
			}
		} else {
			rest = append(rest, arg)
			result["_"] = rest
		}

		i++
	}

	return result
}
