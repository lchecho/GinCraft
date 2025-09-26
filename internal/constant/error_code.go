package constant

// 错误码定义
const (
	// 成功
	Success = 0

	// 系统级错误码
	SystemError     = 10001
	ParamError      = 10002
	DBError         = 10003
	Unauthorized    = 10004
	Forbidden       = 10005
	NotFound        = 10006
	MethodNotAllow  = 10007
	TooManyRequests = 10008
	Timeout         = 10009

	// 业务级错误码 (2xxxx)
	// 用户相关错误码 (200xx)
	UserNotExist         = 20001
	PasswordError        = 20002
	TokenExpired         = 20003
	TokenInvalid         = 20004
	UserAlreadyExist     = 20005
	UsernameAlreadyExist = 20006
	EmailAlreadyExist    = 20007
	UserCreateFailed     = 20008
	UserUpdateFailed     = 20009
	UserDeleteFailed     = 20010

	// 数据库连接错误码 (202xx)
	DBConnectionFailed  = 20201
	DBTransactionFailed = 20202
)

// ErrorMsg 错误码对应的错误信息
var ErrorMsg = map[int]string{
	Success: "成功",

	SystemError:     "系统错误",
	ParamError:      "参数错误",
	DBError:         "数据库错误",
	Unauthorized:    "未授权",
	Forbidden:       "禁止访问",
	NotFound:        "资源不存在",
	MethodNotAllow:  "方法不允许",
	TooManyRequests: "请求过多",
	Timeout:         "请求超时",

	// 用户相关错误信息
	UserNotExist:         "用户不存在",
	PasswordError:        "密码错误",
	TokenExpired:         "令牌已过期",
	TokenInvalid:         "无效的令牌",
	UserAlreadyExist:     "用户已存在",
	UsernameAlreadyExist: "用户名已存在",
	EmailAlreadyExist:    "邮箱已存在",
	UserCreateFailed:     "用户创建失败",
	UserUpdateFailed:     "用户更新失败",
	UserDeleteFailed:     "用户删除失败",

	// 数据库连接错误信息
	DBConnectionFailed:  "数据库连接失败",
	DBTransactionFailed: "数据库事务失败",
}

// GetMsg 获取错误信息
func GetMsg(code int) string {
	msg, ok := ErrorMsg[code]
	if ok {
		return msg
	}
	return ErrorMsg[SystemError]
}
