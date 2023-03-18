package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TEnterpriseIndicator 表名
const TEnterpriseIndicator = "t_enterprise_indicator"

// MainIndicatorRes 主要指标返回体
type MainIndicatorRes struct {
	Result Result `json:"Result"`
}

// Result 主要指标返回值
type Result struct {
	Enterprise []Enterprise `json:"ZhuYaoZhiBiaoList_QiYe"`
	// YinHang   YinHang   `json:"ZhuYaoZhiBiaoList_YinHang"`
	// QuanShang QuanShang `json:"ZhuYaoZhiBiaoList_QuanShang"`
	// BaoXian   BaoXian   `json:"ZhuYaoZhiBiaoList_BaoXian"`
}

// Enterprise 企业指标
type Enterprise struct {
	ID                           *primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	CreateDate                   *time.Time          `json:"CreateDate,omitempty" bson:"create_date,omitempty"`
	Code                         *string             `json:"code,omitempty" bson:"code,omitempty"`                                                     //股票代号
	ReportDate                   *string             `json:"ReportDate,omitempty" bson:"report_date,omitempty"`                                        //报告日期
	Title                        *string             `json:"Title,omitempty" bson:"title,omitempty"`                                                   //报告名称
	Epsjb                        *string             `json:"Epsjb,omitempty" bson:"epsjb,omitempty"`                                                   //基本每股收益(元)
	Epskcjb                      *string             `json:"Epskcjb,omitempty" bson:"epskcjb,omitempty"`                                               //扣非每股收益(元)
	Epsxs                        *string             `json:"Epsxs,omitempty" bson:"epsxs,omitempty"`                                                   //稀释每股收益(元)
	Bps                          *string             `json:"Bps,omitempty" bson:"bps,omitempty"`                                                       //每股净资产(元)
	Mgzbgj                       *string             `json:"Mgzbgj,omitempty" bson:"mgzbgj,omitempty"`                                                 //每股资本公积(元)
	Mgwfplr                      *string             `json:"Mgwfplr,omitempty" bson:"mgwfplr,omitempty"`                                               //每股未分配利润(元)
	Mgjyxjje                     *string             `json:"Mgjyxjje,omitempty" bson:"mgjyxjje,omitempty"`                                             //每股经营现金流(元)
	TotalIncome                  *string             `json:"Totalincome,omitempty" bson:"total_income,omitempty"`                                      //营业总收入(元)
	GrossProfit                  *string             `json:"Grossprofit,omitempty" bson:"gross_profit,omitempty"`                                      //毛利润(元)
	ParentNetProfit              *string             `json:"Parentnetprofit,omitempty" bson:"parent_net_profit,omitempty"`                             //归属净利润(元)
	BuckleNetProfit              *string             `json:"Bucklenetprofit,omitempty" bson:"buckle_net_profit,omitempty"`                             //扣非净利润(元)
	TotalIncomeyoy               *string             `json:"Totalincomeyoy,omitempty" bson:"total_incomeyoy,omitempty"`                                //营业总收入同比增长
	ParentNetProfityoy           *string             `json:"Parentnetprofityoy,omitempty" bson:"parent_net_profityoy,omitempty"`                       //归属净利润同比增长
	BuckleNetProfityoy           *string             `json:"Bucklenetprofityoy,omitempty" bson:"buckle_net_profityoy,omitempty"`                       //扣非净利润同比增长
	TotalIncomeRelativeRatio     *string             `json:"Totalincomerelativeratio,omitempty" bson:"total_income_relative_ratio,omitempty"`          //营业总收入滚动环比增长
	ParentNetProfitRelativeRatio *string             `json:"Parentnetprofitrelativeratio,omitempty" bson:"parent_net_profit_relative_ratio,omitempty"` //归属净利润滚动环比增长
	BuckleNetProfitRelativeRatio *string             `json:"Bucklenetprofitrelativeratio,omitempty" bson:"buckle_net_profit_relative_ratio,omitempty"` //扣非净利润滚动环比增长
	Roejq                        *string             `json:"Roejq,omitempty" bson:"roejq,omitempty"`                                                   //净资产收益率(加权)
	Roekcjq                      *string             `json:"Roekcjq,omitempty" bson:"roekcjq,omitempty"`                                               //净资产收益率(扣非/加权)
	AllCapitalEarningsRate       *string             `json:"Allcapitalearningsrate,omitempty" bson:"all_capital_earnings_rate,omitempty"`              //总资产收益率(加权)
	GrossMargin                  *string             `json:"Grossmargin,omitempty" bson:"gross_margin,omitempty"`                                      //毛利率
	NetInterest                  *string             `json:"Netinterest,omitempty" bson:"net_interest,omitempty"`                                      //净利率
	AccountsRate                 *string             `json:"Accountsrate,omitempty" bson:"accounts_rate,omitempty"`                                    //预收账款/营业收入
	SalesRate                    *string             `json:"Salesrate,omitempty" bson:"sales_rate,omitempty"`                                          //销售净现金流/营业收入
	OperatingRate                *string             `json:"Operatingrate,omitempty" bson:"operating_rate,omitempty"`                                  //经营净现金流/营业收入
	TaxRate                      *string             `json:"Taxrate,omitempty" bson:"tax_rate,omitempty"`                                              //实际税率
	LiquidityRatio               *string             `json:"Liquidityratio,omitempty" bson:"liquidity_ratio,omitempty"`                                //流动比率
	QuickRatio                   *string             `json:"Quickratio,omitempty" bson:"quick_ratio,omitempty"`                                        //速动比率
	CashFlowRatio                *string             `json:"Cashflowratio,omitempty" bson:"cash_flow_ratio,omitempty"`                                 //现金流量比率
	AssetliabilityRatio          *string             `json:"Assetliabilityratio,omitempty" bson:"assetliability_ratio,omitempty"`                      //资产负债率
	EquityMultiplier             *string             `json:"Equitymultiplier,omitempty" bson:"equity_multiplier,omitempty"`                            //权益乘数
	EquityRatio                  *string             `json:"Equityratio,omitempty" bson:"equity_ratio,omitempty"`                                      //产权比率
	TotalAssetsDays              *string             `json:"Totalassetsdays,omitempty" bson:"total_assets_days,omitempty"`                             //总资产周转天数(天)
	InventoryDays                *string             `json:"Inventorydays,omitempty" bson:"inventory_days,omitempty"`                                  //存货周转天数(天)
	AccountsreceivableDays       *string             `json:"Accountsreceivabledays,omitempty" bson:"accountsreceivable_days,omitempty"`                //应收账款周转天数(天)
	TotalassetRate               *string             `json:"Totalassetrate,omitempty" bson:"totalasset_rate,omitempty"`                                //总资产周转率(次)
	InventoryRate                *string             `json:"Inventoryrate,omitempty" bson:"inventory_rate,omitempty"`                                  //存货周转率(次)
	AccountsReceiveableRate      *string             `json:"Accountsreceiveablerate,omitempty" bson:"accounts_receiveable_rate,omitempty"`             //应收账款周转率(次)
}
