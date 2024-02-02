package token

import (
	"errors"
	"sim/app/global/consts"
	"sim/app/util/jwt"
	userPb "sim/idl/pb/user"
	"sync"
	"time"
)

var UserTokenIns *UserToken
var UserTokenOnce sync.Once

type UserToken struct {
	JwtToken jwt.JwtSign
}

// 创建UserToken工厂
func CreateUserTokenFactory() *UserToken {
	UserTokenOnce.Do(func() {
		UserTokenIns = &UserToken{}
	})
	return UserTokenIns
}

func (u *UserToken) GenerateToken(user interface{}, expiredAt int64) (string, error) {
	var claims jwt.Claims
	userInfo, ok := user.(*userPb.UserResponse)
	if !ok {
		return "", errors.New("user转换UserResponse失败")
	}
	claims.Info = *userInfo
	//jsonStr, _ := json.Marshal(user)
	//_ = json.Unmarshal(jsonStr, &claims)
	claims.StandardClaims.NotBefore = time.Now().Unix() - 10        // 生效开始时间 给个10秒的浮动区间
	claims.StandardClaims.ExpiresAt = time.Now().Unix() + expiredAt // 有效时间

	return jwt.CreateJwtSign().CreateToken(claims)
}

func (u *UserToken) IsNotExpired(tokenStr string) (*jwt.Claims, int) {
	claims, err := jwt.CreateJwtSign().ParseToken(tokenStr)
	if err == nil {
		// 校验是否在有效期内
		if time.Now().Unix()-claims.ExpiresAt < 0 {
			return claims, 1
		} else {
			return nil, consts.TokenExpiredError
		}
	}
	return nil, consts.TokenInvalidError
}

// IsEffective 判断token是否有效
func (u *UserToken) IsEffective(tokenStr string) (*jwt.Claims, int) {
	return u.IsNotExpired(tokenStr)
}
