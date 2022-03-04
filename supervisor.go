package main

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	svcPrefix = regexp.MustCompile(`^\d.\-`)
)

type Supervisor struct {
	Config Config

	dir            string
	groupsServices map[string][]string
	services       map[string]*Service
}

type ConfigParseError struct {
	errors map[string]error
}

func (c *ConfigParseError) Append(svc string, err error) {
	if c.errors == nil {
		c.errors = make(map[string]error)
	}

	c.errors[svc] = err
}

func (c ConfigParseError) Error() string {
	out := new(strings.Builder)
	out.WriteString("the following error(s) occurred parsing configs:\n")
	for svc, err := range c.errors {
		out.WriteString(svc + ": " + err.Error() + "\n")
	}

	return out.String()
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
	cpe := new(ConfigParseError)

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		svc, err = LoadService(filepath.Join(s.dir, entry.Name()))
		if err != nil {
			// log error, recover
			svc.isDirty = true
			cpe.Append(entry.Name(), err)

			continue
		}

		name := serviceName(entry.Name())
		groupName := s.Config.ReconcileOverride(name, svc.Config.Grouping.GroupName)
		_, ok := groupsServices[groupName]
		if !ok {
			groupsServices[groupName] = make([]string, 0)
		}

		groupsServices[groupName] = append(groupsServices[groupName], name)

		// If this service already exists/ has some state then copy it over
		// (so we don't lose running state)
		oldSvc := s.services[name]
		if oldSvc != nil {
			svc.status = oldSvc.status
		}

		services[name] = svc
	}

	if len(cpe.errors) > 0 {
		return cpe
	}

	// Only assign new services and groups when everything loads,
	// rather than accidentally returning broken state
	s.groupsServices = groupsServices
	s.services = services

	return
}

func (s *Supervisor) Start(name string, wait bool) error {
	svc, ok := s.services[name]
	if !ok {
		return errServiceNotExist
	}

	return svc.Start(wait)
}

func (s *Supervisor) Status(name string) (ServiceStatus, error) {
	svc, ok := s.services[name]
	if !ok {
		return ServiceStatus{}, errServiceNotExist
	}

	return svc.Status()
}

func (s *Supervisor) Stop(name string) error {
	svc, ok := s.services[name]
	if !ok {
		return errServiceNotExist
	}

	return svc.Stop()
}

func (s *Supervisor) Reload(name string) error {
	svc, ok := s.services[name]
	if !ok {
		return errServiceNotExist
	}

	return svc.Reload()
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
