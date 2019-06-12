# gowall

Simple package to synnchronize your goroutines and collect the results.<br>
Loosely inspired by Promise.all() et similia
# Installation
```
go get github.com/GiacomoGrangia/gowall
```

# Features

gowall create a synchronization point where your goroutine can rendezvous and return their data to you.<br>
The basic element (function) to be passed is of type Brick.
``` golang
type Brick func() (interface{}, error)
```
Any functions will return a slice (or a single element) of type Response.
```golang
type Response struct {
	Resp interface{}
	Err  error
}
```
If a Brick panics, it is recovered and the field Err is populated.

## All
Receives a set of functions and return a slice 

- All(...Brick)
- AllSettled(...Brick)
- Some(int, ...Brick)
- Race(...Brick)

For each of the above functions, there are also "timed" versions. They have the same behaviour but they halts when the timer expires, cancelling any still running goroutines.

- AllTimed(time.Duration, ...Brick)
- AllSettledTimed(time.Duration, ...Brick)
- SomeTimed(time.Duration, int, ...Brick)
- RaceTimed(time.Duration, ...Brick)
