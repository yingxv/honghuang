/*
 * @Author: fuRan NgeKaworu@gmail.com
 * @Date: 2023-03-19 12:29:13
 * @LastEditors: fuRan NgeKaworu@gmail.com
 * @LastEditTime: 2023-03-19 12:37:02
 * @FilePath: /honghuang/util/tool/required.go
 * @Description:
 *
 * Copyright (c) 2023 by ${git_name_email}, All Rights Reserved.
 */
package tool

import "errors"

// Required 判断是否为空
func Required(m map[string]interface{}, r map[string]string) error {
	var s string
	for k, v := range r {
		if _, ok := m[k]; !ok {
			s += v + " "
		}
	}
	if s != "" {
		return errors.New(s)
	}
	return nil
}
