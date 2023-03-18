package db

import (
	"context"
	"time"

	"github.com/yingxv/flashcard-go/src/model"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

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
		panic(err)
	}

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
