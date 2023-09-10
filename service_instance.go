package wheels

import (
	"reflect"
)

type ServiceInstance struct {
	name     string
	typ      reflect.Type
	instance any
	value    reflect.Value
}

func newServiceInstance(name string, val any) Service {
	rv := reflect.ValueOf(val)
	if name == "" {
		name = rv.Type().String()
	}
	return &ServiceInstance{
		name:     name,
		instance: val,
		typ:      rv.Type(),
		value:    rv,
	}
}

func (s *ServiceInstance) getName() string {
	return s.name
}

func (s *ServiceInstance) getType() reflect.Type {
	return s.typ
}

func (s *ServiceInstance) getInstance(i *Injector) (any, error) {
	return s.instance, nil
}

func (s *ServiceInstance) getValue(i *Injector) (val reflect.Value, err error) {
	return s.value, nil
}
