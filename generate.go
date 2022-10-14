package rwalk_peripherals

import (
	"encoding/json"
	"errors"
	"io"
)

func Generate(rwalkPlg, schemasPlg string, w io.Writer) error {
	walk, err := LoadRwalkPlugin(rwalkPlg)
	if err != nil {
		return err
	}

	schemas, err := LoadSchemasPlugin(schemasPlg)
	if err != nil {
		return err
	}

	items := schemas()
	if items == nil {
		return errors.New("schemas invokes should return a non-empty interface{} slice")
	}

	// boolean: ignoreStdPkgPaths, ingoreUnexportedField, ingoreJSONIgnoredField
	result := walk(items, true, true, true)
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "\t")
	return encoder.Encode(result)
}
