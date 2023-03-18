package app

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/NgeKaworu/stock/src/util"
	"github.com/julienschmidt/httprouter"
)

func (app *App) IsLogin(next http.Handler) http.Handler {
	//权限验证
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s, err := app.checkUser(r)

		if err != nil {
			w.Header().Set("WWW-Authenticate", "Bearer realm=Restricted")
			w.WriteHeader(http.StatusUnauthorized)
			util.RetFail(w, err)
			return
		}

		r.Header.Set("uid", *s)
		next.ServeHTTP(w, r)

	})
}

// checkUser
func (app *App) checkUser(r *http.Request) (*string, error) {
	bear, err := app.getBearer(r)
	if err != nil {
		return nil, err
	}

	s, err := app.rdb.Get(context.Background(), *bear).Result()

	if err != nil {
		return nil, err
	}

	return &s, nil
}

// getBearer
func (app *App) getBearer(r *http.Request) (*string, error) {
	auth := r.Header.Get("Authorization")
	if !strings.HasPrefix(auth, "Bearer ") {
		return nil, errors.New("unknown authorization type")
	}
	auth = auth[7:]
	return &auth, nil
}

// perm mid
func (app *App) CheckPerm(perm string) func(httprouter.Handle) httprouter.Handle {
	return func(next httprouter.Handle) httprouter.Handle {
		//权限验证
		return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
			u := r.Header.Get("uid")
			k := u + ":perm"

			const (
				NONEXIST = 0
				EXIST    = 1
			)

			e, _ := app.rdb.Exists(context.Background(), k).Result()

			if e == EXIST {
				if b, _ := app.rdb.SIsMember(context.Background(), k, perm).Result(); b {
					next(w, r, ps)
					return
				}
				w.WriteHeader(http.StatusForbidden)
				return
			}

			if e == NONEXIST {
				client := &http.Client{}
				req, _ := http.NewRequest("HEAD", *app.uc+"/check-perm-rpc/"+perm, nil)
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
