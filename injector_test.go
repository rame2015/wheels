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
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type ServiceTest interface {
	Print() string
}

type ServiceA struct {
	val int
}

func NewServiceA() *ServiceA {
	return &ServiceA{}
}

func (s *ServiceA) Print() string {
	return "A"
}

type ServiceB struct {
	a *ServiceA
	c *ServiceC
}

func NewServiceB(a *ServiceA, c *ServiceC) (*ServiceB, error) {
	return &ServiceB{a: a, c: c}, nil
}

func (s *ServiceB) Print() string {
	return "B"
}

type ServiceC struct {
	D *ServiceD
}

func (s *ServiceC) Print() string {
	return "C"
}

type ServiceD struct {
	C *ServiceC
	A *ServiceA
}

func (s *ServiceD) Print() string {
	return "D"
}

type ServiceE struct {
	B ServiceTest
}

func (s ServiceE) Print() string {
	return "E"
}

type ServiceF struct {
}

var ErrNewServiceF = errors.New("new service f failed")

func newServiceF() (*ServiceF, error) {
	return nil, ErrNewServiceF
}

type ServiceG struct {
	F *ServiceF
}

type ServiceH struct {
	S ServiceTest
}

type ServiceJ struct {
	s ServiceTest
}

func newServiceJ(s ServiceTest) *ServiceJ {
	return &ServiceJ{s: s}
}

func TestInjector_Provide(t *testing.T) {
	i := New()
	type args struct {
		ctor any
		opts []ProvideOption
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name:    "A ctor that has no return value",
			args:    args{ctor: func() {}},
			wantErr: ErrInvalidCtorType,
		},
		{
			name:    "A ctor that returns a single value",
			args:    args{ctor: NewServiceA},
			wantErr: nil,
		},
		{
			name:    "A ctor that returns two values",
			args:    args{ctor: func() (int64, int64) { return 0, 1 }},
			wantErr: ErrInvalidCtorType,
		},
		{
			name:    "A ctor that returns two values, including an error code",
			args:    args{ctor: NewServiceB},
			wantErr: nil,
		},
		{
			name:    "A ctor that returns three values",
			args:    args{ctor: func() (int64, int64, error) { return 0, 1, nil }},
			wantErr: ErrInvalidCtorType,
		},
		{
			name:    "A constructor that returns an existing service instance",
			args:    args{ctor: func() (*ServiceA, error) { return nil, nil }},
			wantErr: ErrServiceAlreadyExists,
		},
		{
			name:    "A constructor that returns a different named instance of an existing service type",
			args:    args{ctor: func() (*ServiceA, error) { return nil, nil }, opts: []ProvideOption{Name("service a")}},
			wantErr: nil,
		},
		{
			name:    "A struct",
			args:    args{ctor: ServiceA{}},
			wantErr: ErrInvalidCtorType,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := i.Provide(tt.args.ctor, tt.args.opts...); !errors.Is(err, tt.wantErr) {
				t.Errorf("Injector.Provide() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInjector_ProvideInstance(t *testing.T) {
	i := New()
	type args struct {
		val  any
		opts []ProvideOption
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name:    "A struct",
			args:    args{val: ServiceA{}},
			wantErr: nil,
		},
		{
			name:    "A pointer",
			args:    args{val: &ServiceA{}},
			wantErr: nil,
		},
		{
			name:    "An existing service instance",
			args:    args{val: &ServiceA{}},
			wantErr: ErrServiceAlreadyExists,
		},
		{
			name:    "A different named instance of an existing service type",
			args:    args{val: &ServiceA{}, opts: []ProvideOption{Name("service a")}},
			wantErr: nil,
		},
		{
			name:    "A func",
			args:    args{val: func() *ServiceA { return nil }},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := i.ProvideInstance(tt.args.val, tt.args.opts...); !errors.Is(err, tt.wantErr) {
				t.Errorf("Injector.ProvideInstance() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInjector_ProvideZero(t *testing.T) {
	i := New()
	type args struct {
		val  any
		opts []ProvideOption
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name:    "A struct",
			args:    args{val: ServiceA{}},
			wantErr: ErrInvalidZeroType,
		},
		{
			name:    "A pointer",
			args:    args{val: &ServiceA{}},
			wantErr: nil,
		},
		{
			name:    "An existing service instance",
			args:    args{val: &ServiceA{}},
			wantErr: ErrServiceAlreadyExists,
		},
		{
			name:    "A different named instance of an existing service type",
			args:    args{val: &ServiceA{}, opts: []ProvideOption{Name("service a")}},
			wantErr: nil,
		},
		{
			name:    "A func",
			args:    args{val: func() *ServiceA { return nil }},
			wantErr: ErrInvalidZeroType,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := i.ProvideZero(tt.args.val, tt.args.opts...); !errors.Is(err, tt.wantErr) {
				t.Errorf("Injector.ProvideZero() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInjector_ProvideAs(t *testing.T) {
	i := New()
	type args struct {
		provide func(val any, opts ...ProvideOption) error
		val     any
		opts    []ProvideOption
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "An option passed an incorrect format for the interface",
			args: args{
				provide: i.Provide,
				val:     NewServiceB,
				opts:    []ProvideOption{As(ServiceTest(nil))},
			},
			wantErr: ErrInvalidAsType,
		},
		{
			name: "An option passed an interface that is not implemented by the service",
			args: args{
				provide: i.ProvideZero,
				val:     &ServiceG{},
				opts:    []ProvideOption{As(new(ServiceTest))},
			},
			wantErr: ErrServiceNotImplementsAs,
		},
		{
			name: "The option did not pass an interface",
			args: args{
				provide: i.ProvideZero,
				val:     &ServiceA{},
				opts:    []ProvideOption{As(new(int))},
			},
			wantErr: ErrInvalidAsType,
		},
		{
			name: "A constructor that returns a service implemented as an interface",
			args: args{
				provide: i.Provide,
				val:     NewServiceA,
				opts:    []ProvideOption{As(new(ServiceTest))},
			},
		},
		{
			name: "A constructor that returns a service implemented as an existing interface",
			args: args{
				provide: i.ProvideZero,
				val:     &ServiceC{},
				opts:    []ProvideOption{As(new(ServiceTest))},
			},
			wantErr: ErrServiceAlreadyExists,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.args.provide(tt.args.val, tt.args.opts...); !errors.Is(err, tt.wantErr) {
				t.Errorf("Injector.ProvideZero() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInjector_Invoke(t *testing.T) {
	i := New()
	_ = i.ProvideInstance(&ServiceA{})
	_ = i.Provide(NewServiceB, As(new(ServiceTest)))
	_ = i.ProvideZero(&ServiceC{})
	_ = i.ProvideZero(&ServiceD{})
	_ = i.ProvideZero(&ServiceE{})
	_ = i.Provide(newServiceF)
	_ = i.ProvideZero(&ServiceG{})
	type args struct {
		name string
		opts []InvokeOption
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name:    "Unknown service",
			args:    args{name: "wheels.ServiceTestA"},
			wantErr: ErrUnknownService,
		},
		{
			name:    "New instance failed",
			args:    args{name: "*wheels.ServiceF"},
			wantErr: ErrNewServiceF,
		},
		{
			name:    "New param instance failed",
			args:    args{name: "*wheels.ServiceG"},
			wantErr: ErrNewServiceF,
		},
		{
			name:    "A",
			args:    args{name: "*wheels.ServiceA"},
			wantErr: nil,
		},
		{
			name:    "B",
			args:    args{name: "*wheels.ServiceB"},
			wantErr: nil,
		},
		{
			name:    "C",
			args:    args{name: "*wheels.ServiceC"},
			wantErr: nil,
		},
		{
			name:    "D",
			args:    args{name: "*wheels.ServiceD"},
			wantErr: nil,
		},
		{
			name:    "E",
			args:    args{name: "*wheels.ServiceE"},
			wantErr: nil,
		},
		{
			name:    "B",
			args:    args{name: "wheels.ServiceTest"},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ins, err := i.Invoke(tt.args.name, tt.args.opts...)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Injector.Invoke() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && ins.(ServiceTest).Print() != tt.name {
				t.Errorf("Injector.Invoke() name = %v, wantName = %v", ins.(ServiceTest).Print(), tt.name)
			}
		})
	}
}

func TestInjector_Override(t *testing.T) {
	i := New()
	_ = i.Provide(NewServiceB, As(new(ServiceTest)))
	_ = i.ProvideInstance(&ServiceA{})
	_ = i.ProvideZero(&ServiceC{})
	_ = i.ProvideZero(&ServiceD{})
	_ = i.ProvideZero(&ServiceE{})
	b, _ := i.Invoke("wheels.ServiceTest")
	assert.Equal(t, "B", b.(ServiceTest).Print())
	type args struct {
		override func() error
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "A",
			args: args{
				override: func() error {
					return i.Override(NewServiceA, As(new(ServiceTest)))
				},
			},
		},
		{
			name: "C",
			args: args{
				override: func() error {
					return i.OverrideZero(&ServiceC{}, As(new(ServiceTest)))
				},
			},
		},
		{
			name: "D",
			args: args{
				override: func() error {
					return i.OverrideInstance(&ServiceD{}, As(new(ServiceTest)))
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.override()
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Injector.Override() error = %v, wantErr %v", err, tt.wantErr)
			}
			ins, err := i.Invoke("wheels.ServiceTest")
			if err != nil || ins.(ServiceTest).Print() != tt.name {
				t.Errorf("Injector.Invoke() name = %v, wantName = %v, err = %v", ins.(ServiceTest).Print(), tt.name, err)
			}
		})
	}
}

func TestInjector_OverrideAssociated(t *testing.T) {
	i := New()
	_ = i.Provide(NewServiceB, As(new(ServiceTest)))
	_ = i.ProvideInstance(&ServiceA{})
	_ = i.ProvideZero(&ServiceC{})
	_ = i.ProvideZero(&ServiceD{})
	_ = i.ProvideZero(&ServiceE{})
	_ = i.ProvideZero(&ServiceH{})
	_ = i.Provide(newServiceJ)
	h, _ := i.Invoke("*wheels.ServiceH")
	j, _ := i.Invoke("*wheels.ServiceJ")
	assert.Equal(t, "B", h.(*ServiceH).S.Print())
	assert.Equal(t, "B", j.(*ServiceJ).s.Print())
	type args struct {
		override func() error
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "A",
			args: args{
				override: func() error {
					return i.Override(NewServiceA, As(new(ServiceTest)))
				},
			},
		},
		{
			name: "C",
			args: args{
				override: func() error {
					return i.OverrideZero(&ServiceC{}, As(new(ServiceTest)))
				},
			},
		},
		{
			name: "D",
			args: args{
				override: func() error {
					return i.OverrideInstance(&ServiceD{}, As(new(ServiceTest)))
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.override()
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Injector.Override() error = %v, wantErr %v", err, tt.wantErr)
			}
			ins, err := i.Invoke("*wheels.ServiceH")
			if err != nil || ins.(*ServiceH).S.Print() != tt.name {
				t.Errorf("Injector.Invoke() name = %v, wantName = %v, err = %v", ins.(*ServiceH).S.Print(), tt.name, err)
			}
			ins, err = i.Invoke("*wheels.ServiceJ")
			if err != nil || ins.(*ServiceJ).s.Print() != tt.name {
				t.Errorf("Injector.Invoke() name = %v, wantName = %v, err = %v", ins.(*ServiceJ).s.Print(), tt.name, err)
			}
		})
	}
}
