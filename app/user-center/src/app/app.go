/*
 * @Author: fuRan NgeKaworu@gmail.com
 * @Date: 2023-03-19 03:04:23
 * @LastEditors: fuRan NgeKaworu@gmail.com
 * @LastEditTime: 2023-03-20 10:32:42
 * @FilePath: /honghuang/app/user-center/src/app/app.go
 * @Description:
 *
 * Copyright (c) 2023 by ${git_name_email}, All Rights Reserved.
 */
package app

import "github.com/NgeKaworu/util/service"

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
