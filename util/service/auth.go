/*
 * @Author: fuRan NgeKaworu@gmail.com
 * @Date: 2023-03-19 03:04:50
 * @LastEditors: fuRan NgeKaworu@gmail.com
 * @LastEditTime: 2023-03-19 23:05:58
 * @FilePath: /honghuang/util/service/auth.go
 * @Description:
 *
 * Copyright (c) 2023 by ${git_name_email}, All Rights Reserved.
 */
package service

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/yingxv/honghuang/util/tool"
)

func (srv *Service) IsLogin(next http.Handler) http.Handler {
	//权限验证
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s, err := srv.checkUser(r)

		if err != nil {
			w.Header().Set("WWW-Authenticate", "Bearer realm=Restricted")
			w.WriteHeader(http.StatusUnauthorized)
			tool.RetFail(w, err)
			return
		}

		r.Header.Set("uid", *s)
		next.ServeHTTP(w, r)

	})
}

// checkUser
func (srv *Service) checkUser(r *http.Request) (*string, error) {
	bear, err := srv.getBearer(r)
	if err != nil {
		return nil, err
	}

	s, err := srv.Rdb.Get(context.Background(), *bear).Result()

	if err != nil {
		return nil, err
	}

	return &s, nil
}

// getBearer
func (srv *Service) getBearer(r *http.Request) (*string, error) {
	auth := r.Header.Get("Authorization")
	if !strings.HasPrefix(auth, "Bearer ") {
		return nil, errors.New("unknown authorization type")
	}
	auth = auth[7:]
	return &auth, nil
}

// perm mid
func (srv *Service) CheckPerm(perm string) func(httprouter.Handle) httprouter.Handle {
	return func(next httprouter.Handle) httprouter.Handle {
		//权限验证
		return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
			u := r.Header.Get("uid")
			k := u + ":perm"

			const (
				NONEXIST = 0
				EXIST    = 1
			)

			e, _ := srv.Rdb.Exists(context.Background(), k).Result()

			if e == EXIST {
				if b, _ := srv.Rdb.SIsMember(context.Background(), k, perm).Result(); b {
					next(w, r, ps)
					return
				}
				w.WriteHeader(http.StatusForbidden)
				return
			}

			if e == NONEXIST {
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
}
