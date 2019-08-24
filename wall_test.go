package wall

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

const (
	timeout  = Error("Timeout Expired")
	panicErr = Error("I am in panic")
)

func ft(text string, milli time.Duration) (interface{}, error) {
	time.Sleep(milli * time.Millisecond)
	return text, nil
}

func panicFunction() (interface{}, error) {
	panic("I am in panic")
}

func TestAll(t *testing.T) {
	t.Run("NormalExecution", func(t *testing.T) {
		want := []Response{{"1", nil}, {"2", nil}}
		wState := true
		var httpf = []Brick{}
		httpf = append(httpf, func() (interface{}, error) {
			return ft("1", 500)
		})
		httpf = append(httpf, func() (interface{}, error) {
			return ft("2", 750)
		})
		r, s := All(httpf...)
		if !cmp.Equal(r, want) || wState != s {
			t.Errorf("NormalExecution is not correct. Want\n%v %t\n got\n%v %t", want, wState, r, s)
		}
	})
	t.Run("PanicExecution", func(t *testing.T) {
		want := []Response{{nil, panicErr}}
		wState := false
		var httpf = []Brick{}
		httpf = append(httpf, func() (interface{}, error) {
			return ft("1", 500)
		})
		httpf = append(httpf, panicFunction)
		r, s := All(httpf...)
		if !cmp.Equal(r, want) || wState != s {
			t.Errorf("PanicExecution is not correct. Want\n%v %t\n got\n%v %t", want, wState, r, s)
		}
	})
}

func TestAllTimed(t *testing.T) {
	// TODO: The tests fails
	t.Run("NormalExecution", func(t *testing.T) {
		want := []Response{{"1", nil}, {"2", nil}}
		wState := true
		var httpf = []Brick{}
		httpf = append(httpf, func() (interface{}, error) {
			return ft("1", 500)
		})
		httpf = append(httpf, func() (interface{}, error) {
			return ft("2", 750)
		})
		r, s := AllTimed(1000, httpf...)
		if !cmp.Equal(r, want) || wState != s {
			t.Errorf("NormalExecution is not correct. Want\n%v %t\n got\n%v %t", want, wState, r, s)
		}
	})
	t.Run("PanicExecution", func(t *testing.T) {
		want := []Response{{nil, panicErr}}
		wState := false
		var httpf = []Brick{}
		httpf = append(httpf, func() (interface{}, error) {
			return ft("1", 500)
		})
		httpf = append(httpf, panicFunction)
		r, s := AllTimed(1000, httpf...)
		if !cmp.Equal(r, want) || wState != s {
			t.Errorf("PanicExecution is not correct. Want\n%v %t\n got\n%v %t", want, wState, r, s)
		}
	})
	t.Run("TimeoutExecution", func(t *testing.T) {
		want := []Response{{nil, timeout}}
		wState := false
		var httpf = []Brick{}
		httpf = append(httpf, func() (interface{}, error) {
			return ft("1", 500)
		})
		httpf = append(httpf, func() (interface{}, error) {
			return ft("2", 750)
		})
		r, s := AllTimed(250, httpf...)
		if !cmp.Equal(r, want) || wState != s {
			t.Errorf("TimeoutExecution is not correct. Want\n%v %t\n got\n%v %t", want, wState, r, s)
		}
	})
}

func TestAllSettled(t *testing.T) {
	t.Run("NormalExecution", func(t *testing.T) {
		want := []Response{{"1", nil}, {"2", nil}}
		var httpf = []Brick{}
		httpf = append(httpf, func() (interface{}, error) {
			return ft("1", 500)
		})
		httpf = append(httpf, func() (interface{}, error) {
			return ft("2", 750)
		})
		r := AllSettled(httpf...)
		if !cmp.Equal(r, want) {
			t.Errorf("NormalExecution is not correct. Want\n%v\n got\n%v", want, r)
		}
	})
	t.Run("PanicExecution", func(t *testing.T) {
		want := []Response{{"1", nil}, {nil, panicErr}}
		var httpf = []Brick{}
		httpf = append(httpf, func() (interface{}, error) {
			return ft("1", 500)
		})
		httpf = append(httpf, panicFunction)
		r := AllSettled(httpf...)
		if !cmp.Equal(r, want) {
			t.Errorf("PanicExecution is not correct. Want\n%v\n got\n%v", want, r)
		}
	})
}
func TestAllSettledTimed(t *testing.T) {

	t.Run("NormalExecution", func(t *testing.T) {
		want := []Response{{"1", nil}, {"2", nil}}
		wState := true
		wall := []Brick{}
		wall = append(wall, func() (interface{}, error) {
			return ft("1", 300)
		})
		wall = append(wall, func() (interface{}, error) {
			return ft("2", 5)
		})

		r, s := AllSettledTimed(500, wall...)
		if !cmp.Equal(r, want) || s != wState {
			t.Errorf("NormalExecution is not correct. Want\n%v %t\n got\n%v %t", want, wState, r, s)
		}
	})
	t.Run("PanicExecution", func(t *testing.T) {
		want := []Response{{nil, panicErr}, {"1", nil}}
		wState := true
		wall := []Brick{}
		wall = append(wall, func() (interface{}, error) {
			return panicFunction()
		})
		wall = append(wall, func() (interface{}, error) {
			return ft("1", 300)
		})

		r, s := AllSettledTimed(500, wall...)
		if !cmp.Equal(r, want) || s != wState {
			t.Errorf("PanicExecution is not correct. Want\n%v %t\n got\n%v %t", want, wState, r, s)
		}
	})
	t.Run("TimeoutExecution", func(t *testing.T) {
		want := []Response{{"1", nil}, {nil, nil}, {nil, nil}}
		wState := false
		wall := []Brick{}
		wall = append(wall, func() (interface{}, error) {
			return ft("1", 100)
		})
		wall = append(wall, func() (interface{}, error) {
			return ft("2", 500)
		})
		wall = append(wall, func() (interface{}, error) {
			return ft("3", 1000)
		})

		r, s := AllSettledTimed(250, wall...)
		if !cmp.Equal(r, want) || s != wState {
			t.Errorf("TimeoutExecution is not correct. Want\n%v %t\n got\n%v %t", want, wState, r, s)
		}
	})

}

func TestRace(t *testing.T) {
	t.Run("NormalExecution", func(t *testing.T) {
		want := Response{"1", nil}
		var httpf = []Brick{}
		httpf = append(httpf, func() (interface{}, error) {
			return ft("1", 100)
		})
		httpf = append(httpf, func() (interface{}, error) {
			return ft("2", 500)
		})
		r := Race(httpf...)
		if !cmp.Equal(r, want) {
			t.Errorf("NormalExecution is not correct. Want\n%v\n got\n%v", want, r)
		}
	})
	t.Run("PanicExecution", func(t *testing.T) {
		want := Response{nil, panicErr}
		var httpf = []Brick{}
		httpf = append(httpf, panicFunction)
		httpf = append(httpf, func() (interface{}, error) {
			return ft("2", 500)
		})
		r := Race(httpf...)
		if !cmp.Equal(r, want) {
			t.Errorf("PanicExecution is not correct. Want\n%v\n got\n%v", want, r)
		}
	})
}

func TestRaceTimed(t *testing.T) {
	t.Run("NormalExecution", func(t *testing.T) {
		want := Response{"1", nil}
		wState := true
		var httpf = []Brick{}
		httpf = append(httpf, func() (interface{}, error) {
			return ft("1", 1000)
		})
		httpf = append(httpf, func() (interface{}, error) {
			return ft("2", 1500)
		})
		r, s := RaceTimed(2000, httpf...)
		if !cmp.Equal(r, want) || s != wState {
			t.Errorf("NormalExecution is not correct. Want\n%v %t\n got\n%v %t", want, wState, r, s)
		}
	})
	t.Run("TimeoutExecution", func(t *testing.T) {
		want := Response{nil, timeout}
		wState := false
		var httpf = []Brick{}
		httpf = append(httpf, func() (interface{}, error) {
			return ft("1", 1000)
		})
		httpf = append(httpf, func() (interface{}, error) {
			return ft("2", 2000)
		})
		r, s := RaceTimed(500, httpf...)
		if !cmp.Equal(r, want) || s != wState {
			t.Errorf("TimeoutExecution is not correct. Want\n%v %t\n got\n%v %t", want, wState, r, s)
		}
	})
	t.Run("PanicExecution", func(t *testing.T) {
		want := Response{nil, panicErr}
		wState := true
		var httpf = []Brick{}
		httpf = append(httpf, panicFunction)
		httpf = append(httpf, func() (interface{}, error) {
			return ft("2", 1500)
		})
		r, s := RaceTimed(2000, httpf...)
		if !cmp.Equal(r, want) || s != wState {
			t.Errorf("PanicExecution is not correct. Want\n%v %t\n got\n%v %t", want, wState, r, s)
		}
	})
}

func TestSome(t *testing.T) {

	t.Run("NormalExecution", func(t *testing.T) {
		want := []Response{{"2", nil}, {"3", nil}}
		wall := []Brick{}
		wall = append(wall, func() (interface{}, error) {
			return ft("1", 300)
		})
		wall = append(wall, func() (interface{}, error) {
			return ft("2", 100)
		})
		wall = append(wall, func() (interface{}, error) {
			return ft("3", 200)
		})

		r := Some(2, wall...)
		if !cmp.Equal(r, want) {
			t.Errorf("NormalExecution is not correct. Want\n%v\n got\n%v", want, r)
		}
	})
	t.Run("NormalExecutionMore", func(t *testing.T) {
		want := []Response{{"2", nil}, {"3", nil}, {"1", nil}}
		wall := []Brick{}
		wall = append(wall, func() (interface{}, error) {
			return ft("1", 300)
		})
		wall = append(wall, func() (interface{}, error) {
			return ft("2", 100)
		})
		wall = append(wall, func() (interface{}, error) {
			return ft("3", 200)
		})

		r := Some(6, wall...)
		if !cmp.Equal(r, want) {
			t.Errorf("NormalExecutionMore is not correct. Want\n%v\n got\n%v", want, r)
		}
	})
	t.Run("PanicExecution", func(t *testing.T) {
		want := []Response{{nil, panicErr}}
		wall := []Brick{}
		wall = append(wall, func() (interface{}, error) {
			return panicFunction()
		})
		wall = append(wall, func() (interface{}, error) {
			return ft("1", 300)
		})

		r := Some(1, wall...)
		if !cmp.Equal(r, want) {
			t.Errorf("PanicExecution is not correct. Want\n%v\n got\n%v", want, r)
		}
	})
}

func TestSomeTimed(t *testing.T) {
	t.Run("NormalExecution", func(t *testing.T) {
		want := []Response{{"2", nil}, {"3", nil}}
		wState := true
		wall := []Brick{}
		wall = append(wall, func() (interface{}, error) {
			return ft("1", 300)
		})
		wall = append(wall, func() (interface{}, error) {
			return ft("2", 100)
		})
		wall = append(wall, func() (interface{}, error) {
			return ft("3", 200)
		})

		r, s := SomeTimed(500, 2, wall...)
		if !cmp.Equal(r, want) || s != wState {
			t.Errorf("NormalExecution is not correct. Want\n%v %t\n got\n%v %t", want, wState, r, s)
		}
	})
	t.Run("NormalExecutionMore", func(t *testing.T) {
		want := []Response{{"2", nil}, {"3", nil}}
		wState := true
		wall := []Brick{}
		wall = append(wall, func() (interface{}, error) {
			return ft("1", 300)
		})
		wall = append(wall, func() (interface{}, error) {
			return ft("2", 100)
		})
		wall = append(wall, func() (interface{}, error) {
			return ft("3", 200)
		})

		r, s := SomeTimed(500, 2, wall...)
		if !cmp.Equal(r, want) || s != wState {
			t.Errorf("NormalExecutionMore is not correct. Want\n%v %t\n got\n%v %t", want, wState, r, s)
		}
	})
	t.Run("PanicExecution", func(t *testing.T) {
		want := []Response{{nil, panicErr}}
		wState := true
		wall := []Brick{}
		wall = append(wall, func() (interface{}, error) {
			return ft("1", 300)
		})
		wall = append(wall, panicFunction)

		r, s := SomeTimed(500, 1, wall...)
		if !cmp.Equal(r, want) || s != wState {
			t.Errorf("PanicExecution is not correct. Want\n%v %t\n got\n%v %t", want, wState, r, s)
		}
	})
	t.Run("TimeoutExecution", func(t *testing.T) {
		want := []Response{{nil, panicErr}, {"2", nil}, {nil, timeout}}
		wState := false
		wall := []Brick{}
		wall = append(wall, func() (interface{}, error) {
			return ft("1", 500)
		})
		wall = append(wall, func() (interface{}, error) {
			return ft("2", 100)
		})
		wall = append(wall, func() (interface{}, error) {
			return ft("3", 200)
		})
		wall = append(wall, panicFunction)

		r, s := SomeTimed(150, 4, wall...)
		if !cmp.Equal(r, want) || s != wState {
			t.Errorf("TimeoutExecution is not correct. Want\n%v %t\n got\n%v %t", want, wState, r, s)
		}
	})
}
