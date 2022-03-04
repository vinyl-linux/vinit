package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"
)

type ServiceStatus struct {
	Running    bool
	Pid        int
	ExitStatus int
	StartTime  time.Time
	EndTime    time.Time
	Success    bool
	Error      error
}

type Service struct {
	Config ServiceConfig
	Env    EnvVars

	// The following only need to be set once, no matter
	// how often a service is restarted
	uid    uint32
	gid    uint32
	bin    string
	wd     string
	logdir string

	// The following get set on Service.Start()
	status ServiceStatus
	proc   *exec.Cmd

	// Service.Stop() sets this 'desired state' value
	// before killing the process. When the Start() loop
	// kicks in, it checks whether this value is true
	// prior to restarting
	desiredRunning bool

	// isDirty is set to true if a call to Supervisor.LoadConfigs
	// fails and so new config hasn't been picked up.
	//
	// we can assume that this means there's a new config that doesn't
	// look right, since the existence of this service means it must
	// have been right first time.
	//
	// but, in any case, isDirty is a flag that means 'Something went
	// wrong parsing this service config, go check the logs'
	isDirty bool
}

func LoadService(dir string) (s *Service, err error) {
	s = new(Service)

	s.Config, err = LoadServiceConfig(filepath.Join(dir, ".config.toml"))
	if err != nil {
		return
	}

	s.Env, err = LoadEnvVars(filepath.Join(dir, "environment"))
	if err != nil {
		return
	}

	overridesFn := filepath.Join(dir, "environment_overrides")
	if _, err = os.Stat(overridesFn); err == nil {
		var overrides EnvVars

		overrides, err = LoadEnvVars(overridesFn)
		if err != nil {
			return
		}

		s.Env = append(s.Env, overrides...)
	}

	err = nil

	uid, err := s.Config.User.Uid()
	if err != nil {
		return
	}
	s.uid = uint32(uid)

	gid, err := s.Config.User.Gid()
	if err != nil {
		return
	}
	s.gid = uint32(gid)

	s.bin = filepath.Join(dir, "bin")
	err = s.validateBin()
	if err != nil {
		err = fmt.Errorf("file %s must exist and be an executable", s.bin)

		return
	}

	s.wd = filepath.Join(dir, "wd")
	s.logdir = filepath.Join(dir, "logs")

	return
}

func (s *Service) Start(wait bool) (err error) {
	if s.isRunning() {
		return fmt.Errorf("service is already running")
	}

	defer func() {
		s.proc = nil
	}()

	s.status.Running = true
	s.status = ServiceStatus{
		StartTime: time.Now(),
	}

	if s.Config.Type == ServiceType_Oneoff && wait {
		s.status.Error = s.start()
		s.status.Success = s.Config.Oneoff.Success(s.status.ExitStatus)

		return s.status.Error
	}

	go func() {
		var cont = true
		for cont {
			s.status.Running = true

			s.status.Error = s.start()

			cont = s.desiredRunning

			if s.Config.Type == ServiceType_Oneoff {
				s.status.Success = s.Config.Oneoff.Success(s.status.ExitStatus)

				return
			}

		}
	}()

	return nil
}

func (s *Service) start() (err error) {
	s.proc = exec.Command(s.bin, s.Config.Command.Args...) // #nosec G204
	s.proc.Env = s.Env
	s.proc.Dir = s.wd
	s.proc.SysProcAttr = &syscall.SysProcAttr{}
	s.proc.SysProcAttr.Credential = &syscall.Credential{Uid: uint32(s.uid), Gid: uint32(s.gid)}

	if !s.Config.Command.IgnoreOutput {
		for _, f := range []func() error{
			s.mkLogdir,
			s.streamStdout,
			s.streamStderr,
		} {
			err = f()
			if err != nil {
				continue
			}
		}
	}

	err = s.proc.Start()
	if err != nil {
		return
	}

	s.status.Pid = s.proc.Process.Pid

	err = s.proc.Wait()

	s.status.Running = false
	s.status.EndTime = time.Now()
	s.status.ExitStatus = s.proc.ProcessState.ExitCode()

	s.proc = nil

	return
}

func (s *Service) Stop() (err error) {
	if !s.isRunning() {
		return fmt.Errorf("service is not running")
	}

	s.desiredRunning = false
	s.status.EndTime = time.Now()

	err = s.proc.Process.Kill()
	if err != nil {
		return
	}

	s.status.Running = false
	s.status.ExitStatus = s.proc.ProcessState.ExitCode()

	return
}

func (s *Service) Status() (status ServiceStatus, err error) {
	return s.status, nil
}

func (s *Service) Reload() (err error) {
	if !s.isRunning() {
		return fmt.Errorf("service is not running")
	}

	return s.proc.Process.Signal(s.Config.ReloadSignal.s)
}

func (s Service) isRunning() bool {
	return s.proc != nil && s.proc.Process != nil
}

func (s *Service) mkLogdir() error {
	return os.MkdirAll(s.logdir, 0700)
}

func (s *Service) streamStdout() (err error) {
	s.proc.Stdout, err = os.OpenFile(filepath.Join(s.logdir, "stdout"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		s.proc.Stdout = io.Discard
	}

	return
}

func (s *Service) streamStderr() (err error) {
	s.proc.Stderr, err = os.OpenFile(filepath.Join(s.logdir, "stderr"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		s.proc.Stderr = io.Discard
	}

	return
}

func (s Service) validateBin() (err error) {
	f, err := os.Stat(s.bin)
	if err != nil {
		return fmt.Errorf("could not open file %s", s.bin)
	}

	if !f.Mode().IsRegular() {
		return fmt.Errorf("file %s is not a regular file", s.bin)
	}

	if f.Mode()&0111 == 0 {
		return fmt.Errorf("file %s is not executable", s.bin)
	}

	return
}
