package core

import "context"

// Source is a generic interface representing a data provider.
// It defines a single method, Load, which retrieves a slice of DATA items.
//
// The Load method receives a context.Context, which allows cancellation,
// timeouts, and propagation of request-scoped values across API boundaries.
//
// The DATA type parameter allows Source to work with arbitrary data types.
//
// Example usage:
//
//	type User struct { Name string }
//	var src Source[User]
//	users, err := src.Load(context.Background())
type Source[DATA any] interface {
	// Load retrieves a slice of data elements of type DATA.
	// It may return an error if loading fails.
	Load(context.Context) ([]DATA, error)
}

// SourceFunc is a function type adapter that allows ordinary functions to
// satisfy the Source interface.
//
// A SourceFunc is simply a function that matches the Load method signature:
//
//	func(context.Context) ([]DATA, error)
//
// You can use SourceFunc to turn any such function into a Source.
type SourceFunc[DATA any] func(context.Context) ([]DATA, error)

// Load implements the Source interface for SourceFunc.
//
// It simply calls the underlying function f with the provided context.
func (f SourceFunc[DATA]) Load(ctx context.Context) ([]DATA, error) {
	return f(ctx)
}

// MakeSourceFunc wraps a plain function of the correct signature and returns
// it as a SourceFunc.
//
// This helper makes it easier to construct a SourceFunc without explicit type
// casting.
//
// Example:
//
//	fetchUsers := func(ctx context.Context) ([]User, error) { ... }
//	src := MakeSourceFunc(fetchUsers)
func MakeSourceFunc[DATA any](f func(context.Context) ([]DATA, error)) SourceFunc[DATA] {
	return f
}

// MakeSource creates a SourceFunc that always returns the provided data slice
// when Load is called.
//
// This is useful for creating simple, static, or mock data sources — for example
// in tests or when you want to wrap in-memory data as a Source.
//
// Example:
//
//	src := MakeSource(User{Name: "Alice"}, User{Name: "Bob"})
//	users, _ := src.Load(context.Background()) // returns the same slice each time
func MakeSource[DATA any](data ...DATA) SourceFunc[DATA] {
	return MakeSourceFunc(func(ctx context.Context) ([]DATA, error) {
		// Return the provided data without error.
		return data, nil
	})
}
