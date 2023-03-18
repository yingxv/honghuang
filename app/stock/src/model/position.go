/*
 * @Author: fuRan NgeKaworu@gmail.com
 * @Date: 2023-01-30 15:19:24
 * @LastEditors: fuRan NgeKaworu@gmail.com
 * @LastEditTime: 2023-02-19 20:24:57
 * @FilePath: /stock/stock-go/src/model/position.go
 * @Description:
 *
 * Copyright (c) 2023 by ${git_name_email}, All Rights Reserved.
 */

package model

var TPosition = "t_position"

// Position 持仓
type Position struct {
	Code          *string `json:"code,omitempty" bson:"_id,omitempty"`    // 交易所代码
	Stock         *Stock  `json:"stock,omitempty" bson:"stock,omitempty"` // 股票
	TotalShare    float64 `json:"totalShare" bson:"total_share"`          // 总股份
	TotalCapital  float64 `json:"totalCapital" bson:"total_capital"`      // 总投入
	TotalDividend float64 `json:"totalDividend" bson:"total_dividend"`    // 总派息
	StopProfit    float64 `json:"stopProfit" bson:"stop_profit"`          // 止盈点
	StopLoss      float64 `json:"stopLoss" bson:"stop_loss"`              // 止损点
	CreateAt      MyTime  `json:"createAt" bson:"createAt"`               // 创建时间
	UpdateAt      MyTime  `json:"updateAt" bson:"updateAt"`               // 更新时间
}
