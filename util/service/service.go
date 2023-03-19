/*
 * @Author: fuRan NgeKaworu@gmail.com
 * @Date: 2023-03-19 19:52:43
 * @LastEditors: fuRan NgeKaworu@gmail.com
 * @LastEditTime: 2023-03-19 23:05:51
 * @FilePath: /honghuang/util/service/service.go
 * @Description:
 *
 * Copyright (c) 2023 by ${git_name_email}, All Rights Reserved.
 */
package service

import (
	"github.com/NgeKaworu/util/db"
	"github.com/NgeKaworu/util/tool"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
)

type Service struct {
	Validate *validator.Validate
	Trans    *ut.Translator
	Mongo    *db.MongoClient
	Rdb      *redis.Client
	UCHost   *string
}

// New 工厂方法
func New(
	uc *string,
	r *string,
) *Service {

	mongo := db.NewMongoClient()
	validate := tool.NewValidator()
	trans := tool.NewValidatorTranslator(validate)
	rdb := redis.NewClient(&redis.Options{
		Addr:     *r,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return &Service{
		validate,
		trans,
		mongo,
		rdb,
		uc,
	}
}

func (srv *Service) Destroy() {
	srv.Mongo.Close()
	srv.Rdb.Close()
}
