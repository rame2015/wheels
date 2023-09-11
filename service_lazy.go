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

var errType = reflect.TypeOf(new(error)).Elem()

type ServiceLazy struct {
	name string
	typ  reflect.Type
	ctor reflect.Value // func(...) (...vals, error) or func (...) vals

	mu         sync.RWMutex
	instance   any
	value      reflect.Value
	built      bool
	paramNames []string
}

func newServiceLazy(name string, ctor any) (Service, error) {
	rv := reflect.ValueOf(ctor)
	rt := reflect.TypeOf(ctor)
	if rt.Kind() != reflect.Func {
		return nil, fmt.Errorf("name: %v, err: %w", name, ErrInvalidCtorType)
	}
	numOut := rt.NumOut()
	if numOut == 0 || numOut > 2 {
		return nil, fmt.Errorf("name: %v, err: %w", name, ErrInvalidCtorType)
	}
	if numOut == 2 && !rt.Out(1).Implements(errType) {
		return nil, fmt.Errorf("name: %v, err: %w", name, ErrInvalidCtorType)
	}
	typ := rt.Out(0)
	if name == "" {
		name = typ.String()
	}
	return &ServiceLazy{
		name: name,
		ctor: rv,
		typ:  typ,
	}, nil
}

func (s *ServiceLazy) reset() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.built {
		return false
	}
	s.built = false
	s.paramNames = nil
	return true
}

func (s *ServiceLazy) getType() reflect.Type {
	return s.typ
}

func (s *ServiceLazy) getName() string {
	return s.name
}

func (s *ServiceLazy) getInstance(i *Injector, insName string) (ins any, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.built {
		return s.instance, nil
	}
	err = s.buildInstanceLocked(i, insName)
	if err != nil {
		return nil, err
	}
	return s.instance, nil
}

// buildInstanceLocked TODO support ctx?
func (s *ServiceLazy) buildInstanceLocked(i *Injector, insName string) (err error) {
	ctype := s.ctor.Type()
	paramValues := make([]reflect.Value, ctype.NumIn())
	for j := 0; j < ctype.NumIn(); j++ {
		ptype := ctype.In(j)
		pname := ptype.String()
		pvalue, err := i.getValueLocked(pname)
		if err != nil {
			return err
		}
		paramValues[j] = pvalue
		s.paramNames = append(s.paramNames, pname)
		i.appendAssociatedService(pname, s)
	}
	retValues := s.ctor.Call(paramValues)
	if len(retValues) == 2 {
		errValue := retValues[1]
		if errValue.IsNil() {
			err = nil
		} else {
			err = errValue.Interface().(error)
		}
	}
	if err != nil {
		return
	}
	s.instance = retValues[0].Interface()
	s.value = retValues[0]
	s.built = true
	i.setInstance(insName, s.instance)
	return nil
}

func (s *ServiceLazy) getValue(i *Injector, insName string) (val reflect.Value, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.built {
		return s.value, nil
	}
	err = s.buildInstanceLocked(i, insName)
	if err != nil {
		return val, err
	}
	return s.value, nil
}
