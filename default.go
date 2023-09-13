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

import "fmt"

var defaultInjector *Injector = New()

func Default() *Injector {
	return defaultInjector
}

func Provide(ctor any, opts ...ProvideOption) error {
	return Default().Provide(ctor, opts...)
}

func ProvideInstance(val any, opts ...ProvideOption) error {
	return Default().ProvideInstance(val, opts...)
}

func ProvideZero(val any, opts ...ProvideOption) error {
	return Default().ProvideZero(val, opts...)
}

func Override(val any, opts ...ProvideOption) error {
	return Default().Override(val, opts...)
}

func OverrideInstance(val any, opts ...ProvideOption) error {
	return Default().OverrideInstance(val, opts...)
}

func OverrideZero(val any, opts ...ProvideOption) error {
	return Default().OverrideZero(val, opts...)
}

func Invoke[T any](opts ...InvokeOption) (ins T, err error) {
	name := fmt.Sprintf("%T", ins)
	val, err := Default().invoke(name, opts...)
	if err != nil {
		return
	}
	ins, ok := val.(T)
	if !ok {
		return ins, fmt.Errorf("name: %v, err: %v", name, ErrInvalidInvokeType)
	}
	return
}
