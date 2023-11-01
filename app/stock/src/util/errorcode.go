/*
 * @Author: fuRan NgeKaworu@gmail.com
 * @Date: 2021-08-30 13:24:40
 * @LastEditors: fuRan NgeKaworu@gmail.com
 * @LastEditTime: 2023-03-19 18:43:23
 * @FilePath: /honghuang/app/stock/src/util/errorcode.go
 * @Description:
 *
 * Copyright (c) 2023 by ${git_name_email}, All Rights Reserved.
 */
package util

import "github.com/yingxv/honghuang/util/tool"

const (
	YEAR_ERR tool.Bits = 1 << iota
	CUR_ERR
)
