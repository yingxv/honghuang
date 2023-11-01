/*
 * @Author: fuRan NgeKaworu@gmail.com
 * @Date: 2020-11-14 11:06:01
 * @LastEditors: fuRan NgeKaworu@gmail.com
 * @LastEditTime: 2023-03-19 23:09:17
 * @FilePath: /honghuang/app/time-mgt/src/creator/mongo.go
 * @Description:
 *
 * Copyright (c) 2023 by ${git_name_email}, All Rights Reserved.
 */
package creator

import (
	"context"
	"log"

	"github.com/yingxv/honghuang/user-center/src/model"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

func Init(session *mongo.Client, mdb string) {
	// 用户表
	t := session.Database(mdb).Collection(model.TUser)

	indexView := t.Indexes()
	_, err := indexView.CreateMany(context.Background(), []mongo.IndexModel{
		{Keys: bsonx.Doc{bsonx.Elem{Key: "email", Value: bsonx.Int32(1)}}, Options: options.Index().SetUnique(true)},
		{Keys: bsonx.Doc{bsonx.Elem{Key: "name", Value: bsonx.Int32(1)}}},
	})

	if err != nil {
		log.Println(err)
	}

	// 角色表
	t = session.Database(mdb).Collection(model.TRole)

	indexView = t.Indexes()
	_, err = indexView.CreateMany(context.Background(), []mongo.IndexModel{
		{Keys: bsonx.Doc{bsonx.Elem{Key: "name", Value: bsonx.Int32(1)}}},
		{Keys: bsonx.Doc{bsonx.Elem{Key: "perms", Value: bsonx.Int32(1)}}},
	})

	if err != nil {
		log.Println(err)
	}

	// 权限表
	t = session.Database(mdb).Collection(model.TPerm)

	indexView = t.Indexes()
	_, err = indexView.CreateMany(context.Background(), []mongo.IndexModel{
		{Keys: bsonx.Doc{bsonx.Elem{Key: "name", Value: bsonx.Int32(1)}}},
		{Keys: bsonx.Doc{bsonx.Elem{Key: "pID", Value: bsonx.Int32(1)}}},
		{Keys: bsonx.Doc{bsonx.Elem{Key: "url", Value: bsonx.Int32(1)}}},
		{Keys: bsonx.Doc{bsonx.Elem{Key: "isMenu", Value: bsonx.Int32(1)}}},
	})

	if err != nil {
		log.Println(err)
	}
}

func WithoutInit(session *mongo.Client, mdb string) {}
