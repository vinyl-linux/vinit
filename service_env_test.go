package main

import (
	"reflect"
	"testing"
)

func TestLoadEnvVars(t *testing.T) {
	for _, test := range []struct {
		name         string
		fn           string
		expect       EnvVars
		expectErrors bool
	}{
		{"file exists, has contents", "testdata/services/00-app/environment", EnvVars{"PATH=/bin:/sbin"}, false},
		{"file doesn't exist", "testdata/services/00-app-cronjob/environment", EnvVars{}, false},
		{"file exists, empty", "testdata/services/00-app-oneoff/environment", EnvVars{}, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			got, err := LoadEnvVars(test.fn)

			if test.expectErrors && err == nil {
				t.Error("expected error, received nothing")
			} else if !test.expectErrors && err != nil {
				t.Errorf("unexpected error: %#v", err)
			}

			if !reflect.DeepEqual(test.expect, got) {
				t.Errorf("expected %#v, received %#v", test.expect, got)
			}
		})
	}
}
