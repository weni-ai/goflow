package utils

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var snakedChars = regexp.MustCompile(`[^\p{L}\d_]+`)

// Snakify turns the passed in string into a context reference. We replace all whitespace
// characters with _ and replace any duplicate underscores
func Snakify(text string) string {
	return strings.Trim(strings.ToLower(snakedChars.ReplaceAllString(text, "_")), "_")
}

// VariableResolver is the interface for objects in the context which can be keyed into, e.g. foo.bar
type VariableResolver interface {
	Resolve(key string) interface{}
}

// VariableIndexer is the interface for objects in the context which can be indexed into, e.g. foo.0
type VariableIndexer interface {
	Index(index int) interface{}
	Length() int
}

// VariableAtomizer is the interface for objects in the context which can reduce themselves to an XAtom primitive
type VariableAtomizer interface {
	Atomize() interface{}
}

// ResolveVariable will resolve the passed in string variable given in dot notation and return
// the value as defined by the VariableResolver passed in.
//
// Example syntaxes:
//      foo.bar.0  - 0th element of bar slice within foo, could also be "0" key in bar map within foo
//      foo.bar[0] - same as above
func ResolveVariable(env Environment, variable interface{}, key string) interface{} {
	var err error

	err, isErr := variable.(error)
	if isErr {
		return err
	}

	// self referencing
	if key == "" {
		return variable
	}

	// strip leading '.'
	if key[0] == '.' {
		key = key[1:]
	}

	rest := key
	for rest != "" {
		key, rest = popNextVariable(rest)

		if IsNil(variable) {
			return fmt.Errorf("can't resolve key '%s' of nil", key)
		}

		// is our key numeric?
		index, err := strconv.Atoi(key)
		if err == nil {
			indexable, isIndexable := variable.(VariableIndexer)
			if isIndexable {
				if index >= indexable.Length() || index < -indexable.Length() {
					return fmt.Errorf("index %d out of range for %d items", index, indexable.Length())
				}
				if index < 0 {
					index += indexable.Length()
				}
				variable = indexable.Index(index)
				continue
			}
		}

		resolver, isResolver := variable.(VariableResolver)

		// look it up in our resolver
		if isResolver {
			variable = resolver.Resolve(key)

			err, isErr := variable.(error)
			if isErr {
				return err
			}

		} else if IsMap(variable) {
			variable, err = LookupKey(variable, key)
			if err != nil {
				return err
			}

		} else {
			return fmt.Errorf("can't resolve key '%s' of type %s", key, reflect.TypeOf(variable))
		}
	}

	return variable
}

// popNextVariable pops the next variable off our string:
//     foo.bar.baz -> "foo", "bar.baz"
//     foo[0].bar -> "foo", "[0].baz"
//     foo.0.bar -> "foo", "0.baz"
//     [0].bar -> "0", "bar"
//     foo["my key"] -> "foo", "my key"
func popNextVariable(input string) (string, string) {
	var keyStart = 0
	var keyEnd = -1
	var restStart = -1

	for i, c := range input {
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
		return input, ""
	}

	key := strings.Trim(input[keyStart:keyEnd], "\"")
	rest := input[restStart:]

	return key, rest
}

type mapResolver struct {
	values map[string]interface{}
}

// NewMapResolver returns a simple resolver that resolves variables according to the values
// passed in
func NewMapResolver(values map[string]interface{}) VariableResolver {
	return &mapResolver{
		values: values,
	}
}

// Resolve resolves the given key when this map is referenced in an expression
func (r *mapResolver) Resolve(key string) interface{} {
	val, found := r.values[key]
	if !found {
		return fmt.Errorf("no key '%s' in map", key)
	}
	return val
}

// Atomize is called when this object needs to be reduced to a primitive
func (r *mapResolver) Atomize() interface{} { return fmt.Sprintf("%s", r.values) }

var _ VariableAtomizer = (*mapResolver)(nil)
var _ VariableResolver = (*mapResolver)(nil)
