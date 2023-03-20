/*
 * @Author: fuRan NgeKaworu@gmail.com
 * @Date: 2023-03-19 19:52:43
 * @LastEditors: fuRan NgeKaworu@gmail.com
 * @LastEditTime: 2023-03-20 10:26:59
 * @FilePath: /honghuang/util/service/service.go
 * @Description:
 *
 * Copyright (c) 2023 by ${git_name_email}, All Rights Reserved.
 */
package service

import (
	"github.com/NgeKaworu/util/db"
	"github.com/NgeKaworu/util/service/key"
	"github.com/NgeKaworu/util/tool"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"github.com/robfig/cron/v3"
	"gopkg.in/gomail.v2"
)

type Service struct {
	Validate *validator.Validate
	Trans    *ut.Translator
	Mongo    *db.MongoClient
	Rdb      *redis.Client
	UCHost   *string
	Key      *key.Key
	Dialer   *gomail.Dialer
	Cron     *cron.Cron
}

type ServiceAug struct {
	UCHost    *string
	RedisAddr *string
	CipherKey *string
	DialerPwd *string
}

// New 工厂方法
func New(s *ServiceAug) *Service {
	// 校验器
	validate := tool.NewValidator()

	srv := &Service{
		Validate: validate,
		// 校验器翻译
		Trans: tool.NewValidatorTranslator(validate),
		// Mongo客户端
		Mongo: db.NewMongoClient(),
		// 用户中心地址
		UCHost: s.UCHost,
		// Cron定时任务
		Cron: cron.New(cron.WithParser(cron.NewParser(cron.Hour | cron.Dow))),
	}

	// redis 客户端
	if nil != s.RedisAddr {
		srv.Rdb = redis.NewClient(&redis.Options{
			Addr:     *s.RedisAddr,
			Password: "", // no password set
			DB:       0,  // use default DB
		})
	}

	// 加密器
	if nil != s.CipherKey {
		srv.Key = key.New(*s.CipherKey)
	}

	// smtp 客户端
	if nil != s.DialerPwd {
		srv.Dialer = gomail.NewDialer("smtp.gmail.com", 587, "ngekaworu@gmail.com", *s.DialerPwd)
	}

	return srv
}

func (srv *Service) Destroy() {
	if nil != srv.Mongo {
		srv.Mongo.Close()
	}

	if nil != srv.Rdb {
		srv.Rdb.Close()
	}

	if nil != srv.Cron {
		srv.Cron.Stop()
	}
}
