package routers

type BaseRouter struct {
	Name_ string `json:"name,omitempty"`
}

func (r *BaseRouter) Name() string { return r.Name_ }
