/*
 * @Author: fuRan NgeKaworu@gmail.com
 * @Date: 2023-03-19 03:04:11
 * @LastEditors: fuRan NgeKaworu@gmail.com
 * @LastEditTime: 2023-03-19 22:33:52
 * @FilePath: /honghuang/app/flashcard/src/controller/controller.go
 * @Description:
 *
 * Copyright (c) 2023 by ${git_name_email}, All Rights Reserved.
 */
package controller

import "github.com/NgeKaworu/util/service"

// Controller 控制器
type Controller struct {
	srv *service.Service
}

// NewController 工厂方法
func NewController(srv *service.Service) *Controller {

	return &Controller{
		srv,
	}
}
