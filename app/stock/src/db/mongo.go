/*
 * @Author: fuRan NgeKaworu@gmail.com
 * @Date: 2020-11-14 11:06:01
 * @LastEditors: fuRan NgeKaworu@gmail.com
 * @LastEditTime: 2023-02-19 20:18:05
 * @FilePath: /stock/stock-go/src/db/mongo.go
 * @Description:
 *
 * Copyright (c) 2023 by ${git_name_email}, All Rights Reserved.
 */
package db

import (
	"context"
	"log"
	"time"

	"github.com/NgeKaworu/stock/src/model"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

// MongoClient 关系型数据库引擎
type MongoClient struct {
	MgEngine *mongo.Client //文档型数据库引擎
	Mdb      string
}

// NewClient 实例工厂
func NewMongoClient() *MongoClient {
	return &MongoClient{}
}

// Open 开启连接池
func (d *MongoClient) Open(mg, mdb string, initdb bool) error {
	d.Mdb = mdb
	ops := options.Client().ApplyURI(mg)
	p := uint64(39000)
	ops.MaxPoolSize = &p
	ops.WriteConcern = writeconcern.New(writeconcern.J(true), writeconcern.W(1))
	ops.ReadPreference = readpref.PrimaryPreferred()
	db, err := mongo.NewClient(ops)

	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.Connect(ctx)
	if err != nil {
		return err
	}

	//err = db.Ping(ctx, readpref.PrimaryPreferred())
	//if err != nil {
	//	log.Println("ping err", err)
	//}

	d.MgEngine = db

	//初始化数据库
	if initdb {
		var session *mongo.Client
		session, err = mongo.NewClient(ops)
		if err != nil {
			panic(err)
		}
		err = session.Connect(context.Background())
		if err != nil {
			panic(err)
		}
		defer session.Disconnect(context.Background())

		// 每股数据
		stock := session.Database(mdb).Collection(model.TStock)
		indexes := stock.Indexes()
		// undo 刷索引
		_, err = indexes.CreateMany(context.Background(), []mongo.IndexModel{
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

	return nil
}

// GetColl 获取表名
func (d *MongoClient) GetColl(coll string) *mongo.Collection {
	col, _ := d.MgEngine.Database(d.Mdb).Collection(coll).Clone()
	return col
}

// Close 关闭连接池
func (d *MongoClient) Close() {
	d.MgEngine.Disconnect(context.Background())
}
