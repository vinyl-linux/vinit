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

type ServiceType int8

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

type ReloadSignal struct {
	s os.Signal
}

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

type Args []string

func (a *Args) UnmarshalText(text []byte) (err error) {
	ss, err := shlex.Split(string(text))
	if err != nil {
		return
	}

	*a = ss

	return
}

type User struct {
	User  string `toml:"user"`
	Group string `toml:"group"`
}

func (u User) Uid() (uid int64, err error) {
	id, err := user.Lookup(u.User)
	if err != nil {
		return
	}

	return strconv.ParseInt(id.Uid, 10, 32)
}

func (u User) Gid() (uid int64, err error) {
	id, err := user.LookupGroup(u.Group)
	if err != nil {
		return
	}

	return strconv.ParseInt(id.Gid, 10, 32)
}

type Grouping struct {
	GroupName string `toml:"name"`
}

type Schedule struct {
	cron.Schedule
}

func (s *Schedule) UnmarshalText(text []byte) (err error) {
	sched, err := cron.ParseStandard(string(text))
	if err != nil {
		return
	}

	s.Schedule = sched

	return
}

type Cron struct {
	Schedule Schedule `toml:"schedule"`
}

type Oneoff struct {
	ValidCodes []int `toml:"valid_exit_codes"`
}

func (o Oneoff) Success(exitCode int) bool {
	for _, c := range o.ValidCodes {
		if c == exitCode {
			return true
		}
	}

	return false
}

type Command struct {
	Args Args `toml:"args"`
}

type ServiceConfig struct {
	Type         ServiceType   `toml:"type"`
	ReloadSignal *ReloadSignal `toml:"reload_signal"`
	User         User          `toml:"user"`
	Grouping     Grouping      `toml:"grouping"`
	Cron         *Cron         `toml:"cron,omitempty"`
	Oneoff       *Oneoff       `toml:"oneoff,omitemoty"`
	Command      Command       `toml:"command"`
}

func LoadServiceConfig(fn string) (s ServiceConfig, err error) {
	_, err = toml.DecodeFile(fn, &s)

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
