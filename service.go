package wheels

import "reflect"

type Service interface {
	getName() string
	getType() reflect.Type
	getValue(*Injector) (reflect.Value, error)
	getInstance(*Injector) (any, error)
}
