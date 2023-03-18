/*
 * @Author: fuRan NgeKaworu@gmail.com
 * @Date: 2023-01-30 18:05:33
 * @LastEditors: fuRan NgeKaworu@gmail.com
 * @LastEditTime: 2023-02-26 15:22:31
 * @FilePath: /stock/stock-go/src/app/exchange_controller.go
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

	"log"
	"net/http"
	"time"

	"github.com/NgeKaworu/stock/src/model"
	"github.com/NgeKaworu/stock/src/util"
	"github.com/hetiansu5/urlquery"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (app *App) ExchangeList(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	code := ps.ByName("code")
	if code == "" {
		util.RetFail(w, errors.New("code is null "))
		return
	}

	query := struct {
		Skip  int64 `query:"skip"`
		Limit int64 `query:"limit"`
	}{
		Limit: 10,
		Skip:  0,
	}

	err := urlquery.Unmarshal([]byte(r.URL.RawQuery), &query)
	if err != nil {
		util.RetFail(w, err)
		return
	}

	tExchange := app.mongo.GetColl(model.TExchange)

	filter := bson.M{
		"code": code,
	}

	res, err := tExchange.Find(
		context.Background(),
		filter,
		options.Find().
			SetLimit(query.Limit).
			SetSkip(query.Skip).
			SetSort(bson.M{"createAt": -1}),
	)

	if err != nil {
		util.RetFail(w, err)
		return
	}

	list := make([]model.Exchange, 0)

	if err := res.All(context.Background(), &list); err != nil {
		util.RetFail(w, err)
		return
	}

	count, err := tExchange.CountDocuments(context.Background(), filter)
	if err != nil {
		util.RetFail(w, err)
		return
	}

	util.RetOkWithTotal(w, list, count)

}

func (app *App) ExchangeUpsert(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	exchange := new(model.Exchange)

	err = json.Unmarshal(body, &exchange)
	if err != nil {
		util.RetFail(w, err)
		return
	}

	err = app.validate.Struct(exchange)
	if err != nil {
		util.RetFailWithTrans(w, err, app.trans)
		return
	}

	tPosition := app.mongo.GetColl(model.TPosition)
	pos := new(model.Position)

	if err := tPosition.FindOne(context.Background(), bson.M{"_id": exchange.Code}).Decode(&pos); err != nil {
		log.Println("find position error: ", err)
	}

	now := time.Now().Local()

	if pos.Code == nil {
		pos.Code = exchange.Code
		pos.CreateAt.Time = now
		pos.StopLoss = -30
		pos.StopProfit = 15
	}

	pos.UpdateAt.Time = now

	id := ps.ByName("id")

	isEdit := id != ""

	tExchange := app.mongo.GetColl(model.TExchange)

	if !isEdit {
		log.Println("creat exchange")
		pos.TotalShare += exchange.CurrentShare
		pos.TotalCapital += exchange.TransactionPrice * exchange.CurrentShare
		pos.TotalDividend += exchange.CurrentDividend

		exchange.UpdateAt.Time = now
		exchange.ID = primitive.NewObjectID()

	}

	if isEdit {
		log.Println("edit exchange")
		oid, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			util.RetFail(w, err)
			return
		}

		exchange.ID = oid

		old := new(model.Exchange)
		err = tExchange.FindOne(context.Background(), bson.M{"_id": oid}).Decode(&old)
		if err != nil {
			util.RetFail(w, err)
			return
		}

		pos.TotalShare += exchange.CurrentShare - old.CurrentShare
		pos.TotalCapital += exchange.TransactionPrice*exchange.CurrentShare - old.TransactionPrice*old.CurrentShare
		pos.TotalDividend += exchange.CurrentDividend - old.CurrentDividend
	}

	if _, err = tExchange.UpdateOne(context.Background(), bson.M{"_id": exchange.ID}, bson.M{"$set": exchange}, options.Update().SetUpsert(true)); err != nil {
		util.RetFail(w, err)
		return
	}

	if _, err = tPosition.UpdateOne(context.Background(), bson.M{"_id": pos.Code}, bson.M{"$set": pos}, options.Update().SetUpsert(true)); err != nil {
		util.RetFail(w, err)
		return
	}

	util.RetOk(w, "ok")

}

func (app *App) ExchangeDelete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	if id == "" {
		util.RetFail(w, errors.New("id is null "))
		return
	}

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		util.RetFail(w, err)
		return
	}

	tExchange := app.mongo.GetColl(model.TExchange)
	exchange := new(model.Exchange)

	err = tExchange.FindOne(context.Background(), bson.M{"_id": oid}).Decode(&exchange)
	if err != nil {
		util.RetFail(w, err)
		return
	}

	tPosition := app.mongo.GetColl(model.TPosition)
	pos := new(model.Position)

	if err := tPosition.FindOne(context.Background(), bson.M{"_id": exchange.Code}).Decode(&pos); err != nil {
		util.RetFail(w, err)
		return
	}

	pos.UpdateAt.Time = time.Now().Local()

	pos.TotalShare -= exchange.CurrentShare
	pos.TotalCapital -= exchange.TransactionPrice * exchange.CurrentShare
	pos.TotalDividend -= exchange.CurrentDividend

	if _, err = tExchange.DeleteOne(context.Background(), bson.M{"_id": exchange.ID}); err != nil {
		util.RetFail(w, err)
		return
	}

	if _, err = tPosition.UpdateOne(context.Background(), bson.M{"_id": pos.Code}, bson.M{"$set": pos}, options.Update().SetUpsert(true)); err != nil {
		util.RetFail(w, err)
		return
	}

	util.RetOk(w, "ok")
}
