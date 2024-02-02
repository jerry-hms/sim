package jwt

import (
	"github.com/dgrijalva/jwt-go"
	userPb "sim/idl/pb/user"
)

type Claims struct {
	Info userPb.UserResponse
	jwt.StandardClaims
}
