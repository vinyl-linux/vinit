package main

import (
	"reflect"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	s, err := New("testdata/services")
	if err != nil {
		t.Errorf("unexpected error %#v", err)
	}

	t.Run("groups and service mappings are correct", func(t *testing.T) {
		got := s.groupsServices
		expect := map[string][]string{
			"system": []string{"app", "app-cronjob", "app-oneoff"},
		}

		if !reflect.DeepEqual(expect, got) {
			t.Errorf("expected\n%#v\n\nreceived%#v", expect, got)
		}
	})
}

func TestSupervisor_RunShell_RunsWithoutPanic(t *testing.T) {
	s, err := New("testdata/services")
	if err != nil {
		t.Errorf("unexpected error %#v", err)
	}

	defer func() {
		err := recover()
		if err != nil {
			t.Errorf("unexpected error: %#v", err)
		}
	}()

	go s.RunShell()

	// give the shell time to do stuff
	time.Sleep(time.Millisecond * 100)

	s.restartShell = false
}
