/*
 * @Author: fuRan NgeKaworu@gmail.com
 * @Date: 2023-01-30 18:05:33
 * @LastEditors: fuRan NgeKaworu@gmail.com
 * @LastEditTime: 2023-02-19 21:51:59
 * @FilePath: /stock/stock-go/src/app/position_controller.go
 * @Description:
 *
 * Copyright (c) 2023 by ${git_name_email}, All Rights Reserved.
 */
package app

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/NgeKaworu/stock/src/model"
	"github.com/NgeKaworu/stock/src/util"
	"github.com/hetiansu5/urlquery"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (app *App) Position(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	code := ps.ByName("code")

	query := struct {
		Skip      int64 `query:"skip"`
		Limit     int64 `query:"limit"`
		Omitempty bool  `query:"omitempty"`
	}{
		Skip:      0,
		Limit:     10,
		Omitempty: true,
	}

	err := urlquery.Unmarshal([]byte(r.URL.RawQuery), &query)
	if err != nil {
		util.RetFail(w, err)
		return
	}

	tPosition := app.mongo.GetColl(model.TPosition)
	match := bson.M{}

	if code != "" {
		match["_id"] = code
	}

	if query.Omitempty {
		match["total_share"] = bson.M{
			"$gt": 0,
		}
	}

	pipeline := []bson.M{
		{"$match": match},
		{"$skip": query.Skip},
		{"$limit": query.Limit},
		{"$lookup": bson.M{
			"from": "t_stock",
			"let": bson.M{
				"stockCode": "$_id",
			},
			"pipeline": []bson.M{
				{
					"$match": bson.M{
						"$expr": bson.M{"$eq": bson.A{"$$stockCode", "$code"}},
					},
				},
				{
					"$sort": bson.M{"createAt": -1},
				},
				{
					"$limit": 1,
				},
			},
			"as": "stock",
		}},
		{"$unwind": "$stock"},
	}

	res, err := tPosition.Aggregate(
		context.Background(),
		pipeline,
	)

	if err != nil {
		util.RetFail(w, err)
		return
	}

	list := make([]model.Position, 0)
	if err := res.All(context.Background(), &list); err != nil {
		util.RetFail(w, err)
		return
	}

	if code != "" && len(list) > 0 {
		util.RetOk(w, list[0])
		return
	}

	count, err := tPosition.CountDocuments(context.Background(), match, options.Count())
	if err != nil {
		util.RetFail(w, err)
		return
	}

	util.RetOkWithTotal(w, list, count)

}

func (app *App) PositionUpsert(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	code := ps.ByName("code")
	if code == "" {
		util.RetFail(w, errors.New("code is null "))
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		util.RetFail(w, err)
		return
	}

	if len(body) == 0 {
		util.RetFail(w, errors.New("not has body"))
		return
	}

	position := struct {
		StopProfit float64 `json:"stopProfit" bson:"stop_profit" validate:"required,numeric,gte=0"` // 止盈点
		StopLoss   float64 `json:"stopLoss" bson:"stop_loss" validate:"required,numeric,lte=0"`     // 止损点
	}{}

	err = json.Unmarshal(body, &position)
	if err != nil {
		util.RetFail(w, err)
		return
	}

	err = app.validate.Struct(position)
	if err != nil {
		util.RetFailWithTrans(w, err, app.trans)
		return
	}

	if res := app.mongo.
		GetColl(model.TPosition).
		FindOneAndUpdate(context.Background(),
			bson.M{"_id": code},
			bson.M{"$set": &position},
		); res.Err() != nil {
		util.RetFail(w, res.Err())
		return
	}

	util.RetOk(w, "ok")
}
