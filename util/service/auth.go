/*
 * @Author: fuRan NgeKaworu@gmail.com
 * @Date: 2023-03-19 03:04:50
 * @LastEditors: fuRan NgeKaworu@gmail.com
 * @LastEditTime: 2023-11-03 20:12:35
 * @FilePath: /honghuang/util/service/auth.go
 * @Description:
 *
 * Copyright (c) 2023 by ${git_name_email}, All Rights Reserved.
 */
package service

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (srv *Service) IsLogin(next http.Handler) http.Handler {
	//权限验证
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		client := &http.Client{}

		req, _ := http.NewRequest("HEAD", *srv.UCHost+"/rpc/is-login-rpc", nil)
		req.Header.Set("Authorization", r.Header.Get("Authorization"))
		res, err := client.Do(req)

		if err != nil || res.StatusCode != http.StatusOK {
			w.WriteHeader(res.StatusCode)
			return
		}

		r.Header.Set("uid", res.Header.Get("uid"))
		next.ServeHTTP(w, r)
	})
}

// perm mid
func (srv *Service) CheckPerm(perm string) func(httprouter.Handle) httprouter.Handle {
	return func(next httprouter.Handle) httprouter.Handle {
		//权限验证
		return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
			client := &http.Client{}

			req, _ := http.NewRequest("HEAD", *srv.UCHost+"/check-perm-rpc/"+perm, nil)
			req.Header.Set("Authorization", r.Header.Get("Authorization"))
			res, err := client.Do(req)

			if err != nil || res.StatusCode != http.StatusOK {
				w.WriteHeader(res.StatusCode)
				return
			}

			next(w, r, ps)
		}
	}
}
