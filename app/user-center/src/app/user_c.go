package app

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/hetiansu5/urlquery"
	"github.com/julienschmidt/httprouter"
	"github.com/yingxv/honghuang/user-center/src/model"
	"github.com/yingxv/honghuang/util/tool"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Login 登录
func (app *App) Login(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	type user struct {
		ID    *primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty" `                               // id
		Pwd   *string             `json:"pwd,omitempty" bson:"pwd,omitempty" validate:"required"`           // 账号
		Email *string             `json:"email,omitempty" bson:"email,omitempty" validate:"email,required"` // 密码
	}

	inputUser := new(user)

	if err := json.Unmarshal(body, &inputUser); err != nil {
		tool.RetFail(w, err)
		return
	}

	if err := app.srv.Validate.Struct(inputUser); err != nil {
		tool.RetFailWithTrans(w, err, app.srv.Trans)
		return
	}

	t := app.srv.Mongo.GetColl(model.TUser)

	email := strings.ToLower(strings.Replace(*inputUser.Email, " ", "", -1))
	res := t.FindOne(context.Background(), bson.M{
		"email": email,
	})

	if res.Err() != nil {
		tool.RetFail(w, errors.New("用户名或密码不正确"))
		return
	}

	outputUser := new(user)

	err = res.Decode(&outputUser)
	if err != nil {
		tool.RetFail(w, err)
		return
	}

	dec, err := app.srv.Key.CFBDecrypter(*outputUser.Pwd)
	if err != nil {
		tool.RetFail(w, err)
		return
	}

	if string(dec) != *inputUser.Pwd {
		tool.RetFail(w, errors.New("用户名或密码不正确"))
		return
	}

	uid := outputUser.ID.Hex()

	_, err = app.srv.Rdb.Del(context.Background(), uid+":perm").Result()
	if err != nil {
		tool.RetFail(w, err)
		return
	}

	app.cacheSign(w, uid)
}

// Regsiter 注册
func (app *App) Regsiter(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	var u model.User

	if err := json.Unmarshal(body, &u); err != nil {
		tool.RetFail(w, err)
		return
	}

	err = app.checkCaptcha(w, r, u.ToCaptcha())

	if err != nil {
		tool.RetFail(w, err)
		return
	}

	u.Roles = []string{"user"}

	res, err := app.insertOneUser(&u)
	if err != nil {
		tool.RetFail(w, err)
		return
	}

	err = app.removeRedisCaptcha(u.Email)

	if err != nil {
		log.Println(err)
	}

	app.cacheSign(w, res.InsertedID.(primitive.ObjectID).Hex())

}

// ForgetPwd 忘记密码
func (app *App) ForgetPwd(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		tool.RetFail(w, err)
		return
	}
	if len(body) == 0 {
		tool.RetFail(w, errors.New("not has body"))
	}

	var u struct {
		ID       *primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`                                // id
		Pwd      *string             `json:"pwd,omitempty" bson:"pwd,omitempty" validate:"required,min=8"`     // 密码
		Email    *string             `json:"email,omitempty" bson:"email,omitempty" validate:"required,email"` // 邮箱
		Captcha  *string             `json:"captcha,omitempty" bson:"-"`                                       // 验证码
		UpdateAt *time.Time          `json:"updateAt,omitempty" bson:"updateAt,omitempty"`                     // 更新时间
	}

	err = json.Unmarshal(body, &u)
	if err != nil {
		tool.RetFail(w, err)
		return
	}

	capcha := model.Captcha{
		Email:   u.Email,
		Captcha: u.Captcha,
	}

	err = app.checkCaptcha(w, r, &capcha)
	if err != nil {
		tool.RetFail(w, err)
		return
	}

	if err := app.srv.Validate.Struct(u); err != nil {
		tool.RetFail(w, err)
		return
	}

	enc, err := app.srv.Key.CFBEncrypter(*u.Pwd)

	if err != nil {
		tool.RetFail(w, err)
	}

	pwd := string(enc)
	u.Pwd = &pwd

	email := *u.Email
	u.Email = nil
	updater := bson.M{"$set": &u}

	updateAt := time.Now().Local()
	u.UpdateAt = &updateAt

	res := app.srv.Mongo.GetColl(model.TUser).FindOneAndUpdate(context.Background(), bson.M{"email": email}, updater)

	if res.Err() != nil {
		tool.RetFail(w, res.Err())
		return
	}

	res.Decode(&u)

	err = app.removeRedisCaptcha(&email)

	if err != nil {
		log.Println(err)
	}

	app.cacheSign(w, u.ID.Hex())

}

func (app *App) cacheSign(w http.ResponseWriter, uid string) {
	dur := time.Hour * 24 * 15
	exp := time.Now().Add(dur).Unix()
	tk, err := app.srv.Key.GenJWT(&jwt.StandardClaims{
		ExpiresAt: exp,
		Issuer:    "fuRan",
		Audience:  uid,
	})

	if err != nil {
		tool.RetFail(w, err)
		return
	}

	sign := strings.Split(tk, ".")[2]
	cmd := app.srv.Rdb.SetEX(context.Background(), sign, uid, dur)

	if cmd.Err() != nil {
		tool.RetFail(w, cmd.Err())
		return
	}

	tool.RetOk(w, sign)
}

// Profile 获取用户档案
func (app *App) Profile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	uid, err := primitive.ObjectIDFromHex(r.Header.Get("uid"))
	if err != nil {
		tool.RetFail(w, err)
		return
	}
	res := app.srv.Mongo.GetColl(model.TUser).FindOne(context.Background(), bson.M{"_id": uid}, options.FindOne().SetProjection(bson.M{
		"pwd": 0,
	}))

	if res.Err() != nil {
		w.WriteHeader(http.StatusUnauthorized)
		tool.RetFail(w, res.Err())
		return
	}

	var u model.User

	res.Decode(&u)

	tool.RetOk(w, u)
}

// CreateUser 新增用户
func (app *App) CreateUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	var u model.User
	err = json.Unmarshal(body, &u)
	if err != nil {
		tool.RetFail(w, err)
		return
	}

	res, err := app.insertOneUser(&u)

	if err != nil {
		tool.RetFail(w, err)
		return
	}

	tool.RetOk(w, res.InsertedID.(primitive.ObjectID).Hex())
}

// RemoveUser 删除用户
func (app *App) RemoveUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	uid, err := primitive.ObjectIDFromHex(ps.ByName("uid"))

	if err != nil {
		tool.RetFail(w, err)
		return
	}

	res := app.srv.Mongo.GetColl(model.TUser).FindOneAndDelete(context.Background(), bson.M{
		"_id": uid,
	})

	if res.Err() != nil {
		tool.RetFail(w, res.Err())
		return
	}

	tool.RetOk(w, "删除成功")
}

// UpdateUser 修改用户
func (app *App) UpdateUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		tool.RetFail(w, err)
		return
	}

	if len(body) == 0 {
		tool.RetFail(w, errors.New("not has body"))
	}

	var u model.User

	err = json.Unmarshal(body, &u)

	if err != nil {
		tool.RetFail(w, err)
		return
	}

	if u.ID == nil {
		tool.RetFail(w, errors.New("用户id不能为空"))
		return
	}

	if u.Pwd != nil {
		enc, err := app.srv.Key.CFBEncrypter(*u.Pwd)

		if err != nil {
			tool.RetFail(w, err)
		}
		pwd := string(enc)
		u.Pwd = &pwd
	}

	u.Email = nil
	updater := bson.M{"$set": &u}

	updateAt := time.Now().Local()
	u.UpdateAt = &updateAt
	if u.Roles == nil || len(u.Roles) == 0 {
		updater["$unset"] = bson.M{"roles": ""}
	}

	res := app.srv.Mongo.GetColl(model.TUser).FindOneAndUpdate(context.Background(), bson.M{"_id": *u.ID}, updater)

	if res.Err() != nil {
		errMsg := res.Err().Error()
		if strings.Contains(errMsg, "dup key") {
			errMsg = "该邮箱已经被注册"
		}

		tool.RetFail(w, errors.New(errMsg))
		return
	}

	tool.RetOk(w, "操作成功")
}

// UserList 查找用户
func (app *App) UserList(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
		params["$or"] = []bson.M{
			{"name": bson.M{"$regex": p.Keyword}},
			{"email": bson.M{"$regex": p.Keyword}},
		}
	}

	opt := options.Find().
		SetProjection(bson.M{
			"pwd": 0,
		})

	if p.Limit != nil {
		opt.SetLimit(*p.Limit)
	} else {
		opt.SetLimit(10)
	}

	if p.Skip != nil {
		opt.SetSkip(*p.Skip)
	}
	t := app.srv.Mongo.GetColl(model.TUser)

	cur, err := t.Find(context.Background(), params, opt)

	if err != nil {
		tool.RetFail(w, err)
		return
	}

	var users []model.User
	err = cur.All(context.Background(), &users)

	if err != nil {
		tool.RetFail(w, err)
		return
	}

	total, err := t.CountDocuments(context.Background(), params)

	if err != nil {
		tool.RetFail(w, err)
		return
	}

	tool.RetOkWithTotal(w, users, total)
}

func (app *App) insertOneUser(u *model.User) (*mongo.InsertOneResult, error) {

	if err := app.srv.Validate.Struct(u); err != nil {
		return nil, err
	}

	enc, err := app.srv.Key.CFBEncrypter(*u.Pwd)

	email := strings.ToLower(strings.Replace(*u.Email, " ", "", -1))
	pwd := string(enc)
	now := time.Now().Local()
	if err != nil {
		return nil, err
	}

	u.Email = &email
	u.Pwd = &pwd
	u.CreateAt = &now

	t := app.srv.Mongo.GetColl(model.TUser)

	res, err := t.InsertOne(context.Background(), u)

	if err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "dup key") {
			errMsg = "该邮箱已经被注册"
		}

		return nil, errors.New(errMsg)

	}
	return res, nil
}

// UserValidateEmail email 校验
func (app *App) UserValidateEmail(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	p := struct {
		Email *string `query:"email,omitempty" validate:"omitempty,email,required"`
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

	total, err := app.srv.Mongo.GetColl(model.TUser).CountDocuments(context.Background(), bson.M{"email": &p.Email})

	if err != nil {
		tool.RetFail(w, err)
		return
	}

	if total != 0 {
		tool.RetFail(w, errors.New("该邮箱已经被注册"))
		return
	}

	tool.RetOk(w, "validate key")
}
