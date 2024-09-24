package lib

import (
	"golang.org/x/sys/windows/svc/mgr"
)

func Resart_service_windows(service string) error {

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

	err = s.Start()

	return err
}

func Restart_jaz_agent_windows() error {
	return Resart_service_windows("Job Arranger Agent")
}
