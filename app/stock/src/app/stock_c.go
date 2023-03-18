package app

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/NgeKaworu/stock/src/model"
	"github.com/NgeKaworu/stock/src/util"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (d *App) StockCrawlMany(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	res, err := d.StockCrawlManyService()
	if err != nil {
		util.RetFail(w, err)
	}

	util.RetOk(w, &res)
}

func (d *App) StockCrawlManyService() (*mongo.InsertManyResult, error) {

	allStock := make([]interface{}, 0)
	pool := make(chan bool, 10)
	now := time.Now().Local()
	format, _ := time.Parse("2006-01-02 15:03:05", now.Format("2006-01-02 00:00:00"))

	for k, v := range model.Stocks {
		pool <- true
		go func(key, val string) {
			s := model.NewStock(key, val)
			s.CreateAt = &format
			s.Crawl()
			allStock = append(allStock, s)
			<-pool
		}(k, v)

	}

	t := d.mongo.GetColl(model.TStock)
	_, err := t.DeleteMany(context.Background(), bson.M{
		"createAt": &format,
	})

	if err != nil {
		return nil, err
	}

	res, err := t.InsertMany(context.Background(), allStock)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (d *App) StockList(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	dataTime := r.URL.Query().Get("dataTime")
	times := make([]time.Time, 2)

	err := json.Unmarshal([]byte(dataTime), &times)
	if err != nil {
		util.RetFail(w, err)
		return
	}

	if len(times) != 2 {
		util.RetFail(w, errors.New("dataTime must a range"))
		return
	}

	query := bson.M{
		"createAt": bson.M{
			"$gte": times[0],
			"$lte": times[1],
		},
	}

	t := d.mongo.GetColl(model.TStock)

	c, err := t.Find(context.Background(), &query)
	if err != nil {
		util.RetFail(w, err)
		return
	}

	res := make([]*model.Stock, 0)
	err = c.All(context.Background(), &res)

	if err != nil {
		util.RetFail(w, err)
		return
	}

	util.RetOkWithTotal(w, res, int64(len(res)))
}
