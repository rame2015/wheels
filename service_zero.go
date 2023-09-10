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

	mu         sync.Mutex
	instance   any
	value      reflect.Value
	built      bool
	paramNames []string
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

func (s *ServiceZero) reset() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.built {
		return false
	}
	s.built = false
	s.paramNames = nil
	return true
}

func (s *ServiceZero) getName() string {
	return s.name
}

func (s *ServiceZero) getType() reflect.Type {
	return s.typ
}

func (s *ServiceZero) getInstance(i *Injector, insName string) (ins any, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.built {
		return s.instance, nil
	}
	err = s.buildInstanceLocked(i, insName)
	if err != nil {
		return
	}
	return s.instance, nil
}

func (s *ServiceZero) buildInstanceLocked(i *Injector, insName string) (err error) {
	val := s.value
	if s.typ.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	for j := 0; j < val.Type().NumField(); j++ {
		fe := val.Field(j)
		if !fe.CanSet() {
			continue
		}
		pname := fe.Type().String()
		param, err := i.getValueLocked(pname)
		if err != nil {
			return err
		}
		fe.Set(param)
		s.paramNames = append(s.paramNames, pname)
		i.appendAssociatedService(pname, s)
	}
	s.built = true
	i.setInstance(insName, s.instance)
	return
}

func (s *ServiceZero) getValue(i *Injector, insName string) (reflect.Value, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.built {
		i.setEarlyService(s)
	}
	return s.value, nil
}
