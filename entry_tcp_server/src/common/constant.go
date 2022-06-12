package constant

var (
	//jwt密钥
	JwtKey         = []byte("entryTask")
	ConfigFileName = "config_local.ini"
)

const (
	Success                   = 200
	InvalidParam              = 400
	InvalidUsernameOrPassword = 401
	UserExist                 = 402
	UpdateFailed              = 403
	SystemError               = 404
	InvalidData               = 405
	UnAuthorized              = 406
)
