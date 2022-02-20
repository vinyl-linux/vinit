package main

import (
	"os"
	"path/filepath"
	"regexp"
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

func (s *Supervisor) StopAll() (err error) {
	for _, svc := range s.services {
		if svc.isRunning() {
			err = svc.Stop()
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
