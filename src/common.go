package src

import "zukigit/remote_run-go/src/dao"

var Auth *dao.Auth

func Set_common_auth(auth *dao.Auth) {
	Auth = auth
}
