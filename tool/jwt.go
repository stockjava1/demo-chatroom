package tool

import (
	"errors"
	"time"

	"fmt"
	t "github.com/dgrijalva/jwt-go"
	"github.com/iris-contrib/middleware/jwt"
)

// 定义一个密钥，用于签名和验证 JWT
var secretKey = []byte("my_secret_key")

// GetJWTString get jwt string with expiration time 20 minutes
func GetJWTString(name string, id int64) (string, error) {
	token := jwt.NewTokenWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		// 根据需求，可以存一些必要的数据
		"userName": name,
		"userId":   id,

		// 签发人
		"iss": "iris",
		// 签发时间
		"iat": time.Now().Unix(),
		// 设定过期时间，设置20分钟过期
		"exp": time.Now().Add(24 * time.Hour * time.Duration(1)).Unix(),
	})

	// 使用设置的秘钥，签名生成jwt字符串
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseToken(tokenString string) (t.MapClaims, error) {

	token, err := t.Parse(tokenString, func(token *t.Token) (interface{}, error) {
		// 校验token的签名
		if _, ok := token.Method.(*t.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		// 返回签名的密钥
		return []byte(secretKey), nil
	})

	if err != nil {
		// 解析出错
		fmt.Println("JWT parsing error:", err)
		return nil, err
	}
	fmt.Println(fmt.Sprintf("User ID >>>>>>>>> : %v", token))
	if claims, ok := token.Claims.(t.MapClaims); ok && token.Valid {
		// 解析出token中的数据
		fmt.Println("User ID:", claims["userId"])
		fmt.Println("Username:", claims["userName"])

		return claims, nil
	} else {
		fmt.Println("Invalid JWT token")
		return nil, errors.New("Invalid JWT token")
	}

}
