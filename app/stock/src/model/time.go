/*
 * @Author: fuRan NgeKaworu@gmail.com
 * @Date: 2023-02-18 20:38:05
 * @LastEditors: fuRan NgeKaworu@gmail.com
 * @LastEditTime: 2023-02-19 18:33:46
 * @FilePath: /stock/stock-go/src/model/time.go
 * @Description:
 *
 * Copyright (c) 2023 by ${git_name_email}, All Rights Reserved.
 */

package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

type MyTime struct {
	time.Time
}

func (p *MyTime) UnmarshalJSON(b []byte) (err error) {
	s := string(b)
	t, err := time.Parse(time.RFC3339, s[1:len(s)-1])

	p.Time = t
	return
}

func (p *MyTime) MarshalJSON() (b []byte, err error) {
	return []byte(`"` + p.Time.Format(time.RFC3339) + `"`), nil
}

func (p *MyTime) UnmarshalBSONValue(btype bsontype.Type, b []byte) (err error) {
	vr := bsonrw.NewBSONValueReader(btype, b)
	dec, err := bson.NewDecoder(vr)
	if err != nil {
		return err
	}
	var t time.Time
	err = dec.Decode(&t)
	if err != nil {
		return err
	}
	p.Time = t
	return nil
}

func (p *MyTime) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bson.MarshalValue(p.Time)
}
