package blueprint

import (
	"encoding/json"
	"errors"
	"fmt"
)

var ErrInvalidScript = errors.New("blueprint script is invalid")

// TODO document
type Blueprint interface {
	Values(property string) []string
	Properties() []string

	Child(property string) Blueprint
	Children() []string

	// TODO Properties() and Children() don't contain those reachable through parent
}

type block struct {
	parent   Blueprint
	values   map[string][]string
	children map[string]Blueprint
}

// TODO document
func Parse(data []byte) (Blueprint, error) {
	var raw map[string]interface{}
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return nil, fmt.Errorf("cannot parse JSON: %w", err)
	}

	return parseRaw(raw, nil)
}

func parseRaw(raw map[string]interface{}, parent Blueprint) (Blueprint, error) {
	bp := &block{
		parent:   parent,
		values:   map[string][]string{},
		children: map[string]Blueprint{},
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

func (b *block) parseValues(raw interface{}, k string, valueCounter *int) ([]string, error) {
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

func (b *block) Values(property string) []string {
	if values, ok := b.values[property]; ok {
		return values
	} else if b.parent != nil {
		return b.parent.Values(property)
	} else {
		return nil
	}
}

func (b *block) Properties() []string {
	properties := make([]string, len(b.values))
	i := 0
	for property := range b.values {
		properties[i] = property
		i++
	}
	return properties
}

func (b *block) Child(property string) Blueprint {
	if child, ok := b.children[property]; ok {
		return child
	} else if b.parent != nil {
		return b.parent.Child(property)
	} else {
		return nil
	}
}

func (b *block) Children() []string {
	children := make([]string, len(b.children))
	i := 0
	for property := range b.children {
		children[i] = property
		i++
	}
	return children
}
