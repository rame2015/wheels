package wheels

import (
	"reflect"
)

type ServiceInstance struct {
	name     string
	typ      reflect.Type
	instance any
	value    reflect.Value
	built    bool
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

func (s *ServiceInstance) reset() bool {
	return false
}

func (s *ServiceInstance) getName() string {
	return s.name
}

func (s *ServiceInstance) getType() reflect.Type {
	return s.typ
}

func (s *ServiceInstance) getInstance(i *Injector, insName string) (any, error) {
	if !s.built {
		i.setInstance(insName, s)
		s.built = true
	}
	return s.instance, nil
}

func (s *ServiceInstance) getValue(i *Injector, insName string) (val reflect.Value, err error) {
	return s.value, nil
}
