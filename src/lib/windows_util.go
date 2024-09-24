package lib

import (
	"time"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
)

func Check_service_status_windows(service string) (bool, error) {
	m, err := mgr.Connect()
	if err != nil {
		return false, err
	}
	defer m.Disconnect()

	s, err := m.OpenService(service)
	if err != nil {
		return false, err
	}
	defer s.Close()

	status, err := s.Query()
	if err != nil {
		return false, err
	}

	if status.State == svc.Running {
		return true, nil
	} else {
		return false, nil
	}
}

func Check_jazagent_status_windows(service string) (bool, error) {
	return Check_service_status_windows("Job Arranger Agent")
}

func Restart_service_windows(service string) error {

	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()

	s, err := m.OpenService(service)
	if err != nil {
		return err
	}
	defer s.Close()

	status, err := s.Query()
	if err != nil {
		return err
	}

	if status.State != svc.Running {
		err = s.Start()
	}

	return err
}

func Stop_service_windows(service string) error {

	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()

	s, err := m.OpenService(service)
	if err != nil {
		return err
	}
	defer s.Close()

	status, err := s.Control(svc.Stop)
	if err != nil {
		return err
	}

	index := 0
	for status.State != svc.Stopped {
		status, err = s.Query()
		if err != nil {
			return err
		}
		Spinner_log(index, Formatted_log(INFO, "Jobarg agent service is stopping."))

		time.Sleep(2 * time.Second)
	}

	return err
}

func Restart_jaz_agent_windows() error {
	return Restart_service_windows("Job Arranger Agent")
}

func Stop_jaz_agent_windows() error {
	return Stop_service_windows("Job Arranger Agent")
}
