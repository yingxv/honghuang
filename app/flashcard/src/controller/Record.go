package controller

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/hetiansu5/urlquery"
	"github.com/julienschmidt/httprouter"
	"github.com/yingxv/flashcard-go/src/model"
	"github.com/yingxv/honghuang/util/tool"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (c *Controller) RecordCreate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	record := new(model.Record)
	err = json.Unmarshal(body, &record)
	if err != nil {
		tool.RetFail(w, err)
		return
	}

	err = c.srv.Validate.Struct(record)
	if err != nil {
		tool.RetFailWithTrans(w, err, c.srv.Trans)
		return
	}

	now := time.Now()
	record.UID = &uid
	record.CreateAt = &now

	res, err := c.srv.Mongo.GetColl(model.TRecord).InsertOne(context.Background(), record)

	if err != nil {
		tool.RetFail(w, err)
		return
	}

	tool.RetOk(w, res.InsertedID)

}

func (c *Controller) RecordRemove(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	ret := c.srv.Mongo.GetColl(model.TRecord).FindOneAndDelete(context.Background(), bson.M{
		"_id": &id,
		"uid": &uid,
	})

	if ret.Err() != nil {
		tool.RetFail(w, ret.Err())
		return
	}

	tool.RetOk(w, "OK")
}

func (c *Controller) RecordUpdate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	record := new(model.Record)

	err = json.Unmarshal(body, &record)

	if err != nil {
		tool.RetFail(w, err)
		return
	}

	if record.ID == nil {
		tool.RetFail(w, errors.New("ID是必填字段"))
		return
	}

	now := time.Now()
	record.UpdateAt = &now

	ret := c.srv.Mongo.GetColl(model.TRecord).FindOneAndUpdate(context.Background(), bson.M{
		"_id": record.ID,
		"uid": &uid,
	}, bson.M{
		"$set": record,
	})

	if ret.Err() != nil {
		tool.RetFail(w, ret.Err())
		return
	}

	tool.RetOk(w, "OK")
}

func (c *Controller) RecordMultiUpdate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	record := new(struct {
		IDs        []string            `json:"ids" validate:"required"`
		UID        *primitive.ObjectID `json:"uid,omitempty" bson:"uid,omitempty" `               // 所有者id
		UpdateAt   *time.Time          `json:"updateAt,omitempty" bson:"updateAt,omitempty" `     // 更新时间
		ReviewAt   *time.Time          `json:"reviewAt,omitempty" bson:"reviewAt,omitempty" `     // 复习时间
		CooldownAt *time.Time          `json:"cooldownAt,omitempty" bson:"cooldownAt,omitempty" ` // 冷却时间
		Exp        int64               `json:"exp,omitempty" bson:"exp,omitempty" `               // 复习熟练度
		Tag        string              `json:"tag,omitempty" bson:"tag,omitempty" `               // 标签
		Mode       int64               `json:"mode,omitempty" bson:"mode,omitempty" `             // 模式: 0: 关键字; 1: 全文
	})

	err = json.Unmarshal(body, &record)

	if err != nil {
		tool.RetFail(w, err)
		return
	}

	now := time.Now()
	record.UpdateAt = &now

	ids := make([]primitive.ObjectID, 0, len(record.IDs))
	for _, id := range record.IDs {
		oid, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			tool.RetFail(w, err)
			return
		}
		ids = append(ids, oid)
	}

	_, err = c.srv.Mongo.GetColl(model.TRecord).UpdateMany(context.Background(), bson.M{
		"_id": bson.M{"$in": ids},
		"uid": &uid,
	}, bson.M{
		"$set": record,
	})

	if err != nil {
		tool.RetFail(w, err)
		return
	}

	tool.RetOk(w, "OK")
}

func (c *Controller) RecordList(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	uid, err := primitive.ObjectIDFromHex(r.Header.Get("uid"))
	if err != nil {
		tool.RetFail(w, err)
		return
	}

	type Convertor struct {
		Type     *string `query:"type,omitempty" validate:"omitempty,oneof=store enable cooling done"`
		Sort     *string `query:"sort,omitempty" validate:"required_with=OrderBy"`
		OrderBy  *int    `query:"orderby,omitempty" validate:"omitempty,oneof=1 -1,required_with=Sort"`
		Skip     *int64  `query:"skip,omitempty" validate:"omitempty,min=0"`
		Limit    *int64  `query:"limit,omitempty" validate:"omitempty,min=0"`
		InReview *bool   `query:"inReview,omitempty"`
	}

	convertor := new(Convertor)

	err = urlquery.Unmarshal([]byte(r.URL.RawQuery), &convertor)
	if err != nil {
		tool.RetFail(w, err)
		return
	}

	err = c.srv.Validate.Struct(convertor)
	if err != nil {
		tool.RetFailWithTrans(w, err, c.srv.Trans)
		return
	}

	filter := bson.M{
		"uid": uid,
	}

	if convertor.InReview != nil {
		filter["inReview"] = convertor.InReview
	}

	if convertor.Type != nil {
		switch *convertor.Type {
		case "store":
			{
				filter["cooldownAt"] = nil
				filter["exp"] = bson.M{"$ne": 100}
				break
			}
		case "enable":
			{
				filter["cooldownAt"] = bson.M{"$lte": time.Now()}
				filter["exp"] = bson.M{"$ne": 100}
				break
			}
		case "cooling":
			{
				filter["cooldownAt"] = bson.M{"$gt": time.Now()}
				break
			}
		case "done":
			{
				filter["exp"] = 100
				break
			}
		}
	}
	opt := options.FindOptions{
		Skip:  convertor.Skip,
		Limit: convertor.Limit,
	}

	if opt.Limit == nil {
		opt.SetLimit(10)
	}

	if convertor.Sort != nil && convertor.OrderBy != nil {
		opt.SetSort(bson.D{
			{Key: *convertor.Sort, Value: *convertor.OrderBy},
		})
	}

	t := c.srv.Mongo.GetColl(model.TRecord)
	cursor, err := t.Find(context.Background(), filter, &opt)
	if err != nil {
		tool.RetFail(w, err)
		return
	}

	res := make([]map[string]interface{}, 0)

	err = cursor.All(context.Background(), &res)
	if err != nil {
		tool.RetFail(w, err)
		return
	}

	count, err := t.CountDocuments(context.Background(), filter)
	if err != nil {
		tool.RetFail(w, err)
		return
	}
	tool.RetOkWithTotal(w, res, count)
}

func (c *Controller) RecordMgtList(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	uid, err := primitive.ObjectIDFromHex(r.Header.Get("uid"))
	if err != nil {
		tool.RetFail(w, err)
		return
	}

	type Convertor struct {
		Sort    *string `query:"sort,omitempty" validate:"required_with=OrderBy"`
		OrderBy *int    `query:"orderby,omitempty" validate:"omitempty,oneof=1 -1,required_with=Sort"`
		Skip    *int64  `query:"skip,omitempty" validate:"omitempty,min=0"`
		Limit   *int64  `query:"limit,omitempty" validate:"omitempty,min=0"`

		CooldownAt  *[]string `query:"cooldownAt,omitempty"`
		CreateAt    *[]string `query:"createAt,omitempty"`
		ReviewAt    *[]string `query:"reviewAt,omitempty"`
		UpdateAt    *[]string `query:"updateAt,omitempty"`
		Exp         *int64    `query:"exp,omitempty"`
		InReview    *bool     `query:"inReview,omitempty"`
		Source      *string   `query:"source,omitempty"`
		Translation *string   `query:"translation,omitempty"`
		Tag         *string   `query:"tag,omitempty"`
		Mode        *int64    `query:"mode,omitempty"`
		Finished    bool      `query:"finished"`
		Planing     bool      `query:"planing"`
	}

	convertor := new(Convertor)

	err = urlquery.Unmarshal([]byte(r.URL.RawQuery), &convertor)
	if err != nil {
		tool.RetFail(w, err)
		return
	}

	err = c.srv.Validate.Struct(convertor)
	if err != nil {
		tool.RetFailWithTrans(w, err, c.srv.Trans)
		return
	}

	filter := bson.M{
		"uid": uid,
	}

	if convertor.Planing {
		filter["cooldownAt"] = bson.M{"$ne": nil}
	} else {
		filter["cooldownAt"] = nil
	}

	if convertor.CooldownAt != nil {
		startTime, err := time.Parse(time.RFC1123, (*convertor.CooldownAt)[0])
		if err != nil {
			tool.RetFail(w, err)
			return
		}
		endTime, err := time.Parse(time.RFC1123, (*convertor.CooldownAt)[1])
		if err != nil {
			tool.RetFail(w, err)
			return
		}

		filter["cooldownAt"] = bson.M{
			"$gte": startTime,
			"$lte": endTime,
		}
	}

	if convertor.CreateAt != nil {
		startTime, err := time.Parse(time.RFC1123, (*convertor.CreateAt)[0])
		if err != nil {
			tool.RetFail(w, err)
			return
		}
		endTime, err := time.Parse(time.RFC1123, (*convertor.CreateAt)[1])
		if err != nil {
			tool.RetFail(w, err)
			return
		}
		filter["createAt"] = bson.M{
			"$gte": startTime,
			"$lte": endTime,
		}
	}

	if convertor.ReviewAt != nil {
		startTime, err := time.Parse(time.RFC1123, (*convertor.ReviewAt)[0])
		if err != nil {
			tool.RetFail(w, err)
			return
		}
		endTime, err := time.Parse(time.RFC1123, (*convertor.ReviewAt)[1])
		if err != nil {
			tool.RetFail(w, err)
			return
		}
		filter["reviewAt"] = bson.M{
			"$gte": startTime,
			"$lte": endTime,
		}
	}

	if convertor.UpdateAt != nil {
		startTime, err := time.Parse(time.RFC1123, (*convertor.UpdateAt)[0])
		if err != nil {
			tool.RetFail(w, err)
			return
		}
		endTime, err := time.Parse(time.RFC1123, (*convertor.UpdateAt)[1])
		if err != nil {
			tool.RetFail(w, err)
			return
		}
		filter["updateAt"] = bson.M{
			"$gte": startTime,
			"$lte": endTime,
		}
	}

	filter["exp"] = bson.M{"$ne": 100}
	if convertor.Finished {
		filter["exp"] = 100
	}

	if convertor.Exp != nil {
		filter["exp"] = convertor.Exp
	}

	if convertor.InReview != nil {
		filter["inReview"] = convertor.InReview
	}

	if convertor.Source != nil {
		filter["source"] = bson.M{
			"$regex":   *convertor.Source,
			"$options": "i", // "i" 表示不区分大小写
		}
	}

	if convertor.Translation != nil {
		filter["translation"] = bson.M{
			"$regex":   *convertor.Translation,
			"$options": "i", // "i" 表示不区分大小写
		}
	}

	if convertor.Tag != nil {
		filter["tag"] = bson.M{
			"$regex":   *convertor.Tag,
			"$options": "i", // "i" 表示不区分大小写
		}
	}

	if convertor.Mode != nil {
		filter["mode"] = convertor.Mode
	}

	opt := options.FindOptions{
		Skip:  convertor.Skip,
		Limit: convertor.Limit,
	}

	if opt.Limit == nil {
		opt.SetLimit(10)
	}

	if convertor.Sort != nil && convertor.OrderBy != nil {
		opt.SetSort(bson.D{
			{Key: *convertor.Sort, Value: *convertor.OrderBy},
		})
	}

	t := c.srv.Mongo.GetColl(model.TRecord)
	cursor, err := t.Find(context.Background(), filter, &opt)
	if err != nil {
		tool.RetFail(w, err)
		return
	}

	res := make([]map[string]interface{}, 0)

	err = cursor.All(context.Background(), &res)
	if err != nil {
		tool.RetFail(w, err)
		return
	}

	count, err := t.CountDocuments(context.Background(), filter)
	if err != nil {
		tool.RetFail(w, err)
		return
	}
	tool.RetOkWithTotal(w, res, count)
}

func (c *Controller) RecordReview(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	converter := struct {
		IDs *[]*primitive.ObjectID `json:"ids,omitempty" bson:"ids" validate:"required"`
	}{}

	if err := json.Unmarshal(body, &converter); err != nil {
		tool.RetFail(w, err)
		return
	}

	if err := c.srv.Validate.Struct(converter); err != nil {
		tool.RetFailWithTrans(w, err, c.srv.Trans)
		return
	}

	filter := bson.M{
		"uid": uid,
		"_id": bson.M{"$in": converter.IDs},
	}
	res, err := c.srv.Mongo.GetColl(model.TRecord).UpdateMany(context.Background(), filter, bson.M{"$set": bson.M{"inReview": true}})

	if err != nil {
		tool.RetFail(w, err)
		return
	}

	tool.RetOk(w, res)

}

func (c *Controller) RecordReviewStop(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	converter := struct {
		IDs *[]*primitive.ObjectID `json:"ids,omitempty" bson:"ids" validate:"required"`
	}{}

	if err := json.Unmarshal(body, &converter); err != nil {
		tool.RetFail(w, err)
		return
	}

	if err := c.srv.Validate.Struct(converter); err != nil {
		tool.RetFailWithTrans(w, err, c.srv.Trans)
		return
	}

	filter := bson.M{
		"uid": uid,
		"_id": bson.M{"$in": converter.IDs},
	}
	res, err := c.srv.Mongo.GetColl(model.TRecord).UpdateMany(context.Background(), filter, bson.M{"$set": bson.M{"inReview": false}})

	if err != nil {
		tool.RetFail(w, err)
		return
	}

	tool.RetOk(w, res)

}

func (c *Controller) RecordReviewCooldownAtClear(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	converter := struct {
		IDs *[]*primitive.ObjectID `json:"ids,omitempty" bson:"ids" validate:"required"`
	}{}

	if err := json.Unmarshal(body, &converter); err != nil {
		tool.RetFail(w, err)
		return
	}

	if err := c.srv.Validate.Struct(converter); err != nil {
		tool.RetFailWithTrans(w, err, c.srv.Trans)
		return
	}

	filter := bson.M{
		"uid": uid,
		"_id": bson.M{"$in": converter.IDs},
	}
	res, err := c.srv.Mongo.GetColl(model.TRecord).UpdateMany(context.Background(), filter, bson.M{"$set": bson.M{"cooldownAt": nil}})

	if err != nil {
		tool.RetFail(w, err)
		return
	}

	tool.RetOk(w, res)

}

func (c *Controller) RecordReviewAll(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	uid, err := primitive.ObjectIDFromHex(r.Header.Get("uid"))
	if err != nil {
		tool.RetFail(w, err)
		return
	}
	filter := bson.M{
		"uid":      uid,
		"exp":      bson.M{"$ne": 100},
		"inReview": false,
		"cooldownAt": bson.M{
			"$lte": time.Now(),
		},
	}

	res, err := c.srv.Mongo.GetColl(model.TRecord).UpdateMany(context.Background(), filter, bson.M{
		"$set": bson.M{
			"inReview": true,
		},
	})
	if err != nil {
		tool.RetFail(w, err)
		return
	}

	tool.RetOk(w, res)

}

func (c *Controller) RecordSetReviewResult(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	var record struct {
		ID         *primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty" validate:"required"`                // id
		CooldownAt *time.Time          `json:"cooldownAt,omitempty" bson:"cooldownAt,omitempty" validate:"required"` // 冷却时间
		Exp        int64               `json:"exp,omitempty" bson:"exp,omitempty" validate:"gte=0"`                  // 复习熟练度
		InReview   bool                `json:"inReview,omitempty" bson:"inReview" `                                  // 是否在复习中
		ReviewAt   *time.Time          `json:"reviewAt,omitempty" bson:"reviewAt,omitempty" `                        // 复习时间
	}

	if err := json.Unmarshal(body, &record); err != nil {
		tool.RetFail(w, err)
		return
	}

	if err := c.srv.Validate.Struct(record); err != nil {
		tool.RetFailWithTrans(w, err, c.srv.Trans)
		return
	}
	now := time.Now()
	record.ReviewAt = &now
	ret := c.srv.Mongo.GetColl(model.TRecord).FindOneAndUpdate(context.Background(), bson.M{
		"_id": record.ID,
		"uid": uid,
	}, bson.M{
		"$set": record,
	})

	if ret.Err() != nil {
		tool.RetFail(w, ret.Err())
		return
	}

	tool.RetOk(w, "OK")

}
