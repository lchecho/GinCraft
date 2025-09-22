package constant

// 错误码定义
const (
	// 成功
	SUCCESS = 0

	// 系统级错误码
	SYSTEM_ERROR      = 10001
	PARAM_ERROR       = 10002
	DB_ERROR          = 10003
	UNAUTHORIZED      = 10004
	FORBIDDEN         = 10005
	NOT_FOUND         = 10006
	METHOD_NOT_ALLOW  = 10007
	TOO_MANY_REQUESTS = 10008
	TIMEOUT           = 10009

	// 业务级错误码 (2xxxx)
	// 用户相关错误码 (200xx)
	USER_NOT_EXIST         = 20001
	PASSWORD_ERROR         = 20002
	TOKEN_EXPIRED          = 20003
	TOKEN_INVALID          = 20004
	USER_ALREADY_EXIST     = 20005
	USERNAME_ALREADY_EXIST = 20006
	EMAIL_ALREADY_EXIST    = 20007
	USER_CREATE_FAILED     = 20008
	USER_UPDATE_FAILED     = 20009
	USER_DELETE_FAILED     = 20010

	// 数据库连接错误码 (202xx)
	DB_CONNECTION_FAILED  = 20201
	DB_TRANSACTION_FAILED = 20202
)

// ErrorMsg 错误码对应的错误信息
var ErrorMsg = map[int]string{
	SUCCESS: "成功",

	SYSTEM_ERROR:      "系统错误",
	PARAM_ERROR:       "参数错误",
	DB_ERROR:          "数据库错误",
	UNAUTHORIZED:      "未授权",
	FORBIDDEN:         "禁止访问",
	NOT_FOUND:         "资源不存在",
	METHOD_NOT_ALLOW:  "方法不允许",
	TOO_MANY_REQUESTS: "请求过多",
	TIMEOUT:           "请求超时",

	// 用户相关错误信息
	USER_NOT_EXIST:         "用户不存在",
	PASSWORD_ERROR:         "密码错误",
	TOKEN_EXPIRED:          "令牌已过期",
	TOKEN_INVALID:          "无效的令牌",
	USER_ALREADY_EXIST:     "用户已存在",
	USERNAME_ALREADY_EXIST: "用户名已存在",
	EMAIL_ALREADY_EXIST:    "邮箱已存在",
	USER_CREATE_FAILED:     "用户创建失败",
	USER_UPDATE_FAILED:     "用户更新失败",
	USER_DELETE_FAILED:     "用户删除失败",

	// 数据库连接错误信息
	DB_CONNECTION_FAILED:  "数据库连接失败",
	DB_TRANSACTION_FAILED: "数据库事务失败",
}

// GetMsg 获取错误信息
func GetMsg(code int) string {
	msg, ok := ErrorMsg[code]
	if ok {
		return msg
	}
	return ErrorMsg[SYSTEM_ERROR]
}
