package typeregister

import (
	"fmt"
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
	Set(source interface{})
	Get(name string) (interface{}, error)
}

// NewRegister creates a new Register
func NewRegister() Register {
	return &register{
		types: make(map[string]reflect.Type),
	}
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

	rawType, ok := e.types[name]
	if !ok {
		return nil, fmt.Errorf("Cannot find %s in registry", name)
	}

	return reflect.New(rawType).Interface(), nil
}
