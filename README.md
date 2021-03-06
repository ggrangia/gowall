# gowall

Simple package to synchronize your goroutines and collect the results.<br>
Inspired by Promise.all() et similia
# Installation
```
go get github.com/ggrangia/gowall
```

# Features

gowall create a synchronization point where your goroutine can rendezvous and return their data to you.<br>
The basic element (function) to be passed is of type Brick.
``` golang
type Brick func() (interface{}, error)
```
Any function will return a slice (or a single element) of type Response.
```golang
type Response struct {
	Resp interface{}
	Err  error
}
```
If a Brick panics, it is recovered and the field Err is populated.<br>
Here the list of available functions:
- <a href="#all"><b>All</b></a>
- <a href="#allsettled"><b>AllSettled</b></a>
- <a href="#some"><b>Some</b></a>
- <a href="#race"><b>Race</b></a>


For each of the above functions, there are also "timed" versions. They have the same behaviour but they halts when the timer expires, cancelling any still running goroutines.

- <a href="#timed-functions"><b>AllTimed</b></a>
- <a href="#timed-functions"><b>AllSettledTimed</b></a>
- <a href="#timed-functions"><b>SomeTimed</b></a>
- <a href="#timed-functions"><b>RaceTimed</b></a>

## All
Receives a set of functions and return a slice containing the returning values from settled functions and a boolean indicating if every function succedeed.
If any error occurs, only the error is returned and the state is set to false.
```golang
func All(...Brick) ([]Response, bool)
```

## AllSettled
Receives a set of functions and return a slice containing the returning values from settled functions. Differently from All, it does not halt when an error occurs.
```golang
func AllSettled(...Brick) []Response
```

## Some
Receives the number of functions (as integer) you want to wait for and a set of functions. Some returns a slice with the first functions that settles (or all the functions if the passed integer is higher than the number of functions).
```golang
func Some(int, ...Brick) []Response
```

## Race
Receives a set of functions and return the first one that settles, independently if an error occured or not.
```golang
func Race(...Brick) Response
```

## Timed functions
All the timed functions receive an additional parameter that is the timeout (in milliseconds).
All the functions returns an additional boolean that indicates if the timeout has expired, killing the running goroutines.
```golang
func AllTimed(time.Duration, ...Brick) ([]Response, bool)
```
```golang
func AllSettledTimed(time.Duration, ...Brick) ([]Response, bool)
```
```golang
func SomeTimed(time.Duration, int, ...Brick) ([]Response, bool)
```
```golang
func RaceTimed(time.Duration,...Brick) (Response, bool)
```
# Examples
Get the content of two http API calls, but stop after 5 seconds.
```golang
import (
	// other imports
	"github.com/ggrangia/gowall"
)

var urls = []string{"http://swapi.co/api/planets/9", "http://swapi.co/api/people/9"}

func httpGet(url string) (interface{}, error) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	return string(bytes), err
}

requests := []wall.Brick{}
for u := range url {
		myurl := u
		httpf = append(httpf, func() (interface{}, error) {
			return httpGet(myurl)
		})
	}
res, isCompleted := wall.AllSettledTimed(5000, requests...)

```
# Contributing
Any help is appreciated! If you see something wrong or you have an improvement, do not hesitate to make a pull request! Thank you in advance.
