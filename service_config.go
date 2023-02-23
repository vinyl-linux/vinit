package main

import (
	"fmt"
	"os"
	"os/user"
	"strconv"
	"syscall"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/google/shlex"
	"github.com/robfig/cron/v3"
	"golang.org/x/sys/unix"
)

const (
	ServiceType_Service ServiceType = iota
	ServiceType_Cron
	ServiceType_Oneoff
)

// ServiceType provides an enum type to track the type of service
// we're dealing with; namely:
//
//  1. ServiceType_Service, represented by "service" in config. Long running, to be restarted
//  2. ServiceType_Cron, represented by "cron" in config. Runs on schedule, expects to finish
//  3. ServiceType_Oneoff, represented by "oneoff" in config. Runs once on boot.
type ServiceType int8

// UnmarshalText provides the Unmarshal interface for ServiceType
func (s *ServiceType) UnmarshalText(text []byte) (err error) {
	t := string(text)

	switch t {
	case "service":
		*s = ServiceType_Service
	case "cron":
		*s = ServiceType_Cron
	case "oneoff":
		*s = ServiceType_Oneoff
	default:
		err = fmt.Errorf("invalid type %q; must be in set (%q,%q,%q)",
			t, "service", "cron", "oneoff")
	}

	return
}

// ReloadSignal holds an os.Signal which is sent to a process on `vinitctl reload process`
type ReloadSignal struct {
	s os.Signal
}

// UnmarshalText provides the Unmarshal interface for ReloadSignal.
//
// The default signal is SIGHUP.
func (r *ReloadSignal) UnmarshalText(text []byte) (err error) {
	if len(text) == 0 {
		r.s = syscall.SIGHUP

		return nil
	}

	s := unix.SignalNum(string(text))
	if s == 0 {
		return fmt.Errorf("invalid signal %q", string(text))
	}

	r.s = s

	return nil
}

// Args are the arguments set for a service.
type Args []string

// UnmarshalText provides the Unmarshal interface
//
// This is because we accept args as a string, to allow users to
// not have to worry about arrays of args, and all the other tedious
// things people hate
func (a *Args) UnmarshalText(text []byte) (err error) {
	ss, err := shlex.Split(string(text))
	if err != nil {
		return
	}

	*a = ss

	return
}

// User is a struct containing a service's defined User and Group.
//
// A User.User can either be a username, such as "root", or a uid, such
// as "1".
//
// Ditto User.Group.
type User struct {
	User  string `toml:"user"`
	Group string `toml:"group"`
}

// Uid returns an int64 of the specified user's uid.
//
// It returns an error if the user doesn't exist
func (u User) Uid() (uid int64, err error) {
	id, err := user.Lookup(u.User)
	if err != nil {
		return u.uidFromID()
	}

	return strconv.ParseInt(id.Uid, 10, 32)
}

func (u User) uidFromID() (uid int64, err error) {
	id, err := user.LookupId(u.User)
	if err != nil {
		return
	}

	return strconv.ParseInt(id.Uid, 10, 32)
}

// Gid returns an int64 of the specified user's gid.
//
// It returns an error if the user doesn't exist
func (u User) Gid() (uid int64, err error) {
	id, err := user.LookupGroup(u.Group)
	if err != nil {
		return u.gidFromID()
	}

	return strconv.ParseInt(id.Gid, 10, 32)
}

func (u User) gidFromID() (uid int64, err error) {
	id, err := user.LookupGroupId(u.Group)
	if err != nil {
		return
	}

	return strconv.ParseInt(id.Gid, 10, 32)
}

// Grouping provides a way of giving a service a named group
// which allows people to oder groups
type Grouping struct {
	GroupName string `toml:"name"`
}

// Schedule wraps a cron schedule so we can write an unmarshaller
type Schedule struct {
	cron.Schedule
}

// UnmarshalText implements the Unmarshal interface
func (s *Schedule) UnmarshalText(text []byte) (err error) {
	sched, err := cron.ParseStandard(string(text))
	if err != nil {
		return
	}

	s.Schedule = sched

	return
}

// Cron holds specific configs used just by services of type
// ServiceType_Cron
type Cron struct {
	Schedule Schedule `toml:"schedule"`
}

// Oneoff holds specific configs used just by services of type
// ServiceType_Oneoff
type Oneoff struct {
	ValidCodes []int `toml:"valid_exit_codes"`
}

// Success returns a bool based on whether a Oneoff service
// completed successfully
func (o Oneoff) Success(exitCode int) bool {
	for _, c := range o.ValidCodes {
		if c == exitCode {
			return true
		}
	}

	return false
}

// Command holds extra arguments and config for the process
// started for the service
type Command struct {
	Args         Args `toml:"args"`
	IgnoreOutput bool `toml:"ignore_output"`
}

// ServiceConfig configures a specific service, and includes
// args, and types, and all that stuff
type ServiceConfig struct {
	Type         ServiceType   `toml:"type"`
	ReloadSignal *ReloadSignal `toml:"reload_signal"`
	User         User          `toml:"user"`
	Grouping     Grouping      `toml:"grouping"`
	Cron         *Cron         `toml:"cron,omitempty"`
	Oneoff       *Oneoff       `toml:"oneoff,omitemoty"`
	Command      Command       `toml:"command"`
}

// LoadServiceConfig decodes a toml file
func LoadServiceConfig(fn string) (s ServiceConfig, err error) {
	_, err = toml.DecodeFile(fn, &s)
	if err != nil {
		return
	}

	switch s.Type {
	case ServiceType_Service:
	case ServiceType_Cron:
		if s.Cron == nil || s.Cron.Schedule.Schedule == nil || (s.Cron.Schedule.Next(time.Now()) == time.Time{}) {
			err = fmt.Errorf("invalid cron schedule")
		}
	case ServiceType_Oneoff:
		if s.Oneoff == nil {
			err = fmt.Errorf("missing oneoff config")

			return
		}

		if len(s.Oneoff.ValidCodes) == 0 {
			s.Oneoff.ValidCodes = []int{0}
		}
	}

	if s.Grouping.GroupName == "" {
		err = fmt.Errorf("missing grouping name")

		return
	}

	if s.User.User == "" {
		s.User.User = "root"
	}

	if s.User.Group == "" {
		s.User.Group = s.User.User
	}

	if s.ReloadSignal == nil || s.ReloadSignal.s == nil {
		s.ReloadSignal = &ReloadSignal{
			s: syscall.SIGHUP,
		}
	}

	return
}
