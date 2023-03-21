package errno

import (
	"errors"
	"fmt"
	"log"
)

const (
	SuccessCode                = 0
	ServiceErrCode             = 10001
	ParamErrCode               = 10002
	UserAlreadyExistErrCode    = 10003
	AuthorizationFailedErrCode = 10004
	RpcConnectErrCode          = 10005
	UploadVideoErrCode         = 10006
	ErrDatabaseCode            = 10007
	debug                      = false
)

var (
	Success                = NewErrNo(SuccessCode, "Success")
	ServiceErr             = NewErrNo(ServiceErrCode, "Service is unable to start successfully")
	ParamErr               = NewErrNo(ParamErrCode, "Wrong Parameter has been given")
	UserAlreadyExistErr    = NewErrNo(UserAlreadyExistErrCode, "User already exists")
	AuthorizationFailedErr = NewErrNo(AuthorizationFailedErrCode, "Authorization failed")
	RpcConnectErr          = NewErrNo(RpcConnectErrCode, "Rpc Connect failed")
	UploadVideoErr         = NewErrNo(UploadVideoErrCode, "Upload Video failed")
	ErrDatabase            = NewErrNo(ErrDatabaseCode, "Database error")
)

func Dprintf(format string, args ...interface{}) {
	if debug {
		log.Printf(format, args...)
	}
	return
}

type ErrNo struct {
	ErrCode int64
	ErrMsg  string
}

func (e ErrNo) Error() string {
	return fmt.Sprintf("err_code=%d, err_msg=%s", e.ErrCode, e.ErrMsg)
}

func NewErrNo(code int64, msg string) ErrNo {
	return ErrNo{code, msg}
}

func (e ErrNo) WithMessage(msg string) ErrNo {
	e.ErrMsg = msg
	return e
}

// ConvertErr convert error to Errno
func ConvertErr(err error) ErrNo {
	Err := ErrNo{}
	if errors.As(err, &Err) {
		return Err
	}

	s := ServiceErr
	s.ErrMsg = err.Error()
	return s
}
