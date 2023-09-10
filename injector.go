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

	"golang.org/x/exp/slices"
)

type Injector struct {
	instances sync.Map

	mu                 sync.RWMutex
	services           map[string]Service
	serviceInstances   map[Service][]string
	earlyServices      map[string]Service
	associatedServices map[string][]Service
}

func New() *Injector {
	return &Injector{
		services:           map[string]Service{},
		serviceInstances:   map[Service][]string{},
		earlyServices:      map[string]Service{},
		associatedServices: map[string][]Service{},
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
	return i.provide(svc, options)
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

func (i *Injector) Override(ctor any, opts ...ProvideOption) error {
	options := &providerOptions{}
	for _, po := range opts {
		po(options)
	}
	svc, err := newServiceLazy(options.Name, ctor)
	if err != nil {
		return err
	}
	return i.override(svc, options)
}

func (i *Injector) OverrideInstance(val any, opts ...ProvideOption) error {
	options := &providerOptions{}
	for _, po := range opts {
		po(options)
	}
	svc := newServiceInstance(options.Name, val)
	return i.override(svc, options)
}

func (i *Injector) OverrideZero(val any, opts ...ProvideOption) error {
	options := &providerOptions{}
	for _, po := range opts {
		po(options)
	}
	svc, err := newServiceZero(options.Name, val)
	if err != nil {
		return err
	}
	return i.override(svc, options)
}
func (i *Injector) provide(svc Service, opts *providerOptions) (err error) {
	i.mu.Lock()
	defer i.mu.Unlock()
	return i.provideLocked(svc, opts)
}

func (i *Injector) override(svc Service, opts *providerOptions) (err error) {
	i.mu.Lock()
	defer i.mu.Unlock()
	opts.IsOverride = true
	err = i.provideLocked(svc, opts)
	if err != nil {
		return err
	}
	return nil
}

func (i *Injector) provideLocked(svc Service, opts *providerOptions) (err error) {
	name := svc.getName()
	oldSvc, ok := i.services[name]
	if !opts.IsOverride && ok {
		return fmt.Errorf("name: %v, err: %w", name, ErrServiceAlreadyExists)
	}
	if opts.IsOverride && ok {
		i.instances.Delete(name)
		i.resetAssociatedService(name)
		i.serviceInstances[oldSvc] = slices.DeleteFunc(i.serviceInstances[oldSvc], func(s string) bool { return s == name })
	}
	insNames := []string{name}
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
		oldAs, ok := i.services[asName]
		if !opts.IsOverride && ok {
			return fmt.Errorf("name: %v, err: %w", asName, ErrServiceAlreadyExists)
		}
		if opts.IsOverride {
			i.instances.Delete(asName)
			i.resetAssociatedService(asName)
			i.serviceInstances[oldAs] = slices.DeleteFunc(i.serviceInstances[oldAs], func(s string) bool { return s == asName })
		}
		insNames = append(insNames, asName)
	}
	i.serviceInstances[svc] = insNames
	for _, v := range insNames {
		i.services[v] = svc
	}
	return
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
	ins, err = svc.getInstance(i, name)
	if err != nil {
		return nil, err
	}
	for len(i.earlyServices) > 0 {
		for k, s := range i.earlyServices {
			_, err = s.getInstance(i, s.getName())
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
	return svc.getValue(i, name)
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

func (i *Injector) appendAssociatedService(paramName string, svc Service) {
	i.associatedServices[paramName] = append(i.associatedServices[paramName], svc)
}

func (i *Injector) resetAssociatedService(name string) {
	svcs := i.associatedServices[name]
	delete(i.associatedServices, name)
	for _, s := range svcs {
		isReset := s.reset()
		if !isReset {
			continue
		}
		for _, insName := range i.serviceInstances[s] {
			i.instances.Delete(insName)
			i.resetAssociatedService(insName)
		}
	}
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
