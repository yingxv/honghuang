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

	"github.com/yingxv/honghuang/to-do-list-go/src/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

func Init(session *mongo.Client, mdb string) {
	// 记录表
	t := session.Database(mdb).Collection(models.TTask)
	indexView := t.Indexes()
	_, err := indexView.CreateMany(context.Background(), []mongo.IndexModel{
		{Keys: bsonx.Doc{bsonx.Elem{Key: "uid", Value: bsonx.Int32(1)}}},
		{Keys: bsonx.Doc{bsonx.Elem{Key: "createAt", Value: bsonx.Int32(-1)}}},
		{Keys: bsonx.Doc{bsonx.Elem{Key: "updateAt", Value: bsonx.Int32(-1)}}},
	})

	if err != nil {
		log.Println(err)
	}
}

func WithoutInit(session *mongo.Client, mdb string) {}
