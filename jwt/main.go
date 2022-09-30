package main

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type JWTCustomClaims struct {
	UserId int
	jwt.StandardClaims
}

func main() {
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/user", userHandler)
	http.ListenAndServe(":8080", nil)
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	token, err := authJWT(r)
	if err != nil {
		w.Write([]byte("权限不足"))
		return
	}
	id := "user是:" + strconv.Itoa(token.Claims.(*JWTCustomClaims).UserId)
	w.Write([]byte(id))
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	token, err := generateToken(1)
	if err != nil {
		w.Write([]byte("注册失败"))
		return
	}
	w.Write([]byte(token))
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	_, err := authJWT(r)
	if err != nil {
		w.Write([]byte("登录失败"))
		return
	}
	w.Write([]byte("登录成功"))
}

func authJWT(r *http.Request) (*jwt.Token, error) {
	tokenString, err := getJwtTokenFromHeader(r)
	if err != nil {
		return nil, err
	}

	token, err := parseToken(tokenString)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func parseToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	return token, err
}

func generateToken(userId int) (string, error) {
	claims := JWTCustomClaims{
		userId,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Unix() + 3600,
			Issuer:    "test",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("secret"))
}

func getJwtTokenFromHeader(r *http.Request) (string, error) {
	tokenString := r.Header.Get("Authorization")
	fmt.Println(tokenString)

	tokenResult := strings.SplitN(tokenString, " ", 2)
	if len(tokenResult) != 2 || tokenResult[0] != "Bearer" {
		return "", fmt.Errorf("token格式错误")
	}

	return tokenResult[1], nil
}
