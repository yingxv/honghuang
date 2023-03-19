/*
 * @Author: fuRan NgeKaworu@gmail.com
 * @Date: 2023-03-19 03:04:32
 * @LastEditors: fuRan NgeKaworu@gmail.com
 * @LastEditTime: 2023-03-19 18:49:24
 * @FilePath: /honghuang/app/user-center/src/app/captcha_c.go
 * @Description:
 *
 * Copyright (c) 2023 by ${git_name_email}, All Rights Reserved.
 */
package app

import (
	"errors"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/NgeKaworu/user-center/src/model"
	"github.com/NgeKaworu/util/tool"
	"github.com/hetiansu5/urlquery"
	"github.com/julienschmidt/httprouter"
)

func (app *App) FetchCaptcha(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	if app.getSetSessionLocked(w, r) {
		tool.RetFail(w, errors.New("请求频繁，请60s后再试"))
		return
	}

	p := struct {
		Email *string `query:"email,omitempty" validate:"required,email"`
	}{}

	err := urlquery.Unmarshal([]byte(r.URL.RawQuery), &p)
	if err != nil {
		tool.RetFail(w, err)
		return
	}

	err = app.validate.Struct(&p)
	if err != nil {
		tool.RetFailWithTrans(w, err, app.trans)
		return
	}

	captcha := padStartZero(rand.Intn(10000))
	app.setRedisCaptcha(p.Email, &captcha)

	w.Header().Set("Cache-Control", "max-age="+strconv.FormatInt(int64(MAX_AGE), 10))
	go app.sendCaptcha(p.Email, &captcha)

	tool.RetOk(w, "验证码已经发送")

}

func (app *App) CheckCaptcha(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var p model.Captcha

	err := urlquery.Unmarshal([]byte(r.URL.RawQuery), &p)
	if err != nil {
		tool.RetFail(w, err)
		return
	}

	err = app.checkCaptcha(w, r, &p)

	if err != nil {
		tool.RetFail(w, errors.New("验证失败"))
		return
	}

	tool.RetOk(w, "验证通过")

}

func padStartZero(i int) string {
	s := strconv.FormatInt(int64(i), 10)
	l := 4 - len(s)
	for l > 0 {
		l--
		s = "0" + s
	}
	return s
}
