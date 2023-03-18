package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TTask 任务表
const TTask = "t_task"

// Task 主任务
type MainTask struct {
	*Task    `json:",inline" bson:",inline"` // 继承
	ID       *primitive.ObjectID             `json:"id,omitempty" bson:"_id,omitempty"`            // id
	UID      *primitive.ObjectID             `json:"uid,omitempty" bson:"uid,omitempty"`           // uid
	CreateAt *time.Time                      `json:"createAt,omitempty" bson:"createAt,omitempty"` // 创建时间
	UpdateAt *time.Time                      `json:"updateAt,omitempty" bson:"updateAt,omitempty"` // 更新时间
	SubTask  *[]Task                         `json:"subTask,omitempty" bson:"subTask,omitempty"`   // 子任务
	Level    *int64                          `json:"level,omitempty" bson:"level,omitempty"`       // 优先级
}

// Task 任务
type Task struct {
	Title *string `json:"title,omitempty" bson:"title,omitempty"` // 任务标题
	Done  *bool   `json:"done,omitempty" bson:"done,omitempty"`   // 任务是否完成
}
