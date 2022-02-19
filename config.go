package main

import (
	"github.com/BurntSushi/toml"
)

type Config struct {
	Groups         []string            `toml:"groups"`
	GroupOverrides map[string][]string `toml:"group_overrides"`
}

func LoadConfig(fn string) (c Config, err error) {
	_, err = toml.DecodeFile(fn, &c)

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
