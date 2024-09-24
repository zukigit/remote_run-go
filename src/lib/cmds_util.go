package lib

func Ssh_exec(command string) ([]byte, error) {
	session, err := Get_session(Common_client)
	if err != nil {
		return nil, err
	}
	defer session.Close()

	return session.Output(command)
}

func Ssh_exec_to_str(command string) (string, error) {
	output, err := Ssh_exec(command)
	return string(output), err
}
