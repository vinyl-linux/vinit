package main

import (
	"reflect"
	"testing"
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
