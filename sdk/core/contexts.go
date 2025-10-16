package core

type SourceContext[DATA any] struct {
	Data  []DATA
	Error error
}

func MakeSourceContext[DATA any](data []DATA, err error) SourceContext[DATA] {
	return SourceContext[DATA]{
		Data:  data,
		Error: err,
	}
}

type FactoryContext[
	POLICY any,
	DATA any,
	IN any,
] struct {
	Source SourceContext[POLICY]
	Data   DATA
	Input  IN
}

func MakeFactoryContext[
	POLICY any,
	DATA any,
	IN any,
](source SourceContext[POLICY], data DATA, in IN) FactoryContext[POLICY, DATA, IN] {
	return FactoryContext[POLICY, DATA, IN]{
		Source: source,
		Data:   data,
		Input:  in,
	}
}
