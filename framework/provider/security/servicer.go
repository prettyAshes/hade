package security

import (
	"errors"
	"hade/framework"
	"hade/framework/contact"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type HadeSecurityServicer struct {
	c framework.Container

	duration time.Duration
	slat     string
}

func NewSecurityServicer(...interface{}) (interface{}, error) {
	return nil, nil
}

func (servicer *HadeSecurityServicer) XSS() {

}

// genToken 生成Token
func (servicer *HadeSecurityServicer) GenToken(origin map[string]interface{}) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(servicer.duration)
	claims := contact.MyCustomClaims{
		Origin: origin,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
		},
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(servicer.slat))
}

// VerifyToken 验证token
func (servicer *HadeSecurityServicer) VerifyToken(token string) (map[string]interface{}, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &contact.MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(servicer.slat), nil
	})
	if err != nil {
		return nil, err
	}

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*contact.MyCustomClaims); ok && tokenClaims.Valid {
			return claims.Origin, nil
		}
	}

	return nil, errors.New("verify token error")
}
