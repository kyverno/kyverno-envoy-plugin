package engine

type Evaluation[
	RESULT any,
] struct {
	Result RESULT
	Error  error
}

func MakeEvaluation[
	RESULT any,
](result RESULT, err error) Evaluation[RESULT] {
	return Evaluation[RESULT]{
		Result: result,
		Error:  err,
	}
}
