package app

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/hetiansu5/urlquery"
	"github.com/julienschmidt/httprouter"
	"github.com/yingxv/honghuang/user-center/src/model"
	"github.com/yingxv/honghuang/util/tool"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// RoleCreate 新增角色
func (app *App) RoleCreate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		tool.RetFail(w, err)
		return
	}

	if len(body) == 0 {
		tool.RetFail(w, errors.New("not has body"))
		return
	}

	var u model.Role
	err = json.Unmarshal(body, &u)
	if err != nil {
		tool.RetFail(w, err)
		return
	}

	if err := app.srv.Validate.Struct(u); err != nil {
		tool.RetFail(w, err)
		return
	}
	time := time.Now().Local()
	u.CreateAt = &time

	t := app.srv.Mongo.GetColl(model.TRole)

	res, err := t.InsertOne(context.Background(), u)

	if err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "dup key") {
			errMsg = "该key已经使用"
		}

		tool.RetFail(w, errors.New(errMsg))
		return

	}

	if err != nil {
		tool.RetFail(w, err)
		return
	}

	tool.RetOk(w, res.InsertedID)
}

// RoleRemove 删除角色
func (app *App) RoleRemove(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	if id == "" {
		tool.RetFail(w, errors.New("ID不能为空"))
		return
	}
	c, err := app.srv.Mongo.GetColl(model.TUser).CountDocuments(context.Background(), bson.M{
		"roles": id,
	})

	if err != nil || c > 0 {
		tool.RetFail(w, errors.New("无法删除使用中角色"))
		return
	}

	res := app.srv.Mongo.GetColl(model.TRole).FindOneAndDelete(context.Background(), bson.M{
		"_id": id,
	})

	if res.Err() != nil {
		tool.RetFail(w, res.Err())
		return
	}

	tool.RetOk(w, "删除成功")
}

// RoleUpdate 修改角色
func (app *App) RoleUpdate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		tool.RetFail(w, err)
		return
	}

	if len(body) == 0 {
		tool.RetFail(w, errors.New("not has body"))
	}

	var u model.Role

	err = json.Unmarshal(body, &u)

	if err != nil {
		tool.RetFail(w, err)
		return
	}

	if u.ID == nil {
		tool.RetFail(w, errors.New("id不能为空"))
		return
	}

	if err := app.srv.Validate.Struct(u); err != nil {
		tool.RetFail(w, err)
		return
	}

	updateAt := time.Now().Local()
	u.UpdateAt = &updateAt
	updater := bson.M{"$set": &u}

	if u.Perms == nil {
		updater["$unset"] = bson.M{"perms": ""}
	}

	res := app.srv.Mongo.GetColl(model.TRole).FindOneAndUpdate(context.Background(), bson.M{"_id": *u.ID}, updater)

	if res.Err() != nil {
		errMsg := res.Err().Error()
		if strings.Contains(errMsg, "dup key") {
			errMsg = "该key已经使用"
		}

		tool.RetFail(w, errors.New(errMsg))
		return
	}

	tool.RetOk(w, "操作成功")
}

// RoleList 查找角色
func (app *App) RoleList(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	p := struct {
		Keyword *string `query:"keyword,omitempty" validate:"omitempty"`
		Skip    *int64  `query:"skip,omitempty" validate:"omitempty,min=0"`
		Limit   *int64  `query:"limit,omitempty" validate:"omitempty,min=0"`
	}{}

	err := urlquery.Unmarshal([]byte(r.URL.RawQuery), &p)
	if err != nil {
		tool.RetFail(w, err)
		return
	}

	err = app.srv.Validate.Struct(&p)
	if err != nil {
		tool.RetFailWithTrans(w, err, app.srv.Trans)
		return
	}

	params := bson.M{}

	if p.Keyword != nil {
		params = bson.M{
			"$or": []bson.M{
				{"name": bson.M{"$regex": p.Keyword}},
				{"_id": bson.M{"$regex": p.Keyword}},
			},
		}
	}

	opt := options.Find()

	if p.Limit != nil {
		opt.SetLimit(*p.Limit)
	} else {
		opt.SetLimit(10)
	}

	if p.Skip != nil {
		opt.SetSkip(*p.Skip)
	}
	t := app.srv.Mongo.GetColl(model.TRole)

	cur, err := t.Find(context.Background(), params, opt)

	if err != nil {
		tool.RetFail(w, err)
		return
	}

	var roles []model.Role
	err = cur.All(context.Background(), &roles)

	if err != nil {
		tool.RetFail(w, err)
		return
	}

	total, err := t.CountDocuments(context.Background(), params)

	if err != nil {
		tool.RetFail(w, err)
		return
	}

	tool.RetOkWithTotal(w, roles, total)
}

// RoleValidateKey key 校验
func (app *App) RoleValidateKey(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	p := struct {
		ID *string `query:"id,omitempty" validate:"omitempty,required"`
	}{}

	err := urlquery.Unmarshal([]byte(r.URL.RawQuery), &p)
	if err != nil {
		tool.RetFail(w, err)
		return
	}

	err = app.srv.Validate.Struct(&p)
	if err != nil {
		tool.RetFailWithTrans(w, err, app.srv.Trans)
		return
	}

	total, err := app.srv.Mongo.GetColl(model.TRole).CountDocuments(context.Background(), bson.M{"_id": p.ID})

	if err != nil {
		tool.RetFail(w, err)
		return
	}

	if total != 0 {
		tool.RetFail(w, errors.New("key 重复"))
		return
	}

	tool.RetOk(w, "validate key")
}
