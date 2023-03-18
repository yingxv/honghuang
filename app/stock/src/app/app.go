/*
 * @Author: fuRan NgeKaworu@gmail.com
 * @Date: 2021-12-27 15:45:42
 * @LastEditors: fuRan NgeKaworu@gmail.com
 * @LastEditTime: 2023-02-13 15:48:54
 * @FilePath: /stock/stock-go/src/app/app.go
 * @Description:
 *
 * Copyright (c) 2023 by ${git_name_email}, All Rights Reserved.
 */
package app

import (
	"github.com/NgeKaworu/stock/src/db"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
)

// App
type App struct {
	validate *validator.Validate
	trans    *ut.Translator
	uc       *string
	mongo    *db.MongoClient
	rdb      *redis.Client
}

// New 工厂方法
func New(
	validate *validator.Validate,
	trans *ut.Translator,
	uc *string,
	mongo *db.MongoClient,
	rdb *redis.Client,
) *App {

	return &App{
		validate,
		trans,
		uc,
		mongo,
		rdb,
	}
}
