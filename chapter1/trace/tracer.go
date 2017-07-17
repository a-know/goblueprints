package trace

type Tracer interface {
	Trace(...interface{})
}
