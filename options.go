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
