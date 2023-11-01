/*
 * @Author: fuRan NgeKaworu@gmail.com
 * @Date: 2021-12-27 15:45:42
 * @LastEditors: fuRan NgeKaworu@gmail.com
 * @LastEditTime: 2023-03-19 22:49:53
 * @FilePath: /honghuang/app/stock/src/app/app.go
 * @Description:
 *
 * Copyright (c) 2023 by ${git_name_email}, All Rights Reserved.
 */
package app

import (
	"github.com/yingxv/honghuang/util/service"
)

// App
type App struct {
	srv *service.Service
}

// New 工厂方法
func New(
	srv *service.Service,
) *App {

	return &App{
		srv,
	}
}
