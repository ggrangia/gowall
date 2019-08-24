package wall

import (
	"context"
	"time"
)

// Response is the type that is build with the values returned by Brick
type Response struct {
	Resp interface{}
	Err  error
}

type indexedResponse struct {
	Response
	position int
}

// Brick is the signature that functions must have
type Brick func() (interface{}, error)

// Error is an immmutable type
type Error string

func (e Error) Error() string { return string(e) }

func fireFunctions(ctx context.Context, in chan indexedResponse, funcs []Brick) {

	for i, f := range funcs {
		go runFunctions(ctx, in, f, i)
	}
}

func runFunctions(ctx context.Context, in chan indexedResponse, f Brick, pos int) {
	// Swallows sending on closed channel "in"
	defer func() { recover() }()

	// This channel will be used to receive the result of the function
	c := make(chan indexedResponse)

	// Fire a new goroutine that wraps f()
	go func() {
		// Recover from panic caused by f()
		defer func() {
			if r := recover(); r != nil {
				c <- indexedResponse{Response: Response{Err: Error(r.(string)), Resp: nil}, position: pos}
			}
		}()
		d, e := f()
		c <- indexedResponse{Response: Response{Resp: d, Err: e}, position: pos}
	}()

	// Blocking until ctx channel is closed or data is returned
	select {
	case <-ctx.Done():
		return
	case r := <-c:
		in <- r
	}
}

func min(a, b int) int {

	if a <= b {
		return a
	}
	return b
}

// AllSettled Receives a set of functions and return a slice containing the returning values from settled functions.
// Differently from All, it does not halt when an error occurs.
func AllSettled(funcs ...Brick) []Response {

	l := len(funcs)

	// Used to communicate with the functions
	in := make(chan indexedResponse, l)
	defer close(in)

	// ctx
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	resps := make([]Response, l)

	// Fire all the functions
	fireFunctions(ctx, in, funcs)

	for i := 0; i < l; i++ {
		r := <-in
		resps[r.position] = r.Response
	}

	return resps
}

// AllSettledTimed works the same way as AllSettled, but it adds a limitation to the goroutines execution time
// Return false if timeout expires.
func AllSettledTimed(millis time.Duration, funcs ...Brick) ([]Response, bool) {
	l := len(funcs)

	// ctx
	ctx, cancel := context.WithTimeout(context.Background(), millis*time.Millisecond)
	defer cancel()
	// Used to communicate with the functions
	in := make(chan indexedResponse, l)
	defer close(in)

	resps := make([]Response, l)
	state := true

	// Fire all the functions
	fireFunctions(ctx, in, funcs)

	for i := 0; i < l; i++ {
		select {
		case r := <-in:
			resps[r.position] = r.Response

		case <-ctx.Done():
			state = false
			break
		}

	}
	return resps, state

}

// Race receives a set of functions and return the first one that settles, independently if an error occured or not.
func Race(funcs ...Brick) Response {
	//l := len(funcs)

	// ctx
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// Used to communicate with the functions
	in := make(chan indexedResponse)
	defer close(in)

	// Fire all the functions
	fireFunctions(ctx, in, funcs)

	return (<-in).Response
}

// RaceTimed works the same way as Race, but it adds a limitation to the goroutines execution time
// Return false if timeout expires.
func RaceTimed(millis time.Duration, funcs ...Brick) (Response, bool) {

	state := true
	const timeout = Error("Timeout Expired")

	timeoutFunc := func() (interface{}, error) {
		time.Sleep(millis * time.Millisecond)
		return nil, Error("Timeout Expired")
	}
	r := Race(append(funcs, timeoutFunc)...)
	if r.Err == timeout {
		state = false
	}
	return r, state
}

// All receives a set of functions and return a slice containing the returning values
// from settled functions and a boolean indicating if every function succedeed.
// If any error occurs, only the error is returned and the state is set to false.
func All(funcs ...Brick) ([]Response, bool) {
	l := len(funcs)

	// Used to communicate with the functions
	in := make(chan indexedResponse, l)
	defer close(in)

	ctx, cancel := context.WithCancel(context.Background())
	// Close the ctx channel at the end of this function
	defer cancel()

	resps := make([]Response, l)

	// Fire all the functions
	fireFunctions(ctx, in, funcs)

	for i := 0; i < l; i++ {
		r := <-in
		resps[r.position] = r.Response
		if r.Err != nil {
			resps = nil // free the reference
			return []Response{r.Response}, false
		}
	}
	return resps, true
}

// AllTimed works the same way as AllTimed, but it adds a limitation to the goroutines execution time
// Return false if timeout expires.
func AllTimed(millis time.Duration, funcs ...Brick) ([]Response, bool) {

	l := len(funcs)

	// Used to communicate with the functions
	in := make(chan indexedResponse, l)
	defer close(in)

	// ctx
	ctx, cancel := context.WithTimeout(context.Background(), millis*time.Millisecond)
	defer cancel()

	resps := make([]Response, l)

	// Fire all the functions
	fireFunctions(ctx, in, funcs)

	for i := 0; i < l; i++ {
		select {

		case r := <-in:
			resps[r.position] = r.Response
			if r.Err != nil {
				resps = nil // free the reference
				return []Response{r.Response}, false
			}
		case <-ctx.Done():
			return []Response{{Resp: nil, Err: Error("Timeout Expired")}}, false
		}
	}

	return resps, true
}

// Some receives the number of functions (as integer) you want to wait for and a set of functions.
// Some returns a slice with the first functions that settled
// (or all the functions if the passed integer is higher than the number of functions).
func Some(total int, funcs ...Brick) []Response {
	l := len(funcs)

	// Used to communicate with the functions
	in := make(chan indexedResponse, l)
	defer close(in)

	// ctx
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	resps := make([]Response, min(total, l))

	// Fire all the functions
	fireFunctions(ctx, in, funcs)

	for i := 0; i < l; i++ {
		if i >= total {
			break
		}
		r := <-in
		resps[i] = r.Response
	}

	return resps
}

// SomeTimed works the same way as Some, but it adds a limitation to the goroutines execution time
// Return false if timeout expires.
func SomeTimed(millis time.Duration, total int, funcs ...Brick) ([]Response, bool) {

	l := len(funcs)

	// Used to communicate with the functions
	in := make(chan indexedResponse, l)
	defer close(in)

	// ctx
	ctx, cancel := context.WithTimeout(context.Background(), millis*time.Millisecond)
	defer cancel()

	var resps []Response

	// Fire all the functions
	fireFunctions(ctx, in, funcs)

	i := 0
	for {
		if i >= min(l, total) {
			break
		}
		select {
		case r := <-in:
			resps = append(resps, r.Response)
		case <-ctx.Done():
			return append(resps, Response{Err: Error("Timeout Expired"), Resp: nil}), false
		}
		i++
	}

	return resps, true
}
