package app

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/yingxv/honghuang/time-mgt-go/src/models"

	"github.com/julienschmidt/httprouter"
	"github.com/yingxv/honghuang/util/tool"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// AddRecord 添加记录
func (d *App) AddRecord(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
		"tid": "请至少选一个标签",
	})

	if err != nil {
		tool.RetFail(w, err)
		return
	}

	t := d.srv.Mongo.GetColl(models.TRecord)
	var deration time.Duration

	last := t.FindOne(context.Background(), bson.M{"uid": uid}, options.FindOne().SetSort(bson.M{"createAt": -1}))
	if last.Err() == nil {
		var record models.Record
		err = last.Decode(&record)
		if err == nil {
			deration = time.Now().Local().Sub(*record.CreateAt)
		}
	}

	p["uid"] = uid
	p["createAt"] = time.Now().Local()
	p["deration"] = deration

	res, err := t.InsertOne(context.Background(), p)
	if err != nil {
		tool.RetFail(w, err)
		return
	}

	tool.RetOk(w, res.InsertedID.(primitive.ObjectID).Hex())
}

// SetRecord 更新记录
func (d *App) SetRecord(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
		"event": "请填写发生了什么",
		"tid":   "请至少选一个标签",
		"id":    "ID不能为空",
	})

	if err != nil {
		tool.RetFail(w, err)
		return
	}

	t := d.srv.Mongo.GetColl(models.TRecord)
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

// RemoveRecord 删除记录
func (d *App) RemoveRecord(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	t := d.srv.Mongo.GetColl(models.TRecord)

	res := t.FindOneAndDelete(context.Background(), bson.M{"_id": id, "uid": uid})

	if res.Err() != nil {
		tool.RetFail(w, res.Err())
		return
	}

	tool.RetOk(w, "删除成功")
}

// ListRecord record列表
func (d *App) ListRecord(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	t := d.srv.Mongo.GetColl(models.TRecord)

	total, err := t.CountDocuments(context.Background(), bson.M{
		"uid": uid,
	})

	if err != nil {
		tool.RetFail(w, err)
		return
	}

	cur, err := t.Find(context.Background(), bson.M{
		"uid": uid,
	}, options.Find().SetSort(bson.M{"createAt": -1}).SetSkip(skip).SetLimit(limit))

	if err != nil {
		tool.RetFail(w, err)
		return
	}

	list := make([]models.Record, 0)

	err = cur.All(context.Background(), &list)
	if err != nil {
		tool.RetFail(w, err)
		return
	}
	tool.RetOkWithTotal(w, list, total)
}

// StatisticRecord 统计record
func (d *App) StatisticRecord(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	match := bson.M{
		"uid": uid,
	}

	if len(body) != 0 {
		p, err := tool.ParSup().ConvJSON(body)
		if err != nil {
			tool.RetFail(w, err)
			return
		}
		if dateRange, ok := p["dateRange"].([]interface{}); ok {
			if len(dateRange) == 2 {
				match["createAt"] = bson.M{
					"$gte": dateRange[0],
					"$lte": dateRange[1],
				}
			}
		}

		if tids, ok := p["tids"].([]interface{}); ok {
			if len(tids) > 0 {
				match["tid"] = bson.M{"$in": tids}
			}
		}
	}

	pipe := []bson.M{
		{"$match": match},
		{
			"$unwind": bson.M{
				"path":                       "$tid",
				"preserveNullAndEmptyArrays": true,
			},
		},
	}

	if tid, ok := match["tid"]; ok {
		pipe = append(pipe, bson.M{
			"$match": bson.M{
				"tid": tid,
			},
		})
	}

	pipe = append(pipe,
		bson.M{"$group": bson.M{
			"_id":      "$tid",
			"deration": bson.M{"$sum": "$deration"},
		}},
		bson.M{"$sort": bson.M{
			"deration": -1,
		}},
	)

	t := d.srv.Mongo.GetColl(models.TRecord)
	cur, err := t.Aggregate(context.Background(), pipe)

	if err != nil {
		tool.RetFail(w, err)
		return
	}
	record := make([]models.Record, 0)
	err = cur.All(context.Background(), &record)
	if err != nil {
		tool.RetFail(w, err)
		return
	}

	tool.RetOk(w, record)
}
