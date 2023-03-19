/*
 * @Author: fuRan NgeKaworu@gmail.com
 * @Date: 2020-11-14 11:06:01
 * @LastEditors: fuRan NgeKaworu@gmail.com
 * @LastEditTime: 2023-03-19 23:09:12
 * @FilePath: /honghuang/app/stock/src/creator/mongo.go
 * @Description:
 *
 * Copyright (c) 2023 by ${git_name_email}, All Rights Reserved.
 */
package creator

import (
	"context"
	"log"

	"github.com/NgeKaworu/stock/src/model"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

func Init(session *mongo.Client, mdb string) {
	// 每股数据
	stock := session.Database(mdb).Collection(model.TStock)
	indexes := stock.Indexes()
	// undo 刷索引
	_, err := indexes.CreateMany(context.Background(), []mongo.IndexModel{
		{Keys: bsonx.Doc{bsonx.Elem{Key: "classify", Value: bsonx.Int32(-1)}}},
		{Keys: bsonx.Doc{bsonx.Elem{Key: "name", Value: bsonx.Int32(1)}}},
		{Keys: bsonx.Doc{bsonx.Elem{Key: "createAt", Value: bsonx.Int32(1)}}},
		{Keys: bsonx.Doc{bsonx.Elem{Key: "bourseCode", Value: bsonx.Int32(1)}}},
	})
	if err != nil {
		log.Println(err)
	}

	// 持仓
	position := session.Database(mdb).Collection(model.TPosition)
	indexes = position.Indexes()
	// undo 刷索引
	_, err = indexes.CreateMany(context.Background(), []mongo.IndexModel{
		{Keys: bsonx.Doc{bsonx.Elem{Key: "createAt", Value: bsonx.Int32(1)}}},
	})
	if err != nil {
		log.Println(err)
	}

	// 交易
	exchange := session.Database(mdb).Collection(model.TExchange)
	indexes = exchange.Indexes()
	// undo 刷索引
	_, err = indexes.CreateMany(context.Background(), []mongo.IndexModel{
		{Keys: bsonx.Doc{bsonx.Elem{Key: "code", Value: bsonx.Int32(-1)}}},
		{Keys: bsonx.Doc{bsonx.Elem{Key: "createAt", Value: bsonx.Int32(1)}}},
	})
	if err != nil {
		log.Println(err)
	}
}

func WithoutInit(session *mongo.Client, mdb string) {}
