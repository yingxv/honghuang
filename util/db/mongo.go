/*
 * @Author: fuRan NgeKaworu@gmail.com
 * @Date: 2023-03-19 19:20:49
 * @LastEditors: fuRan NgeKaworu@gmail.com
 * @LastEditTime: 2023-03-19 22:44:50
 * @FilePath: /honghuang/util/db/mongo.go
 * @Description:
 *
 * Copyright (c) 2023 by ${git_name_email}, All Rights Reserved.
 */
package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

type MongoInit = func(*mongo.Client, string)

// MongoClient 非关系型数据库引擎
type MongoClient struct {
	MgEngine *mongo.Client //关系型数据库引擎
	Mdb      string
}

// NewMongoClient 实例工厂
func NewMongoClient() *MongoClient {
	return &MongoClient{}
}

// Open 开启连接池x
func (d *MongoClient) Open(mg, mdb string, mongoInit MongoInit) error {
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
		panic(err)
	}

	d.MgEngine = db
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

	//初始化数据库
	mongoInit(session, mdb)

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
