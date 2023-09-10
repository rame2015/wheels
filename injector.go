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

type Injector struct {
	instances sync.Map

	mu            sync.RWMutex
	services      map[string]Service
	earlyServices map[string]Service
}

func New() *Injector {
	return &Injector{
		services:      map[string]Service{},
		earlyServices: map[string]Service{},
	}
}

func (i *Injector) Provide(ctor any, opts ...ProvideOption) error {
	options := &providerOptions{}
	for _, po := range opts {
		po(options)
	}
	svc, err := newServiceLazy(options.Name, ctor)
	if err != nil {
		return err
	}
	return i.provide(svc, options)
}

func (i *Injector) ProvideInstance(val any, opts ...ProvideOption) error {
	options := &providerOptions{}
	for _, po := range opts {
		po(options)
	}
	svc := newServiceInstance(options.Name, val)
	err := i.provide(svc, options)
	if err != nil {
		return err
	}
	ins, err := svc.getInstance(i)
	if err != nil {
		return err
	}
	i.setInstance(svc.getName(), ins)
	return nil
}

func (i *Injector) ProvideZero(val any, opts ...ProvideOption) error {
	options := &providerOptions{}
	for _, po := range opts {
		po(options)
	}
	svc, err := newServiceZero(options.Name, val)
	if err != nil {
		return err
	}
	return i.provide(svc, options)
}

func (i *Injector) Invoke(name string, opts ...InvokeOption) (ins any, err error) {
	ins, ok := i.getInstance(name)
	if ok {
		return
	}
	return i.invoke(name)
}

func (i *Injector) provide(svc Service, opts *providerOptions) (err error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	name := svc.getName()
	if i.existLocked(name) {
		return fmt.Errorf("name: %v, err: %w", name, ErrServiceAlreadyExists)
	}

	i.services[name] = svc
	for _, as := range opts.As {
		// check as
		asrv := reflect.ValueOf(as)
		if !asrv.IsValid() {
			return ErrInvalidAsType
		}
		if asrv.Kind() == reflect.Ptr {
			asrv = asrv.Elem()
		}
		err = checkAsType(svc.getType(), asrv.Type())
		if err != nil {
			return err
		}
		asName := asrv.Type().String()
		if i.existLocked(asName) {
			return fmt.Errorf("name: %v, err: %w", asName, ErrServiceAlreadyExists)
		}
		i.services[asrv.Type().String()] = svc
	}
	return
}

func (i *Injector) existLocked(name string) bool {
	_, ok := i.services[name]
	return ok
}

func (i *Injector) invoke(name string, opts ...InvokeOption) (ins any, err error) {
	options := &invokeOptions{}
	for _, io := range opts {
		io(options)
	}
	i.mu.Lock()
	defer i.mu.Unlock()
	svc, ok := i.services[name]
	if !ok {
		return nil, fmt.Errorf("name: %v, err: %w", name, ErrUnknownService)
	}
	ins, err = svc.getInstance(i)
	if err != nil {
		return nil, err
	}
	for len(i.earlyServices) > 0 {
		for k, s := range i.earlyServices {
			_, err = s.getInstance(i)
			if err != nil {
				return nil, err
			}
			delete(i.earlyServices, k)
		}
	}
	return
}

func (i *Injector) getValueLocked(name string) (val reflect.Value, err error) {
	svc, ok := i.services[name]
	if !ok {
		return val, fmt.Errorf("name: %v, err: %w", name, ErrUnknownService)
	}
	return svc.getValue(i)
}

func (i *Injector) getInstance(name string) (any, bool) {
	return i.instances.Load(name)
}

func (i *Injector) setInstance(name string, ins any) {
	i.instances.Store(name, ins)
}

func (i *Injector) setEarlyService(svc Service) {
	i.earlyServices[svc.getName()] = svc
}

func checkAsType(svc, as reflect.Type) error {
	switch as.Kind() {
	case reflect.Interface:
		if svc.Implements(as) {
			return nil
		}
		return fmt.Errorf("service: %v, as: %v, err: %w", svc.String(), as.String(), ErrServiceNotImplementsAs)
	}
	return fmt.Errorf("as: %v, err: %w", as.String(), ErrInvalidAsType)
}
