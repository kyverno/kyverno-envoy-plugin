package matching

import (
	"context"
	"strings"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/engine/assert"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/engine/template"
	"github.com/kyverno/kyverno-json/pkg/apis/policy/v1alpha1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

//nolint:errname
type Result struct {
	field.ErrorList
	Message string
}

func (r Result) Error() string {
	var lines []string
	if r.Message != "" {
		lines = append(lines, "-> "+r.Message)
	}
	for _, err := range r.ErrorList {
		lines = append(lines, " -> "+err.Error())
	}
	return strings.Join(lines, "\n")
}

//nolint:errname
type Results []Result

func (r Results) Error() string {
	var lines []string
	for _, err := range r {
		lines = append(lines, err.Error())
	}
	return strings.Join(lines, "\n")
}

// func MatchAssert(ctx context.Context, path *field.Path, match *v1alpha1.Assert, actual any, bindings binding.Bindings, opts ...template.Option) ([]error, error) {
func MatchAssert(ctx context.Context, path *field.Path, match *v1alpha1.Assert, actual any, bindings binding.Bindings, opts ...template.Option) ([]Result, error) {
	if match == nil || (len(match.Any) == 0 && len(match.All) == 0) {
		return nil, field.Invalid(path, match, "an empty assert is not valid")
	} else {
		if len(match.Any) != 0 {
			var fails []Result
			path := path.Child("any")
			for i, assertion := range match.Any {
				checkFails, err := assert.Assert(ctx, path.Index(i).Child("check"), assert.Parse(ctx, assertion.Check.Value), actual, bindings, opts...)
				if err != nil {
					return fails, err
				}
				if len(checkFails) == 0 {
					fails = nil
					break
				}
				fail := Result{
					ErrorList: checkFails,
				}
				if assertion.Message != "" {
					fail.Message = template.String(ctx, assertion.Message, actual, bindings, opts...)
				}
				fails = append(fails, fail)
			}
			if fails != nil {
				return fails, nil
			}
		}
		if len(match.All) != 0 {
			var fails []Result
			path := path.Child("all")
			for i, assertion := range match.All {
				checkFails, err := assert.Assert(ctx, path.Index(i).Child("check"), assert.Parse(ctx, assertion.Check.Value), actual, bindings, opts...)
				if err != nil {
					return fails, err
				}
				if len(checkFails) > 0 {
					fail := Result{
						ErrorList: checkFails,
					}
					if assertion.Message != "" {
						fail.Message = template.String(ctx, assertion.Message, actual, bindings, opts...)
					}
					fails = append(fails, fail)
				}
			}
			return fails, nil
		}
		return nil, nil
	}
}

func Match(ctx context.Context, path *field.Path, match *v1alpha1.Match, actual any, bindings binding.Bindings, opts ...template.Option) (field.ErrorList, error) {
	if match == nil || (len(match.Any) == 0 && len(match.All) == 0) {
		return nil, field.Invalid(path, match, "an empty match is not valid")
	} else {
		var errs field.ErrorList
		if len(match.Any) != 0 {
			_errs, err := MatchAny(ctx, path.Child("any"), match.Any, actual, bindings, opts...)
			if err != nil {
				return errs, err
			}
			errs = append(errs, _errs...)
		}
		if len(match.All) != 0 {
			_errs, err := MatchAll(ctx, path.Child("all"), match.All, actual, bindings, opts...)
			if err != nil {
				return errs, err
			}
			errs = append(errs, _errs...)
		}
		return errs, nil
	}
}

func MatchAny(ctx context.Context, path *field.Path, assertions []v1alpha1.Any, actual any, bindings binding.Bindings, opts ...template.Option) (field.ErrorList, error) {
	var errs field.ErrorList
	for i, assertion := range assertions {
		_errs, err := assert.Assert(ctx, path.Index(i), assert.Parse(ctx, assertion.Value), actual, bindings, opts...)
		if err != nil {
			return errs, err
		}
		if len(_errs) == 0 {
			return nil, nil
		}
		errs = append(errs, _errs...)
	}
	return errs, nil
}

func MatchAll(ctx context.Context, path *field.Path, assertions []v1alpha1.Any, actual any, bindings binding.Bindings, opts ...template.Option) (field.ErrorList, error) {
	var errs field.ErrorList
	for i, assertion := range assertions {
		_errs, err := assert.Assert(ctx, path.Index(i), assert.Parse(ctx, assertion.Value), actual, bindings, opts...)
		if err != nil {
			return errs, err
		}
		errs = append(errs, _errs...)
	}
	return errs, nil
}
