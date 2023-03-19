/*
 * @Author: fuRan NgeKaworu@gmail.com
 * @Date: 2023-03-19 03:04:11
 * @LastEditors: fuRan NgeKaworu@gmail.com
 * @LastEditTime: 2023-03-19 22:45:13
 * @FilePath: /honghuang/app/flashcard/src/creator/mongo.go
 * @Description:
 *
 * Copyright (c) 2023 by ${git_name_email}, All Rights Reserved.
 */
package creator

import (
	"context"

	"github.com/yingxv/flashcard-go/src/model"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

func Init(session *mongo.Client, mdb string) {
	// 记录表
	t := session.Database(mdb).Collection(model.TRecord)
	indexView := t.Indexes()
	_, err := indexView.CreateMany(context.Background(), []mongo.IndexModel{
		{Keys: bsonx.Doc{bsonx.Elem{Key: "createAt", Value: bsonx.Int32(-1)}}},
		{Keys: bsonx.Doc{bsonx.Elem{Key: "reviewAt", Value: bsonx.Int32(-1)}}},
		{Keys: bsonx.Doc{bsonx.Elem{Key: "cooldownAt", Value: bsonx.Int32(-1)}}},
		{Keys: bsonx.Doc{bsonx.Elem{Key: "inReview", Value: bsonx.Int32(1)}}},
		{Keys: bsonx.Doc{bsonx.Elem{Key: "exp", Value: bsonx.Int32(1)}}},
	})

	if err != nil {
		panic(err)
	}

}

func WithoutInit(session *mongo.Client, mdb string) {}
