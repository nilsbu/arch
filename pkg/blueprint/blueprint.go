package blueprint

import (
	"encoding/json"
	"errors"
	"fmt"
)

// ErrInvalidScript is returned when a script isn't valid.
var ErrInvalidScript = errors.New("blueprint script is invalid")

// A Blueprint describes an aspect of a level.
// Properties of a Blueprint are organized in two ways: children and values. Children are themselves blueprints.
// Values are of type []string. They both are accessibly via their property name. All names of children start with '*'
// and none of the names of values do.
//
// Through the children, a tree structure is defined. Each node may have an arbitrary number of children and the
// tree may have any height.
//
// Properties have visibility beyond the blueprint they are defined in. If a property is queried in one blueprint but
// no property of that name has been defined there, the parents are recursively called until a property of that name is
// found.
//
// Properties can get parsed from files.
type Blueprint struct {
	parent   *Blueprint
	values   map[string][]string
	children map[string]*Blueprint
}

// TODO Properties() and Children() don't contain those reachable through parent

// Parse creates a Blueprint from the content of a file.
//
// An *json.InvalidUnmarshalError is returned when the file is no valid JSON.
// An ErrInvalidScript is returned when the content isn't valid.
func Parse(data []byte) (*Blueprint, error) {
	var raw map[string]interface{}
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return nil, fmt.Errorf("cannot parse JSON: %w", err)
	}

	return parseRaw(raw, nil)
}

func parseRaw(raw map[string]interface{}, parent *Blueprint) (*Blueprint, error) {
	bp := &Blueprint{
		parent:   parent,
		values:   map[string][]string{},
		children: map[string]*Blueprint{},
	}

	for k, v := range raw {
		valueCounter := 0
		if strs, err := bp.parseValues(v, k, &valueCounter); err != nil {
			return nil, err
		} else {
			bp.values[k] = strs
		}
	}

	return bp, nil
}

func (b *Blueprint) parseValues(raw interface{}, k string, valueCounter *int) ([]string, error) {
	switch value := raw.(type) {
	case string:
		return []string{value}, nil
	case map[string]interface{}:
		str := fmt.Sprintf("*%v%v", k, *valueCounter)
		*valueCounter++
		if child, err := parseRaw(value, b); err != nil {
			return nil, err
		} else {
			b.children[str] = child
		}
		return []string{str}, nil
	case []interface{}:
		values := make([]string, 0)
		for _, elem := range value {
			if strs, err := b.parseValues(elem, k, valueCounter); err != nil {
				return nil, err
			} else {
				values = append(values, strs...)
			}
		}
		return values, nil
	default:
		return nil, fmt.Errorf("%w: '%v' is not valid type for a property", ErrInvalidScript, raw)
	}
}

// Values returnes the values associated with a property.
// If the queried Blueprint doesn't contain that property, the parents are called recursively.
func (b *Blueprint) Values(property string) []string {
	if values, ok := b.values[property]; ok {
		return values
	} else if b.parent != nil {
		return b.parent.Values(property)
	} else {
		return nil
	}
}

// Properties returns all the properties defined in the Blueprint, that have values as data.
// Since Values() additionally does recursive calls, the list returned here doesn't match the properties that are
// accessible through Values().
func (b *Blueprint) Properties() []string {
	properties := make([]string, len(b.values))
	i := 0
	for property := range b.values {
		properties[i] = property
		i++
	}
	return properties
}

// Child returnes the child associated with a property.
// If the queried Blueprint doesn't contain that property, the parents are called recursively.
func (b *Blueprint) Child(property string) *Blueprint {
	if child, ok := b.children[property]; ok {
		return child
	} else if b.parent != nil {
		return b.parent.Child(property)
	} else {
		return nil
	}
}

// Children returns all the names of the children defined in the Blueprint.
// Since Child() additionally does recursive calls, the list returned here doesn't match the properties that are
// accessible through Child().
func (b *Blueprint) Children() []string {
	children := make([]string, len(b.children))
	i := 0
	for property := range b.children {
		children[i] = property
		i++
	}
	return children
}
