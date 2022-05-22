package Jwt

import (
	"github.com/dgrijalva/jwt-go"
	"golandproject/Class"
	"time"
)

func GenerateToken(uid string, ip string) (string, error) {
	expireTime := time.Now().Add(time.Second * 60 * 60 * 24 * 7) //登录有效期为7天
	claims := Class.CustomClaims{
		Uid: uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Audience:  ip,
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte("MTY1MjgxMDg2MHxOd3dBTkVKRFFrdEZUa0ZJTkVwVVFWcEdSMHhCUVRNeVZrMUhSbEpMUkVoUU0wZEdVMUJCVWpKYVNrd3lXVnBKUlVkRlF6TlBURkU9fPlbyxwil3sCL6pwYb_U6xI0PgydY-wGXL5_W06841Gd"))
	return token, err
}

func TokenValid(token string, ip string) bool {
	tokenClaims, err := jwt.ParseWithClaims(token, &Class.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("MTY1MjgxMDg2MHxOd3dBTkVKRFFrdEZUa0ZJTkVwVVFWcEdSMHhCUVRNeVZrMUhSbEpMUkVoUU0wZEdVMUJCVWpKYVNrd3lXVnBKUlVkRlF6TlBURkU9fPlbyxwil3sCL6pwYb_U6xI0PgydY-wGXL5_W06841Gd"), nil
	})
	if err != nil {
		return false
	}
	if tokenClaims != nil {
		claims, ok := tokenClaims.Claims.(*Class.CustomClaims)
		//fmt.Printf("%#v %#v %#v %#v", ok, tokenClaims.Valid, claims, ip)
		if ok && tokenClaims.Valid && claims.StandardClaims.Audience == ip {
			return true
		}
	}
	return false
}

func TokenGetUid(token string) (uid string) {
	tokenClaims, _ := jwt.ParseWithClaims(token, &Class.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("MTY1MjgxMDg2MHxOd3dBTkVKRFFrdEZUa0ZJTkVwVVFWcEdSMHhCUVRNeVZrMUhSbEpMUkVoUU0wZEdVMUJCVWpKYVNrd3lXVnBKUlVkRlF6TlBURkU9fPlbyxwil3sCL6pwYb_U6xI0PgydY-wGXL5_W06841Gd"), nil
	})
	claims, _ := tokenClaims.Claims.(*Class.CustomClaims)
	return claims.Uid
}

func TokenGetIp(token string) (uid string) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Class.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("MTY1MjgxMDg2MHxOd3dBTkVKRFFrdEZUa0ZJTkVwVVFWcEdSMHhCUVRNeVZrMUhSbEpMUkVoUU0wZEdVMUJCVWpKYVNrd3lXVnBKUlVkRlF6TlBURkU9fPlbyxwil3sCL6pwYb_U6xI0PgydY-wGXL5_W06841Gd"), nil
	})
	if err != nil {
		return "0.0.0.0"
	}
	claims, _ := tokenClaims.Claims.(*Class.CustomClaims)
	return claims.StandardClaims.Audience
}

func ParseToken(token string) (*Class.CustomClaims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Class.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("MTY1MjgxMDg2MHxOd3dBTkVKRFFrdEZUa0ZJTkVwVVFWcEdSMHhCUVRNeVZrMUhSbEpMUkVoUU0wZEdVMUJCVWpKYVNrd3lXVnBKUlVkRlF6TlBURkU9fPlbyxwil3sCL6pwYb_U6xI0PgydY-wGXL5_W06841Gd"), nil
	})
	if err != nil {
		return nil, err
	}

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Class.CustomClaims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}
