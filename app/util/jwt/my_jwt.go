package jwt

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"sim/app/global/variable"
	"sync"
)

var JwtSignIns *JwtSign
var JwtSignOnce sync.Once

func CreateJwtSign() *JwtSign {
	JwtSignOnce.Do(func() {
		JwtSignIns = &JwtSign{
			SignKey: []byte(variable.ConfigYml.GetString("jwt.signKey")),
		}
	})

	return JwtSignIns
}

type JwtSign struct {
	SignKey []byte
}

// 创建token
func (j *JwtSign) CreateToken(c Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)

	signToken, err := token.SignedString(j.SignKey)
	if err != nil {
		return "", err
	}
	return signToken, nil
}

// ParseToken 解析token
func (j *JwtSign) ParseToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SignKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
