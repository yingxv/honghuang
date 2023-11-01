/*
 * @Author: fuRan NgeKaworu@gmail.com
 * @Date: 2023-03-20 09:58:42
 * @LastEditors: fuRan NgeKaworu@gmail.com
 * @LastEditTime: 2023-11-01 13:54:30
 * @FilePath: /honghuang/util/service/key/key.go
 * @Description:
 *
 * Copyright (c) 2023 by ${git_name_email}, All Rights Reserved.
 */
package key

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/yingxv/honghuang/util/tool"
)

// Auth 加解密相关
type Key struct {
	key []byte
}

// New 工厂方法
func New(k string) *Key {
	return &Key{
		key: []byte(k),
	}
}

func (k *Key) GenJWT(claims *jwt.StandardClaims) (string, error) {
	return tool.GenJWT(claims, k.key)
}
func (k *Key) CheckTokenAudience(auth string) (audience *string, err error) {
	return tool.CheckTokenAudience(auth, k.key)
}
func (k *Key) CFBEncrypter(s string) ([]byte, error) {
	return tool.CFBEncrypter(s, k.key)
}
func (k *Key) CFBDecrypter(s string) ([]byte, error) {
	return tool.CFBDecrypter(s, k.key)
}
