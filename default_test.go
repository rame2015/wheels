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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInvoke(t *testing.T) {
	_ = ProvideInstance(&ServiceA{})
	_ = Provide(NewServiceB, As(new(ServiceTest)))
	_ = ProvideZero(&ServiceC{})
	_ = ProvideZero(&ServiceD{})
	a, _ := Invoke[*ServiceA]()
	assert.Equal(t, "A", a.Print())
	b, _ := Invoke[*ServiceB]()
	assert.Equal(t, "B", b.Print())
	c, _ := Invoke[*ServiceC]()
	assert.Equal(t, "C", c.Print())
	d, _ := Invoke[*ServiceD]()
	assert.Equal(t, "D", d.Print())
	assert.Equal(t, "C", d.C.Print())

	err := OverrideInstance(&ServiceA{val: 100})
	assert.NoError(t, err)
	err = Override(NewServiceB, As(new(ServiceTest)))
	assert.NoError(t, err)
	err = OverrideZero(&ServiceC{})
	assert.NoError(t, err)
	err = OverrideZero(&ServiceD{})
	assert.NoError(t, err)

	na, _ := Invoke[*ServiceA]()
	assert.NotSame(t, na, a)
	nb, _ := Invoke[*ServiceB]()
	assert.NotSame(t, nb, b)
	nc, _ := Invoke[*ServiceC]()
	assert.NotSame(t, nc, c)
	nd, _ := Invoke[*ServiceD]()
	assert.NotSame(t, nd, d)
}
