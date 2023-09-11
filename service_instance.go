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
