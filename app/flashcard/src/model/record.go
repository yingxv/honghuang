/*
 * @Author: fuRan NgeKaworu@gmail.com
 * @Date: 2023-03-19 03:04:11
 * @LastEditors: fuRan NgeKaworu@gmail.com
 * @LastEditTime: 2023-10-31 09:42:10
 * @FilePath: /honghuang/app/flashcard/src/model/record.go
 * @Description:
 *
 * Copyright (c) 2023 by ${git_name_email}, All Rights Reserved.
 */
package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// 表名
const TRecord = "t_record"

type Record struct {
	ID          *primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty" `                                                           // id
	UID         *primitive.ObjectID `json:"uid,omitempty" bson:"uid,omitempty" `                                                          // 所有者id
	CreateAt    *time.Time          `json:"createAt,omitempty" bson:"createAt,omitempty" `                                                // 创建时间
	UpdateAt    *time.Time          `json:"updateAt,omitempty" bson:"updateAt,omitempty" `                                                // 更新时间
	ReviewAt    *time.Time          `json:"reviewAt,omitempty" bson:"reviewAt,omitempty" `                                                // 复习时间
	CooldownAt  *time.Time          `json:"cooldownAt,omitempty" bson:"cooldownAt,omitempty" `                                            // 冷却时间
	Source      string              `json:"source,omitempty" bson:"source,omitempty" validate:"required_without=ID" label:"原文"`           // 原文
	Translation string              `json:"translation,omitempty" bson:"translation,omitempty" validate:"required_without=ID" label:"译文"` // 译文
	InReview    bool                `json:"inReview,omitempty" bson:"inReview" `                                                          // 是否在复习中
	Exp         int64               `json:"exp,omitempty" bson:"exp" `                                                                    // 复习熟练度
	Tag         string              `json:"tag" bson:"tag" `                                                                              // 标签
	Mode        int64               `json:"mode" bson:"mode" `                                                                            // 模式: 0: 关键字; 1: 全文
}
