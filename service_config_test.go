package main

import (
	"testing"
)

func TestServiceConfig(t *testing.T) {
	for _, test := range []struct {
		name        string
		fn          string
		expectError bool
	}{
		{"service happy path", "testdata/services/00-app/.config.toml", false},
		{"cronjob happy path", "testdata/services/00-app-cronjob/.config.toml", false},
		{"oneoff happy path", "testdata/services/00-app-oneoff/.config.toml", false},

		// edge cases and errors
		{"invalid type", "testdata/erroring/wrong-type.toml", true},
		{"missing cron schedule", "testdata/erroring/missing-cron.toml", true},
		{"invalid cron schedule", "testdata/erroring/invalid-cron.toml", true},
		{"missing oneoff", "testdata/erroring/missing-oneoff.toml", true},
		{"missing group name", "testdata/erroring/missing-groupname.toml", true},
		{"invalid signal errors out", "testdata/erroring/invalid-signal.toml", true},
		{"missing args is fine", "testdata/successing/missing-args.toml", false},
		{"missing user sets user to root", "testdata/successing/missing-user.toml", false},
		{"empty validcodes gets a default", "testdata/successing/empty-validcodes.toml", false},
		{"empty reload signal gets a default", "testdata/successing/empty-reloadsignal.toml", false},

		// minimal viable configs
		{"minimal viable service", "testdata/mvs/service.toml", false},
		{"minimal viable cronjob", "testdata/mvs/cron.toml", false},
		{"minimal viable oneoff", "testdata/mvs/oneoff.toml", false},
	} {
		t.Run(test.name, func(t *testing.T) {
			_, err := LoadServiceConfig(test.fn)

			if test.expectError && err == nil {
				t.Errorf("expected error, received none")
			} else if !test.expectError && err != nil {
				t.Errorf("unexpected error %#v", err)
			}
		})
	}
}
