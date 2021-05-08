package models

import xt "github.com/free-utils-go/xorm_type_assist"

//Unicode
type UnihanChar struct {
	Unid       rune   `xorm:"not null pk comment('Unicode') INT(32)"`
	UnicodeHex string `xorm:"not null comment('Unicode十六进制表示') unique VARCHAR(6)"`
}

//汉程链接
type HanChengChar struct {
	Unid rune   `xorm:"not null pk comment('Unicode') INT(32)"`
	Url  string `xorm:"not null comment('汉程链接') unique TEXT"`
}

//基本解释
type HanCheng struct {
	Unid    rune     `xorm:"not null pk comment('Unicode') INT(32)"`
	Content []string `xorm:"comment('汉程字典') TEXT"`
}

//繁体字映射
type TraditionalId struct {
	Rid   int  `xorm:"not null pk autoincr INT(11)"`
	Unid  rune `xorm:"not null comment('Unicode') INT(32)"`
	UnidS rune `xorm:"not null comment('最简字Unicode') index INT(32)"`
}

//变体字映射
type VariantId struct {
	Rid   int  `xorm:"not null pk autoincr INT(11)"`
	Unid  rune `xorm:"not null comment('Unicode') INT(32)"`
	UnidS rune `xorm:"not null comment('最简字Unicode') index INT(32)"`
}

//古异体字
type VariantGu struct {
	Unid rune `xorm:"not null pk comment('Unicode') INT(32)"`
}

//常用编码表
type BianMa struct {
	Unid    rune   `xorm:"not null pk comment('Unicode') INT(32)"`
	WuBi86  string `xorm:"comment('五笔86码') TEXT"`
	WuBi98  string `xorm:"comment('五笔98码') TEXT"`
	CangJie string `xorm:"comment('仓颉码') TEXT"`
	SiJiao  string `xorm:"comment('四角码') TEXT"`
	GuiFan  string `xorm:"comment('规范汉字码') TEXT"`
}

//民俗
type MinSu struct {
	Msid          int         `xorm:"not null pk autoincr INT(11)"`
	IsSurname     xt.BoolType `xorm:"comment('姓名学（姓氏）,~bool') VARCHAR(1)"`
	SurnameGender string      `xorm:"comment('姓氏性别') VARCHAR(1)"`
	WuXing        string      `xorm:"comment('五行') VARCHAR(1)"`
	Lucky         string      `xorm:"comment('幸运') VARCHAR(4)"`
	Regular       xt.BoolType `xorm:"comment('常用,~bool') VARCHAR(1)"`
}

//民俗关联
type MinSuId struct {
	Rid  int  `xorm:"not null pk autoincr INT(11)"`
	Msid int  `xorm:"not null comment('民俗id') index INT(11)"`
	Unid rune `xorm:"not null comment('Unicode') index INT(32)"`
}

//字形结构
type Glyph struct {
	Unid                     rune        `xorm:"not null pk comment('Unicode') INT(32)"`
	AsRadical                xt.BoolType `xorm:"comment('用作偏旁,~bool') VARCHAR(1)"`
	Radical                  string      `xorm:"comment('部首') VARCHAR(1)"`
	RadicalStroke            int         `xorm:"comment('部首笔画') TINYINT(4)"`
	Stroke                   int         `xorm:"comment('总笔画') INT(11)"`
	SimplifiedRadical        string      `xorm:"comment('简体部首') VARCHAR(1)"`
	SimplifiedRadicalStroke  int         `xorm:"comment('简体部首笔画') TINYINT(4)"`
	SimplifiedTotalStroke    int         `xorm:"comment('简体总笔画') INT(11)"`
	TraditionalRadical       string      `xorm:"comment('繁体部首') VARCHAR(1)"`
	TraditionalRadicalStroke int         `xorm:"comment('繁体部首笔画') TINYINT(4)"`
	TraditionalTotalStroke   int         `xorm:"comment('总笔画') INT(11)"`
	ShouWei                  string      `xorm:"not null comment('首尾分解查字') TEXT"`
	BuJian                   string      `xorm:"comment('汉字部件构造') TEXT"`
	BiHao                    string      `xorm:"comment('笔顺编号') TEXT"`
	BiDu                     string      `xorm:"comment('笔顺读写') TEXT"`
}

//姓名学笔画数
type ScienceStroke struct {
	Unid          rune `xorm:"not null pk comment('Unicode') INT(11)"`
	ScienceStroke int  `xorm:"default 0 comment('姓名学笔画数') INT(11)"`
}

//音韵参考
type YinYun struct {
	Unid     rune   `xorm:"not null pk comment('Unicode') INT(32)"`
	ShangGu  string `xorm:"comment('上古音') TEXT"`
	GuangYun string `xorm:"comment('广韵') TEXT"`
	PingShui string `xorm:"comment('平水韵') TEXT"`
}

//拼音
type PinYin struct {
	Pid    int    `xorm:"not null pk autoincr INT(11)"`
	PinYin string `xorm:"not null comment('拼音') unique TEXT"`
}

//拼音关联
type PinYinId struct {
	Rid  int  `xorm:"not null pk autoincr INT(11)"`
	Pid  int  `xorm:"not null comment('拼音id') index INT(32)"`
	Unid rune `xorm:"not null comment('Unicode') index INT(32)"`
}

//国语拼音关联
type GuoYuId struct {
	Rid  int  `xorm:"not null pk autoincr INT(11)"`
	Pid  int  `xorm:"not null comment('拼音id') index INT(32)"`
	Unid rune `xorm:"not null comment('Unicode') index INT(32)"`
}

//注音
type ZhuYin struct {
	Zid    int    `xorm:"not null pk autoincr INT(11)"`
	ZhuYin string `xorm:"not null comment('注音') unique TEXT"`
}

//注音关联
type ZhuYinId struct {
	Rid  int  `xorm:"not null pk autoincr INT(11)"`
	Zid  int  `xorm:"not null comment('注音id') index INT(32)"`
	Unid rune `xorm:"not null comment('Unicode') index INT(32)"`
}

//唐音
type TangYin struct {
	Tid     int    `xorm:"not null pk autoincr INT(11)"`
	TangYin string `xorm:"not null comment('唐音') unique TEXT"`
}

//唐音关联
type TangYinId struct {
	Rid  int  `xorm:"not null pk autoincr INT(11)"`
	Tid  int  `xorm:"not null comment('唐音id') index INT(32)"`
	Unid rune `xorm:"not null comment('Unicode') index INT(32)"`
}

//粤音
type YueYin struct {
	Yid    int    `xorm:"not null pk autoincr INT(11)"`
	YueYin string `xorm:"not null comment('粤音') unique TEXT"`
}

//粤音关联
type YueYinId struct {
	Rid  int  `xorm:"not null pk autoincr INT(11)"`
	Yid  int  `xorm:"not null comment('粤音id') index INT(32)"`
	Unid rune `xorm:"not null comment('Unicode') index INT(32)"`
}

//闽南
type MinNanYin struct {
	Mnid      int    `xorm:"not null pk autoincr INT(11)"`
	MinNanYin string `xorm:"not null comment('闽南音') unique TEXT"`
}

//闽南关联
type MinNanYinId struct {
	Rid  int  `xorm:"not null pk autoincr INT(11)"`
	Mnid int  `xorm:"not null comment('Unicode') index INT(32)"`
	Unid rune `xorm:"not null comment('成语id') index INT(32)"`
}

//索引参考
type SuoYin struct {
	Unid rune   `xorm:"not null pk comment('Unicode') INT(32)"`
	Gu   string `xorm:"comment('古文字诂林') TEXT"`
	Xun  string `xorm:"comment('故训彙纂') TEXT"`
	Shuo string `xorm:"comment('说文解字') TEXT"`
	Kang string `xorm:"comment('康熙字典') TEXT"`
	Han  string `xorm:"comment('汉语字典') TEXT"`
	Ci   string `xorm:"comment('辞海') TEXT"`
}

//详细解释（新华字典）
type XinHua struct {
	Unid    rune     `xorm:"not null pk comment('Unicode') INT(32)"`
	Content []string `xorm:"comment('新华字典') TEXT"`
}

//汉语字典（汉语大字典）
type HanDa struct {
	Unid    rune     `xorm:"not null pk comment('Unicode') INT(32)"`
	Content []string `xorm:"comment('汉语大字典') TEXT"`
}

//康熙字典
type KangXi struct {
	Unid    rune     `xorm:"not null pk comment('Unicode') INT(32)"`
	Stroke  int      `xorm:"comment('康熙笔画') INT(11)"`
	Brief   string   `xorm:"comment('简要') TEXT"`
	Content []string `xorm:"comment('解释') TEXT"`
	Url     string   `xorm:"comment('扫描资料链接') TEXT"`
}

//说文解字
type ShuoWen struct {
	Unid    rune     `xorm:"not null pk comment('Unicode') INT(32)"`
	Brief   string   `xorm:"comment('简要') TEXT"`
	Content []string `xorm:"comment('详解') TEXT"`
	Url     string   `xorm:"comment('图像链接') TEXT"`
}

//字源演变
type YanBian struct {
	Unid         rune   `xorm:"not null pk comment('Unicode') INT(32)"`
	JiaGuUrl     string `xorm:"comment('甲骨文图像链接') TEXT"`
	JinUrl       string `xorm:"comment('金文图像链接') TEXT"`
	XiaoZhuanUrl string `xorm:"comment('小篆图像链接') TEXT"`
	KaiUrl       string `xorm:"comment('楷体图像链接') TEXT"`
}

//成语
type ChengYu struct {
	Cyid    int    `xorm:"not null pk autoincr INT(11)"`
	ChengYu string `xorm:"not null comment('成语') unique TEXT"`
	Url     string `xorm:"not null comment('链接') unique TEXT"`
}

//相关成语
type ChengYuId struct {
	Rid  int  `xorm:"not null pk autoincr INT(11)"`
	Cyid int  `xorm:"not null comment('成语id') index INT(32)"`
	Unid rune `xorm:"not null comment('Unicode') index INT(32)"`
}

//诗词
type ShiCi struct {
	Scid  int    `xorm:"not null pk autoincr INT(11)"`
	ShiCi string `xorm:"not null comment('诗词') unique TEXT"`
	Url   string `xorm:"not null comment('链接') unique TEXT"`
}

//相关诗词
type ShiCiId struct {
	Rid  int  `xorm:"not null pk autoincr INT(11)"`
	Scid int  `xorm:"not null comment('诗词id') index INT(32)"`
	Unid rune `xorm:"not null comment('Unicode') index INT(32)"`
}

//词语
type CiYu struct {
	Cid  int    `xorm:"not null pk autoincr INT(11)"`
	CiYu string `xorm:"not null comment('词语') unique TEXT"`
	Url  string `xorm:"not null comment('链接') unique TEXT"`
}

//相关词语
type CiYuId struct {
	Rid  int  `xorm:"not null pk autoincr INT(11)"`
	Cid  int  `xorm:"not null comment('词语id') index INT(32)"`
	Unid rune `xorm:"not null comment('Unicode') index INT(32)"`
}

//最简汉字表
type HanChar struct {
	Unid rune   `xorm:"not null pk comment('Unicode') INT(32)"`
	Ch   string `xorm:"not null comment('汉字') unique VARCHAR(1)"`
}
