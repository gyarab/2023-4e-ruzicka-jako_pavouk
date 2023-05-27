package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

var privatniKlic []byte = []byte(os.Getenv("KLIC"))

type Data struct {
	jwt.StandardClaims
	Email string
	Id    uint
}

func HashPassword(heslo string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(heslo), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func CheckPassword(hesloRequest string, hesloDB string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hesloDB), []byte(hesloRequest))
	if err != nil {
		return err
	}
	return nil
}

func GenerovatToken(email string, id uint) (string, error) {
	data := Data{Email: email, Id: id, StandardClaims: jwt.StandardClaims{ExpiresAt: time.Now().Add(time.Hour * 24 * 7).Unix() /*Týden asi good*/}}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, data)
	s, err := token.SignedString(privatniKlic)
	return s, err
}

func ValidovatToken(tokenString string) (bool, uint, error) {
	data := Data{}
	token, err := jwt.ParseWithClaims(tokenString, &data, func(token *jwt.Token) (interface{}, error) {
		return privatniKlic, nil
	})
	if err != nil {
		return false, 0, err
	}
	//fmt.Printf("%+v\n", token)
	return token.Valid, data.Id, err
}
