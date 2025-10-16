package defaults

type SourceResult[DATA any] struct {
	Data  []DATA
	Error error
}

func MakeSourceResult[DATA any](data []DATA, err error) SourceResult[DATA] {
	return SourceResult[DATA]{
		Data:  data,
		Error: err,
	}
}

type PolicyResult[
	POLICY any,
	IN any,
	OUT any,
] struct {
	Policy POLICY
	Input  IN
	Out    OUT
}

func MakePolicyResult[
	POLICY any,
	IN any,
	OUT any,
](policy POLICY, in IN, out OUT) PolicyResult[POLICY, IN, OUT] {
	return PolicyResult[POLICY, IN, OUT]{
		Policy: policy,
		Input:  in,
		Out:    out,
	}
}

type Result[
	POLICY any,
	DATA any,
	IN any,
	OUT any,
] struct {
	Source   SourceResult[POLICY]
	Data     DATA
	Input    IN
	Policies []PolicyResult[POLICY, IN, OUT]
}

func MakeResult[
	POLICY any,
	DATA any,
	IN any,
	OUT any,
](source SourceResult[POLICY], in IN, data DATA, policies []PolicyResult[POLICY, IN, OUT]) Result[POLICY, DATA, IN, OUT] {
	return Result[POLICY, DATA, IN, OUT]{
		Source:   source,
		Input:    in,
		Data:     data,
		Policies: policies,
	}
}
