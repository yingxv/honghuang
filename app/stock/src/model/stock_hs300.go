package model

// Hs300 沪深300
var Hs300 = map[string]string{
	"600000": "01", //浦发银行
	"600008": "01", //首创股份
	"600009": "01", //上海机场
	"600010": "01", //包钢股份
	"600011": "01", //华能国际
	"600015": "01", //华夏银行
	"600016": "01", //民生银行
	"600018": "01", //上港集团
	"600019": "01", //宝钢股份
	"600021": "01", //上海电力
	"600023": "01", //浙能电力
	"600028": "01", //中国石化
	"600029": "01", //南方航空
	"600030": "01", //中信证券
	"600031": "01", //三一重工
	"600036": "01", //招商银行
	"600038": "01", //中直股份
	"600048": "01", //保利地产
	"600050": "01", //中国联通
	"600061": "01", //国投资本
	"600066": "01", //宇通客车
	"600068": "01", //葛洲坝
	"600074": "01", //ST保千里
	"600085": "01", //同仁堂
	"600089": "01", //特变电工
	"600100": "01", //同方股份
	"600104": "01", //上汽集团
	"600109": "01", //国金证券
	"600111": "01", //北方稀土
	"600115": "01", //东方航空
	"600118": "01", //中国卫星
	"600153": "01", //建发股份
	"600157": "01", //永泰能源
	"600170": "01", //上海建工
	"600177": "01", //雅戈尔
	"600188": "01", //兖州煤业
	"600196": "01", //复星医药
	"600208": "01", //新湖中宝
	"600219": "01", //南山铝业
	"600221": "01", //海航控股
	"600233": "01", //圆通速递
	"600271": "01", //航天信息
	"600276": "01", //恒瑞医药
	"600297": "01", //广汇汽车
	"600309": "01", //万华化学
	"600332": "01", //白云山
	"600340": "01", //华夏幸福
	"600352": "01", //浙江龙盛
	"600362": "01", //江西铜业
	"600369": "01", //西南证券
	"600372": "01", //中航电子
	"600373": "01", //中文传媒
	"600376": "01", //首开股份
	"600383": "01", //金地集团
	"600390": "01", //五矿资本
	"600406": "01", //国电南瑞
	"600415": "01", //小商品城
	"600436": "01", //片仔癀
	"600482": "01", //中国动力
	"600485": "01", //信威集团
	"600489": "01", //中金黄金
	"600498": "01", //烽火通信
	"600518": "01", //康美药业
	"600519": "01", //贵州茅台
	"600522": "01", //中天科技
	"600535": "01", //天士力
	"600547": "01", //山东黄金
	"600549": "01", //厦门钨业
	"600570": "01", //恒生电子
	"600583": "01", //海油工程
	"600585": "01", //海螺水泥
	"600588": "01", //用友网络
	"600606": "01", //绿地控股
	"600637": "01", //东方明珠
	"600649": "01", //城投控股
	"600660": "01", //福耀玻璃
	"600663": "01", //陆家嘴
	"600674": "01", //川投能源
	"600682": "01", //南京新百
	"600685": "01", //中船防务
	"600688": "01", //上海石化
	"600690": "01", //青岛海尔
	"600703": "01", //三安光电
	"600704": "01", //物产中大
	"600705": "01", //中航资本
	"600739": "01", //辽宁成大
	"600741": "01", //华域汽车
	"600795": "01", //国电电力
	"600804": "01", //鹏博士
	"600816": "01", //安信信托
	"600820": "01", //隧道股份
	"600827": "01", //百联股份
	"600837": "01", //海通证券
	"600871": "01", //石化油服
	"600886": "01", //国投电力
	"600887": "01", //伊利股份
	"600893": "01", //航发动力
	"600895": "01", //张江高科
	"600900": "01", //长江电力
	"600909": "01", //华安证券
	"600919": "01", //江苏银行
	"600926": "01", //杭州银行
	"600958": "01", //东方证券
	"600959": "01", //江苏有线
	"600977": "01", //中国电影
	"600999": "01", //招商证券
	"601006": "01", //大秦铁路
	"601009": "01", //南京银行
	"601012": "01", //隆基股份
	"601018": "01", //宁波港
	"601021": "01", //春秋航空
	"601088": "01", //中国神华
	"601099": "01", //太平洋
	"601111": "01", //中国国航
	"601117": "01", //中国化学
	"601118": "01", //海南橡胶
	"601155": "01", //新城控股
	"601163": "01", //三角轮胎
	"601166": "01", //兴业银行
	"601169": "01", //北京银行
	"601186": "01", //中国铁建
	"601198": "01", //东兴证券
	"601211": "01", //国泰君安
	"601212": "01", //白银有色
	"601216": "01", //君正集团
	"601225": "01", //陕西煤业
	"601228": "01", //广州港
	"601229": "01", //上海银行
	"601288": "01", //农业银行
	"601318": "01", //中国平安
	"601328": "01", //交通银行
	"601333": "01", //广深铁路
	"601336": "01", //新华保险
	"601375": "01", //中原证券
	"601377": "01", //兴业证券
	"601390": "01", //中国中铁
	"601398": "01", //工商银行
	"601555": "01", //东吴证券
	"601600": "01", //中国铝业
	"601601": "01", //中国太保
	"601607": "01", //上海医药
	"601608": "01", //中信重工
	"601611": "01", //中国核建
	"601618": "01", //中国中冶
	"601628": "01", //中国人寿
	"601633": "01", //长城汽车
	"601668": "01", //中国建筑
	"601669": "01", //中国电建
	"601688": "01", //华泰证券
	"601718": "01", //际华集团
	"601727": "01", //上海电气
	"601766": "01", //中国中车
	"601788": "01", //光大证券
	"601800": "01", //中国交建
	"601818": "01", //光大银行
	"601857": "01", //中国石油
	"601866": "01", //中远海发
	"601872": "01", //招商轮船
	"601877": "01", //正泰电器
	"601878": "01", //浙商证券
	"601881": "01", //中国银河
	"601888": "01", //中国国旅
	"601898": "01", //中煤能源
	"601899": "01", //紫金矿业
	"601901": "01", //方正证券
	"601919": "01", //中远海控
	"601933": "01", //永辉超市
	"601939": "01", //建设银行
	"601958": "01", //金钼股份
	"601966": "01", //玲珑轮胎
	"601985": "01", //中国核电
	"601988": "01", //中国银行
	"601989": "01", //中国重工
	"601991": "01", //大唐发电
	"601992": "01", //金隅集团
	"601997": "01", //贵阳银行
	"601998": "01", //中信银行
	"603160": "01", //汇顶科技
	"603799": "01", //华友钴业
	"603833": "01", //欧派家居
	"603858": "01", //步长制药
	"603993": "01", //洛阳钼业
	"000001": "02", //平安银行
	"000002": "02", //万科A
	"000008": "02", //神州高铁
	"000060": "02", //中金岭南
	"000063": "02", //中兴通讯
	"000069": "02", //华侨城A
	"000100": "02", //TCL集团
	"000157": "02", //中联重科
	"000166": "02", //申万宏源
	"000333": "02", //美的集团
	"000338": "02", //潍柴动力
	"000402": "02", //金融街
	"000413": "02", //东旭光电
	"000415": "02", //渤海金控
	"000423": "02", //东阿阿胶
	"000425": "02", //徐工机械
	"000503": "02", //海虹控股
	"000538": "02", //云南白药
	"000540": "02", //中天金融
	"000559": "02", //万向钱潮
	"000568": "02", //泸州老窖
	"000623": "02", //吉林敖东
	"000625": "02", //长安汽车
	"000627": "02", //天茂集团
	"000630": "02", //铜陵有色
	"000651": "02", //格力电器
	"000671": "02", //阳光城
	"000686": "02", //东北证券
	"000709": "02", //河钢股份
	"000723": "02", //美锦能源
	"000725": "02", //京东方A
	"000728": "02", //国元证券
	"000738": "02", //航发控制
	"000750": "02", //国海证券
	"000768": "02", //中航飞机
	"000776": "02", //广发证券
	"000783": "02", //长江证券
	"000792": "02", //盐湖股份
	"000826": "02", //启迪桑德
	"000839": "02", //中信国安
	"000858": "02", //五粮液
	"000876": "02", //新希望
	"000895": "02", //双汇发展
	"000898": "02", //鞍钢股份
	"000938": "02", //紫光股份
	"000959": "02", //首钢股份
	"000961": "02", //中南建设
	"000963": "02", //华东医药
	"000983": "02", //西山煤电
	"001979": "02", //招商蛇口
	"002007": "02", //华兰生物
	"002008": "02", //大族激光
	"002024": "02", //苏宁易购
	"002027": "02", //分众传媒
	"002044": "02", //美年健康
	"002065": "02", //东华软件
	"002074": "02", //国轩高科
	"002081": "02", //金螳螂
	"002142": "02", //宁波银行
	"002146": "02", //荣盛发展
	"002153": "02", //石基信息
	"002174": "02", //游族网络
	"002202": "02", //金风科技
	"002230": "02", //科大讯飞
	"002236": "02", //大华股份
	"002241": "02", //歌尔股份
	"002252": "02", //上海莱士
	"002292": "02", //奥飞娱乐
	"002294": "02", //信立泰
	"002304": "02", //洋河股份
	"002310": "02", //东方园林
	"002352": "02", //顺丰控股
	"002385": "02", //大北农
	"002411": "02", //必康股份
	"002415": "02", //海康威视
	"002424": "02", //贵州百灵
	"002426": "02", //胜利精密
	"002450": "02", //康得新
	"002456": "02", //欧菲科技
	"002460": "02", //赣锋锂业
	"002465": "02", //海格通信
	"002466": "02", //天齐锂业
	"002468": "02", //申通快递
	"002470": "02", //金正大
	"002475": "02", //立讯精密
	"002500": "02", //山西证券
	"002508": "02", //老板电器
	"002555": "02", //三七互娱
	"002558": "02", //巨人网络
	"002572": "02", //索菲亚
	"002594": "02", //比亚迪
	"002601": "02", //龙蟒佰利
	"002602": "02", //世纪华通
	"002608": "02", //江苏国信
	"002624": "02", //完美世界
	"002673": "02", //西部证券
	"002714": "02", //牧原股份
	"002736": "02", //国信证券
	"002739": "02", //万达电影
	"002797": "02", //第一创业
	"002831": "02", //裕同科技
	"002839": "02", //张家港行
	"002841": "02", //视源股份
	"300003": "02", //乐普医疗
	"300015": "02", //爱尔眼科
	"300017": "02", //网宿科技
	"300024": "02", //机器人
	"300027": "02", //华谊兄弟
	"300033": "02", //同花顺
	"300059": "02", //东方财富
	"300070": "02", //碧水源
	"300072": "02", //三聚环保
	"300122": "02", //智飞生物
	"300124": "02", //汇川技术
	"300136": "02", //信维通信
	"300144": "02", //宋城演艺
	"300251": "02", //光线传媒
	"300315": "02", //掌趣科技
}
