package jwt

import (
	"github.com/dgrijalva/jwt-go"
	userPb "sim/idl/user"
)

type Claims struct {
	Info userPb.UserResponse
	jwt.StandardClaims
}
