package policy

// Evaluation is a generic container that represents the result of evaluating
// a policy or decision function. It encapsulates both the computed result and
// any associated error.
//
// This type is particularly useful for returning results from policy engines,
// where both a successful outcome and potential evaluation errors must be
// reported together.
//
// The type parameter RESULT represents the type of the policy’s output.
//
// Example:
//
//	eval := policy.Evaluation[bool]{Result: true, Error: nil}
//	if eval.Error != nil {
//	    log.Println("policy evaluation failed:", eval.Error)
//	} else if eval.Result {
//	    fmt.Println("policy passed")
//	}
//
// Evaluation types can also be returned directly from helper functions or
// pipelines that combine multiple policy checks.
type Evaluation[
	RESULT any,
] struct {
	// Result holds the evaluated outcome value.
	Result RESULT

	// Error captures any error that occurred during evaluation.
	Error error
}

// MakeEvaluation constructs a new Evaluation instance from a given result
// and error value. This helper improves readability and consistency when
// returning evaluation outcomes from functions.
//
// Example:
//
//	func CheckAccess(user string) policy.Evaluation[bool] {
//	    if user == "" {
//	        return policy.MakeEvaluation(false, errors.New("user not specified"))
//	    }
//	    return policy.MakeEvaluation(true, nil)
//	}
//
//	eval := CheckAccess("alice")
//	if eval.Error != nil {
//	    log.Println("error:", eval.Error)
//	} else if eval.Result {
//	    fmt.Println("access granted")
//	}
//
// The MakeEvaluation function is purely syntactic sugar for:
//
//	policy.Evaluation[bool]{Result: value, Error: err}
func MakeEvaluation[
	RESULT any,
](result RESULT, err error) Evaluation[RESULT] {
	return Evaluation[RESULT]{
		Result: result,
		Error:  err,
	}
}
