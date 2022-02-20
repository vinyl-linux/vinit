package main

import (
	"github.com/BurntSushi/toml"
	"github.com/google/shlex"
)

var (
	defaultStartupScript = &StartupScript{
		cmd: "/sbin/agetty",
		args: []string{
			"-L",
			"-8",
			"--autologin",
			"root",
			"115200",
			"tty1",
			"linux",
		},
	}
)

type StartupScript struct {
	cmd  string
	args []string
}

func (s *StartupScript) UnmarshalText(text []byte) (err error) {
	if len(text) == 0 {
		*s = *defaultStartupScript

		return
	}

	ss, err := shlex.Split(string(text))
	if err != nil {
		return
	}

	switch len(ss) {
	case 0:
		*s = *defaultStartupScript

	case 1:
		s.cmd = ss[0]

	default:
		s.cmd = ss[0]
		s.args = ss[1:]
	}

	return
}

type Config struct {
	Groups         []string            `toml:"groups"`
	GroupOverrides map[string][]string `toml:"group_overrides"`
	StartupScript  *StartupScript      `toml:"startup_script"`
}

func LoadConfig(fn string) (c Config, err error) {
	_, err = toml.DecodeFile(fn, &c)
	if err != nil {
		return
	}

	if c.StartupScript == nil {
		c.StartupScript = defaultStartupScript
	}

	return
}

// HasOverride returns an optional 'override group' for a service.
//
// An override group is a local configuration option which allows the owner
// of a system to override which group a service belongs to, even (effectively)
// creating a brand new group and moving services into it.
//
// This allows users to tweak the order in which boot
func (c Config) ReconcileOverride(svc, initialGroup string) (group string) {
	var members []string

	for group, members = range c.GroupOverrides {
		if contains(members, svc) {
			return
		}
	}

	return initialGroup
}

func contains(ss []string, s string) bool {
	for _, elem := range ss {
		if elem == s {
			return true
		}
	}

	return false
}
