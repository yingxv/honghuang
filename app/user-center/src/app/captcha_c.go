package app

import (
	"errors"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/NgeKaworu/user-center/src/model"
	"github.com/NgeKaworu/user-center/src/util/responser"
	"github.com/hetiansu5/urlquery"
	"github.com/julienschmidt/httprouter"
)

func (app *App) FetchCaptcha(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	if app.getSetSessionLocked(w, r) {
		responser.RetFail(w, errors.New("请求频繁，请60s后再试"))
		return
	}

	p := struct {
		Email *string `query:"email,omitempty" validate:"required,email"`
	}{}

	err := urlquery.Unmarshal([]byte(r.URL.RawQuery), &p)
	if err != nil {
		responser.RetFail(w, err)
		return
	}

	err = app.validate.Struct(&p)
	if err != nil {
		responser.RetFailWithTrans(w, err, app.trans)
		return
	}

	captcha := padStartZero(rand.Intn(10000))
	app.setRedisCaptcha(p.Email, &captcha)

	w.Header().Set("Cache-Control", "max-age="+strconv.FormatInt(int64(MAX_AGE), 10))
	go app.sendCaptcha(p.Email, &captcha)

	responser.RetOk(w, "验证码已经发送")

}

func (app *App) CheckCaptcha(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var p model.Captcha

	err := urlquery.Unmarshal([]byte(r.URL.RawQuery), &p)
	if err != nil {
		responser.RetFail(w, err)
		return
	}

	err = app.checkCaptcha(w, r, &p)

	if err != nil {
		responser.RetFail(w, errors.New("验证失败"))
		return
	}

	responser.RetOk(w, "验证通过")

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
