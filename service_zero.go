/*
   Copyright (c) 2023, rame2015

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package wheels

import (
	"fmt"
	"reflect"
	"sync"
)

type ServiceZero struct {
	typ  reflect.Type
	name string

	mu         sync.Mutex
	built      bool
	paramNames []string
	value      reflect.Value
	instance   any
}

func newServiceZero(name string, val any) (Service, error) {
	rt := reflect.TypeOf(val)
	if rt.Kind() != reflect.Pointer || rt.Elem().Kind() != reflect.Struct {
		return nil, fmt.Errorf("name: %v, err: %w", name, ErrInvalidZeroType)
	}
	if name == "" {
		name = rt.String()
	}

	s := &ServiceZero{
		name: name,
		typ:  rt,
	}
	s.init()
	return s, nil
}

func (s *ServiceZero) reset() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.built {
		return false
	}
	s.built = false
	s.paramNames = nil
	s.value = reflect.Value{}
	s.init()
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

func (s *ServiceZero) init() {
	switch s.typ.Kind() {
	case reflect.Struct:
		s.value = reflect.New(s.typ).Elem()
	case reflect.Pointer:
		s.value = reflect.New(s.typ.Elem())
	}
	s.instance = s.value.Interface()
}

func (s *ServiceZero) buildInstanceLocked(i *Injector, insName string) (err error) {
	val := s.value
	if s.typ.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	for j := 0; j < val.Type().NumField(); j++ {
		fe := val.Field(j)
		// TODO: check can set before
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
