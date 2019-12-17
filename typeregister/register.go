package typeregister

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"
	"sync"
)

// GetTypeName of given struct
func GetTypeName(source interface{}) (reflect.Type, string) {
	rawType := reflect.TypeOf(source)

	// source is a pointer, convert to its value
	if rawType.Kind() == reflect.Ptr {
		rawType = rawType.Elem()
	}

	name := rawType.String()
	// we need to extract only the name without the package
	// name currently follows the format `package.StructName`
	parts := strings.Split(name, ".")
	return rawType, parts[1]
}

// Register stores events so we can deserialize from datastores
type Register interface {
	Set(interface{})
	Get(string) (interface{}, error)
	Decode(string, io.Reader) (interface{}, error)
}

// NewRegister creates a new Register
func NewRegister(items ...interface{}) Register {
	reg := &register{
		types: make(map[string]reflect.Type),
	}

	for _, item := range items {
		reg.Set(item)
	}

	return reg
}

type register struct {
	sync.RWMutex
	types map[string]reflect.Type
}

// Set a new type
func (e *register) Set(source interface{}) {
	e.Lock()
	defer e.Unlock()

	rawType, name := GetTypeName(source)
	e.types[name] = rawType
}

// Get a type based on its name
func (e *register) Get(name string) (interface{}, error) {
	e.RLock()
	defer e.RUnlock()

	for key, value := range e.types {
		if strings.EqualFold(name, key) {
			return reflect.New(value).Interface(), nil
		}
	}

	return nil, fmt.Errorf("Cannot find %s in registry", name)
}

// Decode builds an object and fills it
func (e *register) Decode(name string, body io.Reader) (interface{}, error) {
	val, err := e.Get(name)
	if err != nil {
		return nil, err
	}

	if err := json.NewDecoder(body).Decode(&val); err != nil {
		return nil, err
	}
	return val, nil
}
