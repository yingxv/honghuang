package app

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/NgeKaworu/to-do-list-go/src/models"
	"github.com/NgeKaworu/util/tool"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// AddTask 添加记录
func (d *App) AddTask(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	uid, err := primitive.ObjectIDFromHex(r.Header.Get("uid"))
	if err != nil {
		tool.RetFail(w, err)
		return
	}

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

	p, err := tool.ParSup().ConvJSON(body)
	if err != nil {
		tool.RetFail(w, err)
		return
	}

	err = tool.Required(p, map[string]string{
		"title": "请填写任务名",
		"level": "请选一个优先级",
	})

	if err != nil {
		tool.RetFail(w, err)
		return
	}

	t := d.srv.Mongo.GetColl(models.TTask)

	p["uid"] = uid
	p["createAt"] = time.Now().Local()

	res, err := t.InsertOne(context.Background(), p)
	if err != nil {
		tool.RetFail(w, err)
		return
	}

	tool.RetOk(w, res.InsertedID.(primitive.ObjectID).Hex())
}

// SetTask 更新记录
func (d *App) SetTask(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	uid, err := primitive.ObjectIDFromHex(r.Header.Get("uid"))
	if err != nil {
		tool.RetFail(w, err)
		return
	}

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

	p, err := tool.ParSup().ConvJSON(body)
	if err != nil {
		tool.RetFail(w, err)
		return
	}

	err = tool.Required(p, map[string]string{
		"id": "ID不能为空",
	})

	if err != nil {
		tool.RetFail(w, err)
		return
	}

	t := d.srv.Mongo.GetColl(models.TTask)

	p["uid"] = uid
	p["updateAt"] = time.Now().Local()

	id := p["id"]
	delete(p, "id")

	res := t.FindOneAndUpdate(context.Background(),
		bson.M{"_id": id},
		bson.M{"$set": p},
	)
	if res.Err() != nil {
		tool.RetFail(w, res.Err())
		return
	}

	tool.RetOk(w, "修改成功")
}

// RemoveTask 删除记录
func (d *App) RemoveTask(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	uid, err := primitive.ObjectIDFromHex(r.Header.Get("uid"))
	if err != nil {
		tool.RetFail(w, err)
		return
	}
	id, err := primitive.ObjectIDFromHex(ps.ByName("id"))
	if err != nil {
		tool.RetFail(w, err)
		return
	}

	t := d.srv.Mongo.GetColl(models.TTask)

	res := t.FindOneAndDelete(context.Background(), bson.M{"_id": id, "uid": uid})

	if res.Err() != nil {
		tool.RetFail(w, res.Err())
		return
	}

	tool.RetOk(w, "删除成功")
}

// ListTask record列表
func (d *App) ListTask(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	q := r.URL.Query()
	l := q.Get("limit")
	s := q.Get("skip")

	uid, err := primitive.ObjectIDFromHex(r.Header.Get("uid"))
	if err != nil {
		tool.RetFail(w, err)
		return
	}

	limit, _ := strconv.ParseInt(l, 10, 64)
	skip, _ := strconv.ParseInt(s, 10, 64)

	done, _ := strconv.ParseBool(q.Get("done"))

	t := d.srv.Mongo.GetColl(models.TTask)

	cur, err := t.Find(context.Background(), bson.M{
		"uid":  uid,
		"done": done,
	}, options.Find().SetSort(bson.D{
		{Key: "level", Value: -1},
		{Key: "createAt", Value: -1},
	}).SetSkip(skip).SetLimit(limit))

	if err != nil {
		tool.RetFail(w, err)
		return
	}

	list := make([]models.MainTask, 0)

	err = cur.All(context.Background(), &list)
	if err != nil {
		tool.RetFail(w, err)
		return
	}

	total, err := t.CountDocuments(context.Background(), bson.M{"uid": uid})

	if err != nil {
		tool.RetFail(w, err)
		return
	}

	tool.RetOkWithTotal(w, list, total)
}
