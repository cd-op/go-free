package metadata_test

import (
	"fmt"
	"testing"

	"cdop.pt/go/free/platepipe/metadata"
	. "cdop.pt/go/open/assertive"
)

func TestHasMetadata(t *testing.T) {
	cases := []struct {
		buf []byte
		ret bool
		pos int
	}{
		{[]byte(""), false, 0},
		{[]byte("test"), false, 0},
		{[]byte("test\ntest"), false, 0},

		{[]byte(" barekey = 'value'\n\n"), false, 0},

		{[]byte("barekey = 'value'\n\n"), true, 19},
		{[]byte("1barekey = 'value'\n\n"), true, 20},
		{[]byte("'key' = 'value'\n\n"), true, 17},
		{[]byte(`"key" = 'value'` + "\n\n"), true, 17},

		{[]byte("barekey = 'value'\n\nother text"), true, 19},
		{[]byte("barekey = 'value'\n\r\nother text"), true, 20},
	}

	for _, c := range cases {
		ret, pos := metadata.IsPresent(c.buf)

		if ret != c.ret || pos != c.pos {
			t.Errorf("IsPresent(%s) returned %v, %v",
				string(c.buf), ret, pos)
		}
	}
}

func TestFromTomlBuffer(t *testing.T) {
	t.Run("valid metadata", func(t *testing.T) {
		ret, err := metadata.FromTomlBuffer([]byte("strkey='value'\nintkey = 10"))
		Want(t, err == nil)
		Want(t, fmt.Sprint(ret) == fmt.Sprint(map[string]any{
			"strkey": "value",
			"intkey": 10,
		}))
	})

	t.Run("invalid metadata", func(t *testing.T) {
		ret, err := metadata.FromTomlBuffer([]byte("strkey: 'value'\nintkey = 10"))
		Need(t, err != nil)
		Want(t, err.Error() ==
			"toml: line 1: expected '.' or '=', but got ':' instead")
		Want(t, fmt.Sprint(ret) == fmt.Sprint(map[string]any{}))
	})
}
