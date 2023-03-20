/*
 * @Author: fuRan NgeKaworu@gmail.com
 * @Date: 2023-03-20 10:04:31
 * @LastEditors: fuRan NgeKaworu@gmail.com
 * @LastEditTime: 2023-03-20 10:13:47
 * @FilePath: /honghuang/util/tool/jwt.go
 * @Description:
 *
 * Copyright (c) 2023 by ${git_name_email}, All Rights Reserved.
 */
package tool

import (
	"errors"

	"github.com/dgrijalva/jwt-go"
)

// GenJWT generate jwt
func GenJWT(claims *jwt.StandardClaims, key []byte) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(key)
}

func CheckTokenAudience(auth string, key []byte) (audience *string, err error) {
	if auth == "" {
		err = errors.New("auth is empty")
		return
	}

	token, err := jwt.ParseWithClaims(auth, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})

	if err != nil {
		return
	}

	if tk, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
		audience = &tk.Audience
	}

	return
}
