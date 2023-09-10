package wheels

// TODO auth
import (
	"fmt"
	"reflect"
	"sync"
)

type ServiceZero struct {
	typ  reflect.Type
	name string

	mu       sync.Mutex
	instance any
	value    reflect.Value
	built    bool
}

func newServiceZero(name string, val any) (Service, error) {
	rt := reflect.TypeOf(val)
	if rt.Kind() != reflect.Struct && (rt.Kind() != reflect.Pointer || rt.Elem().Kind() != reflect.Struct) {
		return nil, fmt.Errorf("name: %v, err: %w", name, ErrInvalidZeroType)
	}
	if name == "" {
		name = rt.String()
	}
	return &ServiceZero{
		name:     name,
		instance: val,
		typ:      rt,
		value:    reflect.ValueOf(val),
	}, nil
}

func (s *ServiceZero) getName() string {
	return s.name
}

func (s *ServiceZero) getType() reflect.Type {
	return s.typ
}

func (s *ServiceZero) getInstance(i *Injector) (ins any, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.built {
		return s.instance, nil
	}
	err = s.buildInstanceLocked(i)
	if err != nil {
		return
	}
	return s.instance, nil
}

func (s *ServiceZero) buildInstanceLocked(i *Injector) (err error) {
	val := s.value
	if s.typ.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	for j := 0; j < val.Type().NumField(); j++ {
		fe := val.Field(j)
		if !fe.CanSet() {
			continue
		}
		param, err := i.getValueLocked(fe.Type().String())
		if err != nil {
			return err
		}
		fe.Set(param)
	}
	s.built = true
	i.setInstance(s.name, s.instance)
	return
}

func (s *ServiceZero) getValue(i *Injector) (reflect.Value, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.built {
		i.setEarlyService(s)
	}
	return s.value, nil
}
