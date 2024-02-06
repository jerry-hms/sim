package user

type BaseField struct {
	Username string `json:"username" form:"username" binding:"required" err_msg:"请输入用户名"`
	Password string `json:"password" form:"password" binding:"required" err_msg:"请输入密码"`
}
