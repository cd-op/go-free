package variables_test

import (
	"fmt"
	"testing"

	"cdop.pt/go/free/platepipe/variables"
	. "cdop.pt/go/open/assertive"
)

type m = map[string]any

func TestCoalesce(t *testing.T) {
	cases := []struct {
		args []m
		ret  m
	}{
		{
			[]m{},
			m{},
		},
		{
			[]m{{"k": "v"}},
			m{"k": "v"},
		},
		{
			[]m{{"k": 1}, {"k": 2}},
			m{"k": 1},
		},
		{
			[]m{{"k1": 1}, {"k2": 2}},
			m{"k1": 1, "k2": 2},
		},
		{
			[]m{
				{"k1": 1, "k2": 2},
				{"k2": "ignored", "k3": 3},
				{"k3": "ignored", "k4": 4},
			},
			m{"k1": 1, "k2": 2, "k3": 3, "k4": 4},
		},
	}

	for _, c := range cases {
		ret := variables.Coalesce(c.args...)
		Want(t, fmt.Sprint(ret) == fmt.Sprint(c.ret))
	}
}
