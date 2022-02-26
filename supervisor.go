package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var (
	svcPrefix = regexp.MustCompile(`^\d.\-`)
)

type Supervisor struct {
	Config Config

	dir            string
	groupsServices map[string][]string
	services       map[string]*Service
	restartShell   bool
}

func New(dir string) (s *Supervisor, err error) {
	s = &Supervisor{
		dir: dir,
	}

	err = s.LoadConfigs()

	return
}

func (s *Supervisor) LoadConfigs() (err error) {
	s.Config, err = LoadConfig(filepath.Join(s.dir, ".config.toml"))
	if err != nil {
		return
	}

	groupsServices := make(map[string][]string)
	services := make(map[string]*Service)

	entries, err := os.ReadDir(s.dir)
	if err != nil {
		return
	}

	var svc *Service

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		svc, err = LoadService(filepath.Join(s.dir, entry.Name()))
		if err != nil {
			return
		}

		name := serviceName(entry.Name())
		groupName := s.Config.ReconcileOverride(name, svc.Config.Grouping.GroupName)
		_, ok := groupsServices[groupName]
		if !ok {
			groupsServices[groupName] = make([]string, 0)
		}

		groupsServices[groupName] = append(groupsServices[groupName], name)
		services[name] = svc
	}

	s.groupsServices = groupsServices
	s.services = services

	return
}

func (s *Supervisor) Start(name string, wait bool) error {
	return s.services[name].Start(wait)
}

func (s *Supervisor) Status(name string) (ServiceStatus, error) {
	return s.services[name].Status()
}

func (s *Supervisor) Stop(name string) error {
	return s.services[name].Stop()
}

func (s *Supervisor) Reload(name string) error {
	return s.services[name].Reload()
}

func (s *Supervisor) StartAll() {
	var err error

	for _, group := range s.Config.Groups {
		services, ok := s.groupsServices[group]
		if !ok {
			sugar.Errorw("group either has no services or does not exist",
				"group", group,
			)

			continue
		}

		for _, service := range services {
			sugar.Infow("starting",
				"group", group,
				"service", service,
			)

			err = s.Start(service, true)
			if err != nil {
				sugar.Errorw("failed!",
					"group", group,
					"service", service,
					"error", err.Error(),
				)

				continue
			}

			sugar.Infow("started!",
				"group", group,
				"service", service,
			)

		}
	}
}

// StopAll does the opposite of StartAll; it reverses the order of
// s.Config.Groups, then reverses the order of those services in order
// to stop them all
func (s *Supervisor) StopAll() (err error) {
	var svc *Service

	for _, group := range reverse(s.Config.Groups) {
		for _, svcName := range reverse(s.groupsServices[group]) {
			svc = s.services[svcName]

			if svc == nil || !svc.isRunning() {
				continue
			}

			err = s.Stop(svcName)
			if err != nil {
				return
			}
		}
	}

	return
}

func (s *Supervisor) RunShell() {
	sc := s.Config.StartupScript

	// Keep restarting shell if it crashes
	//
	// This is used during shutdown, for instance,
	// to stop the inital shell constantly restarting
	s.restartShell = true

	for s.restartShell {
		c := exec.Command(sc.cmd, sc.args...) //#nosec: G204
		c.Env = os.Environ()
		c.Dir = "/"

		c.Stdin = os.Stdin
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr

		err := c.Run()
		if err != nil {
			sugar.Errorw("shell restarting",
				"cmd", sc.cmd,
				"args", strings.Join(sc.args, " "),
				"error", err.Error(),
			)
		}

		time.Sleep(time.Second)
	}
}

func serviceName(s string) string {
	if !svcPrefix.Match([]byte(s)) {
		return s
	}

	return svcPrefix.ReplaceAllString(s, "")
}

func reverse(in []string) (out []string) {
	out = make([]string, len(in))
	copy(out, in)

	for i := len(out)/2 - 1; i >= 0; i-- {
		opp := len(out) - 1 - i
		out[i], out[opp] = out[opp], out[i]
	}

	return
}
