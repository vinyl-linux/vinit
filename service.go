package main

import (
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
	s.wd = filepath.Join(dir, "wd")
	s.logdir = filepath.Join(dir, "logs")

	return
}

func (s *Service) Start() (err error) {
	s.status = ServiceStatus{
		StartTime: time.Now(),
	}

	go func() {
		for {
			s.status.Running = true
			s.status.Error = s.start()
			if s.Config.Type == ServiceType_Oneoff {
				s.status.Success = s.Config.Oneoff.Success(s.status.ExitStatus)

				return
			}
		}
	}()

	return nil
}

func (s *Service) start() (err error) {
	s.proc = exec.Command(s.bin, s.Config.Command.Args...)
	s.proc.Env = s.Env
	s.proc.Dir = s.wd
	s.proc.SysProcAttr = &syscall.SysProcAttr{}
	s.proc.SysProcAttr.Credential = &syscall.Credential{Uid: uint32(s.uid), Gid: uint32(s.gid)}

	for _, f := range []func() error{
		s.mkLogdir,
		s.streamStdout,
		s.streamStderr,
	} {
		err = f()
		if err != nil {
			return
		}
	}

	return s.runloop()
}

func (s *Service) Stop() (err error) {
	err = s.proc.Process.Kill()
	if err != nil {
		return
	}

	s.status.Running = false
	s.status.EndTime = time.Now()
	s.status.ExitStatus = s.proc.ProcessState.ExitCode()

	return
}

func (s *Service) Status() (status ServiceStatus, err error) {
	return s.status, nil
}

func (s *Service) mkLogdir() error {
	return os.MkdirAll(s.logdir, 0700)
}

func (s *Service) streamStdout() (err error) {
	s.proc.Stdout, err = os.OpenFile(filepath.Join(s.logdir, "stdout"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	return
}

func (s *Service) streamStderr() (err error) {
	s.proc.Stderr, err = os.OpenFile(filepath.Join(s.logdir, "stderr"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	return
}

func (s *Service) runloop() (err error) {
	err = s.proc.Start()
	if err != nil {
		return
	}

	s.status.Pid = s.proc.Process.Pid

	err = s.proc.Wait()
	s.proc.Process.Release()

	s.status.Running = false
	s.status.EndTime = time.Now()
	s.status.ExitStatus = s.proc.ProcessState.ExitCode()

	return
}
