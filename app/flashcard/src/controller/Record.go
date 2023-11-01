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
	record.CooldownAt = &now

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
	record.CooldownAt = &now
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
func (c *Controller) RecordList(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	uid, err := primitive.ObjectIDFromHex(r.Header.Get("uid"))
	if err != nil {
		tool.RetFail(w, err)
		return
	}

	type Convertor struct {
		Type     *string `query:"type,omitempty" validate:"omitempty,oneof=enable cooling done"`
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

func (c *Controller) RecordRandomReview(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	converter := struct {
		Num *int `json:"num,omitempty"`
	}{}

	if len(body) != 0 {
		if err := json.Unmarshal(body, &converter); err != nil {
			tool.RetFail(w, err)
			return
		}
	}

	t := c.srv.Mongo.GetColl(model.TRecord)

	if converter.Num == nil {
		num := 3
		converter.Num = &num
	}

	cursor, err := t.Aggregate(context.Background(), []bson.M{
		{
			"$match": bson.M{
				"uid":        uid,
				"inReview":   false,
				"cooldownAt": bson.M{"$lte": time.Now()},
			},
		},
		{
			"$sample": bson.M{"size": converter.Num},
		},
		{
			"$project": bson.M{"_id": 1},
		},
	})

	if err != nil {
		tool.RetFail(w, err)
		return
	}

	random := make([]model.Record, 0)
	err = cursor.All(context.Background(), &random)
	if err != nil {
		tool.RetFail(w, err)
		return
	}
	ids := make([]primitive.ObjectID, 0)
	for _, record := range random {
		ids = append(ids, *record.ID)
	}

	filter := bson.M{
		"uid": uid,
		"_id": bson.M{"$in": ids},
	}

	res, err := t.UpdateMany(context.Background(), filter, bson.M{
		"$set": bson.M{"inReview": true},
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
