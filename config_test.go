package main

import (
	"testing"
)

func TestConfig_ReconcileOverride(t *testing.T) {
	defaultGroup := "testing"
	overrideGroup := "my-group"
	svc := "my-super-service"

	for _, test := range []struct {
		name   string
		c      Config
		expect string
	}{
		{"empty config returns inital group", Config{}, defaultGroup},
		{"where no override exists, initial group returned", Config{GroupOverrides: map[string][]string{}}, defaultGroup},
		{"where no relevant override exists, initial group returned", Config{GroupOverrides: map[string][]string{overrideGroup: []string{"another-service"}}}, defaultGroup},
		{"override is returned", Config{GroupOverrides: map[string][]string{overrideGroup: []string{svc}}}, overrideGroup},
	} {
		t.Run(test.name, func(t *testing.T) {
			grp := test.c.ReconcileOverride(svc, defaultGroup)
			if test.expect != grp {
				t.Errorf("expected %q, received %q", test.expect, grp)
			}
		})
	}
}
