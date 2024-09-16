package lib

func Ssh_exec(command string) ([]byte, error) {
	return ticket.Get_auth().Session.Output(command)
}

func Exec(command string) (string, error) {
	output, err := ticket.Get_auth().Session.Output(command)
	return string(output), err
}

func Restart_jaz_server() ([]byte, error) {
	return Ssh_exec("systemctl restart jobarg-server")
}

func Restart_docker() ([]byte, error) {
	return Ssh_exec("systemctl restart docker")
}

func Stop_docker() ([]byte, error) {
	return Ssh_exec("systemctl stop docker")
}

func Stop_jaz_server() ([]byte, error) {
	return Ssh_exec("systemctl stop jobarg-server")
}

func Start_jaz_server() ([]byte, error) {
	return Ssh_exec("systemctl start jobarg-server")
}
