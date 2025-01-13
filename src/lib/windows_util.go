package lib

import (
	"errors"
)

// Deprecated: stopped supporting windows as localhost.
func Check_service_status_windows(service string) (bool, error) {
	// m, err := mgr.Connect()
	// if err != nil {
	// 	return false, err
	// }
	// defer m.Disconnect()

	// s, err := m.OpenService(service)
	// if err != nil {
	// 	return false, err
	// }
	// defer s.Close()

	// status, err := s.Query()
	// if err != nil {
	// 	return false, err
	// }

	// if status.State == svc.Running {
	// 	return true, nil
	// } else {
	// 	return false, nil
	// }

	return false, errors.New("errr: Deprecated: stopped supporting windows as localhost")
}

// Deprecated: stopped supporting windows as localhost.
func Check_jazagent_status_windows(service string) (bool, error) {
	return Check_service_status_windows("Job Arranger Agent")
}

// Deprecated: stopped supporting windows as localhost.
func Restart_service_windows(service string) error {

	// m, err := mgr.Connect()
	// if err != nil {
	// 	return err
	// }
	// defer m.Disconnect()

	// s, err := m.OpenService(service)
	// if err != nil {
	// 	return err
	// }
	// defer s.Close()

	// status, err := s.Query()
	// if err != nil {
	// 	return err
	// }

	// if status.State != svc.Running {
	// 	err = s.Start()
	// }

	return errors.New("errr: Deprecated: stopped supporting windows as localhost")
}

// Deprecated: stopped supporting windows as localhost.
func Stop_service_windows(service string) error {

	// m, err := mgr.Connect()
	// if err != nil {
	// 	return err
	// }
	// defer m.Disconnect()

	// s, err := m.OpenService(service)
	// if err != nil {
	// 	return err
	// }
	// defer s.Close()

	// status, err := s.Control(svc.Stop)
	// if err != nil {
	// 	return err
	// }

	// index := 0
	// for status.State != svc.Stopped {
	// 	status, err = s.Query()
	// 	if err != nil {
	// 		return err
	// 	}
	// 	Spinner_log(index, Formatted_log(common.LOG_LEVEL_INFO, "Jobarg agent service is stopping."))

	// 	time.Sleep(2 * time.Second)
	// }

	// return err

	return errors.New("errr: Deprecated: stopped supporting windows as localhost")
}

// Deprecated: stopped supporting windows as localhost.
func Restart_jaz_agent_windows() error {
	return Restart_service_windows("Job Arranger Agent")
}

// Deprecated: stopped supporting windows as localhost.
func Stop_jaz_agent_windows() error {
	return Stop_service_windows("Job Arranger Agent")
}

// Deprecated: stopped supporting windows as localhost.
func Execute_cmd_window(command string) ([]byte, error) {
	// cmd := exec.Command("cmd", "/C", command)

	// return cmd.CombinedOutput()

	return nil, errors.New("errr: Deprecated: stopped supporting windows as localhost")
}

// Deprecated: stopped supporting windows as localhost.
func Execute_cmd_window_str(command string) (string, error) {
	output, err := Execute_cmd_window(command)
	return string(output), err
}

// To use this function, you must have jobarg_agentd default filepath.
//
// Deprecated: stopped supporting windows as localhost.
func Ja_set_agent_config_windows(key string, value string) error {
	// var lines []string
	// var updated bool
	// file_location := filepath.Join("C:\\", "Program Files", "Job Arranger", "Job Arranger Agent", "conf", "jobarg_agentd.conf")

	// file, err := os.OpenFile(file_location, os.O_RDWR, 0644)
	// if err != nil {
	// 	return err
	// }
	// defer file.Close()

	// scanner := bufio.NewScanner(file)

	// for scanner.Scan() {
	// 	line := scanner.Text()
	// 	line = strings.TrimSpace(line)

	// 	if strings.HasPrefix(line, key+"=") {
	// 		lines = append(lines, key+"="+value)
	// 		updated = true
	// 	} else {
	// 		lines = append(lines, line)
	// 	}
	// }

	// if !updated {
	// 	for index, l := range lines {
	// 		if strings.HasPrefix(l, "# "+key+"=") {
	// 			lines[index] = fmt.Sprintf("%s=%s", key, value)
	// 			updated = true
	// 		}
	// 	}
	// }

	// // Check for scanner errors
	// if err := scanner.Err(); err != nil {
	// 	return err
	// }

	// if updated {
	// 	// Seek to the beginning of the file to overwrite it
	// 	if _, err := file.Seek(0, 0); err != nil {
	// 		return err
	// 	}

	// 	// Truncate the file to remove any leftover content
	// 	if err := file.Truncate(0); err != nil {
	// 		return err
	// 	}

	// 	// Write the updated lines back to the file
	// 	for _, line := range lines {
	// 		if _, err := file.WriteString(line + "\n"); err != nil {
	// 			return err
	// 		}
	// 	}
	// }
	//
	// return nil

	return errors.New("errr: Deprecated: stopped supporting windows as localhost")
}

// Deprecated: stopped supporting windows as localhost.
func Cleanup_agent_windows() error {
	// dir := "C:\\Program Files\\Job Arranger\\Job Arranger Agent\\temp"
	// files, err := os.ReadDir(dir)
	// if err != nil {
	// 	return err
	// }

	// for _, file := range files {
	// 	path := filepath.Join(dir, file.Name())
	// 	if file.IsDir() {
	// 		// Recursively delete subdirectory
	// 		err = os.RemoveAll(path)
	// 	} else {
	// 		// Delete file
	// 		err = os.Remove(path)
	// 	}
	// 	if err != nil {
	// 		return err
	// 	}
	// }
	// return nil

	return errors.New("errr: Deprecated: stopped supporting windows as localhost")
}

// To use this function, your jobarranger agent's TmpDir must be default(TmpDir=C:\Program Files\Job Arranger\Job Arranger Agent\temp).
//
// Jobarg_cleanup() cleans jobarg-server and jobarg-agentd(windows) data.
// Since this is testcase utility funtion, you must use it in testcase function.
//
// Deprecated: stopped supporting windows as localhost.
func Jobarg_cleanup_windows() error {
	// if err := Stop_jaz_server(); err != nil {
	// 	return fmt.Errorf("failed to stop JAZ server: %w", err)
	// }
	// if err := Stop_jaz_agent_windows(); err != nil {
	// 	return fmt.Errorf("failed to stop JAZ agent: %w", err)
	// }
	// if _, err := DBexec("delete from ja_run_jobnet_table;"); err != nil {
	// 	return fmt.Errorf("failed to execute DB command: %w", err)
	// }
	// if err := Cleanup_agent_windows(); err != nil {
	// 	return fmt.Errorf("failed to clean up agent: %w", err)
	// }
	// if err := Restart_jaz_server(); err != nil {
	// 	return fmt.Errorf("failed to stop JAZ server: %w", err)
	// }
	// if err := Restart_jaz_agent_windows(); err != nil {
	// 	return fmt.Errorf("failed to stop JAZ server: %w", err)
	// }

	// return nil

	return errors.New("errr: Deprecated: stopped supporting windows as localhost")
}
