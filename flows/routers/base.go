package routers

type BaseRouter struct {
	ResultName_ string `json:"result_name,omitempty"`
}

func (r *BaseRouter) ResultName() string { return r.ResultName_ }
