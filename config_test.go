package main

import (
	"reflect"
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

func TestConfig_StartupScript(t *testing.T) {
	for _, test := range []struct {
		name        string
		fn          string
		expectCmd   string
		expectArgs  []string
		expectError bool
	}{
		{"Empty startup script loads defaults", "testdata/successing/empty-startupscript.toml", defaultStartupScript.cmd, defaultStartupScript.args, false},
		{"Undefined startup script loads defaults", "testdata/successing/undefined-startupscript.toml", defaultStartupScript.cmd, defaultStartupScript.args, false},
		{"Only setting cmd means no args", "testdata/services/.config.toml", "/bin/date", []string(nil), false},
		{"Setting cmd and args works correctly", "testdata/successing/full-startupscript.toml", "/sbin/agetty", []string{"tty1", "linux"}, false},
		{"Dodgy/ unexpected error fails accordingly", "testdata/erroring/dodgy-startupscript.toml", "", []string(nil), true},
	} {
		t.Run(test.name, func(t *testing.T) {
			c, err := LoadConfig(test.fn)

			if test.expectError && err == nil {
				t.Errorf("expected error, received none")
			} else if !test.expectError && err != nil {
				t.Errorf("unexpected error %#v", err)
			}

			if test.expectCmd != c.StartupScript.cmd {
				t.Errorf("expected %q, received %q", test.expectCmd, c.StartupScript.cmd)
			}

			if !reflect.DeepEqual(test.expectArgs, c.StartupScript.args) {
				t.Errorf("expected %#v, received %#v", test.expectArgs, c.StartupScript.args)
			}
		})
	}
}
