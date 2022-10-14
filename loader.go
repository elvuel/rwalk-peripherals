package rwalk_peripherals

import (
	"errors"
	"fmt"
	"os"
	"plugin"
)

const (
	RWalkFnSymbol   = "Walk"
	SchemasFNSymbol = "Schemas"
)

type RWalkFn func([]interface{}, bool, bool, bool) map[string]interface{}
type SchemasFn func() []interface{}

func LoadRwalkPlugin(path string) (RWalkFn, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, errors.New("reflection walk plugin not found")
	}

	plg, err := plugin.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open reflection walk plugin: %v", err)
	}

	symbol, err := plg.Lookup(RWalkFnSymbol)

	if err != nil {
		return nil, fmt.Errorf("failed to lookup reflection walk plugin's `Walk` symbol: %v", err)
	}

	var fn RWalkFn
	var ok bool

	if fn, ok = symbol.(func([]interface{}, bool, bool, bool) map[string]interface{}); !ok {
		return nil, errors.New("invalid assertion for reflection walk plugin's `Walk` symbol")
	}

	return fn, nil
}

func LoadSchemasPlugin(path string) (SchemasFn, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, errors.New("schemas plugin not found")
	}

	plg, err := plugin.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open schemas plugin: %v", err)
	}

	symbol, err := plg.Lookup(SchemasFNSymbol)

	if err != nil {
		return nil, fmt.Errorf("failed to lookup schemas plugin's `Schemas` symbol: %v", err)
	}

	var fn SchemasFn
	var ok bool
	if fn, ok = symbol.(func() []interface{}); !ok {
		return nil, errors.New("invalid assertion for schemas plugin's `Schemas` symbol")
	}

	return fn, nil
}
