package main

import (
	"bytes"
	"io"
	"reflect"
	"testing"
	"time"
)

func TestNewLogger(t *testing.T) {
	for _, test := range []struct {
		name        string
		kmesgFile   string
		expectError bool
	}{
		{"Happy path", "testdata/kmesg", false},
		{"Unable to open kmsg", "/this/path/does/not/exist/i/bloody/well/hope", true},
	} {
		t.Run(test.name, func(t *testing.T) {
			_, err := NewLogger(test.kmesgFile)
			if test.expectError && err == nil {
				t.Errorf("expected error, received none")
			} else if !test.expectError && err != nil {
				t.Errorf("unexpected error %#v", err)
			}

		})
	}
}

func TestLogger_addLog(t *testing.T) {
	oldKmesgF := kmesgF
	defer func() {
		kmesgF = oldKmesgF
	}()

	l := Logger{
		c: make(chan string),
	}

	go l.Start()

	for _, test := range []struct {
		name   string
		msg    string
		args   []interface{}
		expect string
	}{
		{"no kvs", "hello, world!", []interface{}{}, "vinit test: \"hello, world!\""},
		{"erroneous, single arg only writes msg", "hello, world!", []interface{}{"foo"}, "vinit test: \"hello, world!\""},
		{"odd args (gt 1) skips last arg", "hello, world!", []interface{}{"foo", 123, nil, 455, "abc"}, "vinit test: \"hello, world!\", foo=\"123\", <nil>=\"455\""},
		{"many, even args sets all", "hello, world!", []interface{}{"foo", 123, nil, 455}, "vinit test: \"hello, world!\", foo=\"123\", <nil>=\"455\""},
	} {
		t.Run(test.name, func(t *testing.T) {
			l.f = &bytes.Buffer{}
			l.Buffer = make([]string, maxLogLines)

			l.addLog("test", test.msg, test.args...)

			// Let buffer sync
			time.Sleep(time.Millisecond * 100)

			buf := new(bytes.Buffer)
			io.Copy(buf, l.f)

			got := buf.String()
			if l.Buffer[0] != got {
				t.Errorf("output mismatch; l.Buffer[0]: %q, l.f.String(): %q", l.Buffer[0], got)
			}

			if test.expect != got {
				t.Errorf("expected %q, got %q", test.expect, got)
			}
		})
	}
}

func TestLogger_Buffer(t *testing.T) {
	maxLogLines = 10
	l := Logger{
		Buffer: make([]string, maxLogLines),

		c: make(chan string),
		f: &bytes.Buffer{},
	}

	go l.Start()

	for i := 0; i < 10; i++ {
		l.addLog("test", "iter", "i", i)
	}

	// sync
	time.Sleep(time.Millisecond * 100)

	t.Run("new lines go to front of buffer", func(t *testing.T) {
		got := l.Buffer
		expect := []string{"vinit test: \"iter\", i=\"9\"", "vinit test: \"iter\", i=\"8\"", "vinit test: \"iter\", i=\"7\"", "vinit test: \"iter\", i=\"6\"", "vinit test: \"iter\", i=\"5\"", "vinit test: \"iter\", i=\"4\"", "vinit test: \"iter\", i=\"3\"", "vinit test: \"iter\", i=\"2\"", "vinit test: \"iter\", i=\"1\""}

		if !reflect.DeepEqual(expect, got) {
			t.Errorf("expected\n%#v\n\nreceived\n%#v", expect, got)
		}
	})
}
