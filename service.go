package wheels

import "reflect"

type Service interface {
	getName() string
	getType() reflect.Type
	getValue(*Injector, string) (reflect.Value, error)
	getInstance(*Injector, string) (any, error)
	reset() bool
}
