package main

import (
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	s, _ := New("testdata/services")

	t.Run("groups and service mappings are correct", func(t *testing.T) {
		got := s.groupsServices
		expect := map[string][]string{
			"":       []string{"broken"},
			"system": []string{"app", "app-cronjob", "app-oneoff"},
		}

		if !reflect.DeepEqual(expect, got) {
			t.Errorf("expected\n%#v\n\nreceived%#v", expect, got)
		}
	})
}

func TestNew_IncorrectConfigs(t *testing.T) {
	_, err := New("testdata/mixed-status-services")
	if err == nil {
		t.Fatal("expected error, received none")
	}

	expect := `the following error(s) occurred parsing configs:
01-broken-app: invalid signal "BAT_SIGNAL"
`

	if expect != err.Error() {
		t.Errorf("expected\n%s\n\nreceived\n%s", expect, err.Error())
	}
}
