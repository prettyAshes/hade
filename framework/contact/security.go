package contact

import "github.com/dgrijalva/jwt-go"

// SecurityKey 代表 安全服务
const SecurityKey = "hade:security"

// SecurityService
type SecurityService interface {
	XSS()
	GenToken(originVal map[string]interface{}) string // CSRF生成token
	VerifyToken(token string) bool                    // CSRF验证token
}

type MyCustomClaims struct {
	jwt.StandardClaims
	Origin map[string]interface{}
}
