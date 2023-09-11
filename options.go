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

type providerOptions struct {
	Name       string
	As         []any
	IsOverride bool
}

type ProvideOption func(*providerOptions)

func Name(name string) ProvideOption {
	return func(po *providerOptions) {
		po.Name = name
	}
}

func As(ifaceOrAOP ...any) ProvideOption {
	return func(po *providerOptions) {
		po.As = append(po.As, ifaceOrAOP...)
	}
}

type invokeOptions struct {
}

type InvokeOption func(*invokeOptions)
