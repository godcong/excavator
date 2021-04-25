package excavator

//database

import (
	"excavator/models"
	"fmt"
	"net/url"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/godcong/fate"
	"github.com/godcong/fate/config"
	_ "github.com/mattn/go-sqlite3"
	"xorm.io/xorm"
)

// InitSqlite3 ...
func initSqlite3(sqfile string) *xorm.Engine {
	eng, e := xorm.NewEngine("sqlite3", sqfile)
	if e != nil {
		panic(e)
	}
	eng.ShowSQL(true)

	_, e = eng.Exec("PRAGMA journal_mode = OFF;")
	if e != nil {
		Log.Fatal(e)
	}
	return eng
}

const sqlURL = "%s:%s@tcp(%s)/%s?loc=%s&charset=utf8mb4&parseTime=true"

func initMysql(addr, name, pass string) *xorm.Engine {
	u := fmt.Sprintf(sqlURL, name, pass, addr, "excavator", url.QueryEscape("Asia/Shanghai"))
	before := time.Now()
	eng, e := xorm.NewEngine("mysql", u)
	if e != nil {
		Log.Fatal(e)
	}

	eng.ShowSQL(true)
	fmt.Printf("took %v\n", time.Since(before))
	return eng
}

func InitXorm(cfg *config.DatabaseConfig) (eng *xorm.Engine) {
	if cfg.Driver == "sqlite3" {
		return initSqlite3(cfg.File)
	} else if cfg.Driver == "mysql" {
		return initMysql(cfg.Host, cfg.User, cfg.Pwd)
	} else {
		panic("数据库类型不支持")
	}
}

func ResetExc(engine *xorm.Engine) error {
	engine.Query("DELETE FROM bian_ma")
	engine.Query("DELETE FROM cheng_yu")
	engine.Query("DELETE FROM cheng_yu_id")
	engine.Query("DELETE FROM ci_yu")
	engine.Query("DELETE FROM ci_yu_id")
	engine.Query("DELETE FROM glyph")
	engine.Query("DELETE FROM guo_yu_id")
	engine.Query("DELETE FROM han_char")
	engine.Query("DELETE FROM han_cheng")
	engine.Query("DELETE FROM han_da")
	engine.Query("DELETE FROM kang_xi")
	engine.Query("DELETE FROM min_nan_yin")
	engine.Query("DELETE FROM min_nan_yin_id")
	engine.Query("DELETE FROM min_su")
	engine.Query("DELETE FROM min_su_id")
	engine.Query("DELETE FROM pin_yin")
	engine.Query("DELETE FROM pin_yin_id")
	engine.Query("DELETE FROM science_stroke")
	engine.Query("DELETE FROM shi_ci")
	engine.Query("DELETE FROM shi_ci_id")
	engine.Query("DELETE FROM shuo_wen")
	engine.Query("DELETE FROM suo_yin")
	engine.Query("DELETE FROM tang_yin")
	engine.Query("DELETE FROM tang_yin_id")
	engine.Query("DELETE FROM variant_gu")
	engine.Query("DELETE FROM variant_id")
	engine.Query("DELETE FROM xin_hua")
	engine.Query("DELETE FROM yan_bian")
	engine.Query("DELETE FROM yin_yun")
	engine.Query("DELETE FROM yue_yin")
	engine.Query("DELETE FROM yue_yin_id")
	engine.Query("DELETE FROM zhu_yin")
	engine.Query("DELETE FROM zhu_yin")

	return nil
}

func ResetFate(engine *xorm.Engine) error {
	engine.Query("DELETE FROM character")

	return nil
}

// InsertOrUpdate
func InsertOrUpdate(engine *xorm.Engine, c interface{}) (err error) {

	var has bool
	var id int

	switch v := c.(type) {
	case *models.UnihanChar:
		id = c.(*models.UnihanChar).Unid

		if id == 0 {
			has, err = engine.Get(c)
			id = c.(*models.UnihanChar).Unid
		} else {
			has, err = engine.ID(id).Get(&models.UnihanChar{})
		}
	case *models.HanChengChar:
		id = c.(*models.HanChengChar).Unid

		if id == 0 {
			has, err = engine.Get(c)
			id = c.(*models.HanChengChar).Unid
		} else {
			has, err = engine.ID(id).Get(&models.HanChengChar{})
		}
	case *models.HanCheng:
		id = c.(*models.HanCheng).Unid

		if id == 0 {
			has, err = engine.Get(c)
			id = c.(*models.HanCheng).Unid
		} else {
			has, err = engine.ID(id).Get(&models.HanCheng{})
		}
	case *models.TraditionalId:
		id = c.(*models.TraditionalId).Rid

		if id == 0 {
			has, err = engine.Get(c)
			id = c.(*models.TraditionalId).Rid
		} else {
			has, err = engine.ID(id).Get(&models.TraditionalId{})
		}
	case *models.VariantId:
		id = c.(*models.VariantId).Rid

		if id == 0 {
			has, err = engine.Get(c)
			id = c.(*models.VariantId).Rid
		} else {
			has, err = engine.ID(id).Get(&models.VariantId{})
		}
	case *models.VariantGu:
		id = c.(*models.VariantGu).Unid

		if id != 0 {
			has, err = engine.ID(id).Get(&models.VariantGu{})
		}

		if err != nil {
			panic(err)
		}

		if !has {
			_, err = engine.InsertOne(c)
		}

		if err != nil {
			panic(err)
		}

		return nil
	case *models.BianMa:
		id = c.(*models.BianMa).Unid

		if id == 0 {
			has, err = engine.Get(c)
			id = c.(*models.BianMa).Unid
		} else {
			has, err = engine.ID(id).Get(&models.BianMa{})
		}
	case *models.MinSu:
		id = c.(*models.MinSu).Msid

		if id == 0 {
			has, err = engine.Get(c)
			id = c.(*models.MinSu).Msid
		} else {
			has, err = engine.ID(id).Get(&models.MinSu{})
		}
	case *models.MinSuId:
		id = c.(*models.MinSuId).Rid

		if id == 0 {
			has, err = engine.Get(c)
			id = c.(*models.MinSuId).Rid
		} else {
			has, err = engine.ID(id).Get(&models.MinSuId{})
		}
	case *models.Glyph:
		id = c.(*models.Glyph).Unid

		if id == 0 {
			has, err = engine.Get(c)
			id = c.(*models.Glyph).Unid
		} else {
			has, err = engine.ID(id).Get(&models.Glyph{})
		}
	case *models.ScienceStroke:
		id = c.(*models.ScienceStroke).Unid

		if id == 0 {
			has, err = engine.Get(c)
			id = c.(*models.ScienceStroke).Unid
		} else {
			has, err = engine.ID(id).Get(&models.ScienceStroke{})
		}
	case *models.YinYun:
		id = c.(*models.YinYun).Unid

		if id == 0 {
			has, err = engine.Get(c)
			id = c.(*models.YinYun).Unid
		} else {
			has, err = engine.ID(id).Get(&models.YinYun{})
		}
	case *models.PinYin:
		id = c.(*models.PinYin).Pid

		if id == 0 {
			has, err = engine.Get(c)
			id = c.(*models.PinYin).Pid
		} else {
			has, err = engine.ID(id).Get(&models.PinYin{})
		}
	case *models.PinYinId:
		id = c.(*models.PinYinId).Rid

		if id == 0 {
			has, err = engine.Get(c)
			id = c.(*models.PinYinId).Rid
		} else {
			has, err = engine.ID(id).Get(&models.PinYinId{})
		}
	case *models.GuoYuId:
		id = c.(*models.GuoYuId).Rid

		if id == 0 {
			has, err = engine.Get(c)
			id = c.(*models.GuoYuId).Rid
		} else {
			has, err = engine.ID(id).Get(&models.GuoYuId{})
		}
	case *models.ZhuYin:
		id = c.(*models.ZhuYin).Zid

		if id == 0 {
			has, err = engine.Get(c)
			id = c.(*models.ZhuYin).Zid
		} else {
			has, err = engine.ID(id).Get(&models.ZhuYin{})
		}
	case *models.ZhuYinId:
		id = c.(*models.ZhuYinId).Rid

		if id == 0 {
			has, err = engine.Get(c)
			id = c.(*models.ZhuYinId).Rid
		} else {
			has, err = engine.ID(id).Get(&models.ZhuYinId{})
		}
	case *models.TangYin:
		id = c.(*models.TangYin).Tid

		if id == 0 {
			has, err = engine.Get(c)
			id = c.(*models.TangYin).Tid
		} else {
			has, err = engine.ID(id).Get(&models.TangYin{})
		}
	case *models.TangYinId:
		id = c.(*models.TangYinId).Rid

		if id == 0 {
			has, err = engine.Get(c)
			id = c.(*models.TangYinId).Rid
		} else {
			has, err = engine.ID(id).Get(&models.TangYinId{})
		}
	case *models.YueYin:
		id = c.(*models.YueYin).Yid

		if id == 0 {
			has, err = engine.Get(c)
			id = c.(*models.YueYin).Yid
		} else {
			has, err = engine.ID(id).Get(&models.YueYin{})
		}
	case *models.YueYinId:
		id = c.(*models.YueYinId).Rid

		if id == 0 {
			has, err = engine.Get(c)
			id = c.(*models.YueYinId).Rid
		} else {
			has, err = engine.ID(id).Get(&models.YueYinId{})
		}
	case *models.MinNanYin:
		id = c.(*models.MinNanYin).Mnid

		if id == 0 {
			has, err = engine.Get(c)
			id = c.(*models.MinNanYin).Mnid
		} else {
			has, err = engine.ID(id).Get(&models.MinNanYin{})
		}
	case *models.MinNanYinId:
		id = c.(*models.MinNanYinId).Rid

		if id == 0 {
			has, err = engine.Get(c)
			id = c.(*models.MinNanYinId).Rid
		} else {
			has, err = engine.ID(id).Get(&models.MinNanYinId{})
		}
	case *models.SuoYin:
		id = c.(*models.SuoYin).Unid

		if id == 0 {
			has, err = engine.Get(c)
			id = c.(*models.SuoYin).Unid
		} else {
			has, err = engine.ID(id).Get(&models.SuoYin{})
		}
	case *models.XinHua:
		id = c.(*models.XinHua).Unid

		if id == 0 {
			has, err = engine.Get(c)
			id = c.(*models.XinHua).Unid
		} else {
			has, err = engine.ID(id).Get(&models.XinHua{})
		}
	case *models.HanDa:
		id = c.(*models.HanDa).Unid

		if id == 0 {
			has, err = engine.Get(c)
			id = c.(*models.HanDa).Unid
		} else {
			has, err = engine.ID(id).Get(&models.HanDa{})
		}
	case *models.KangXi:
		id = c.(*models.KangXi).Unid

		if id == 0 {
			has, err = engine.Get(c)
			id = c.(*models.KangXi).Unid
		} else {
			has, err = engine.ID(id).Get(&models.KangXi{})
		}
	case *models.ShuoWen:
		id = c.(*models.ShuoWen).Unid

		if id == 0 {
			has, err = engine.Get(c)
			id = c.(*models.ShuoWen).Unid
		} else {
			has, err = engine.ID(id).Get(&models.ShuoWen{})
		}
	case *models.YanBian:
		id = c.(*models.YanBian).Unid

		if id == 0 {
			has, err = engine.Get(c)
			id = c.(*models.YanBian).Unid
		} else {
			has, err = engine.ID(id).Get(&models.YanBian{})
		}
	case *models.ChengYu:
		id = c.(*models.ChengYu).Cyid

		if id == 0 {
			has, err = engine.Get(c)
			id = c.(*models.ChengYu).Cyid
		} else {
			has, err = engine.ID(id).Get(&models.ChengYu{})
		}
	case *models.ChengYuId:
		id = c.(*models.ChengYuId).Rid

		if id == 0 {
			has, err = engine.Get(c)
			id = c.(*models.ChengYuId).Rid
		} else {
			has, err = engine.ID(id).Get(&models.ChengYuId{})
		}
	case *models.ShiCi:
		id = c.(*models.ShiCi).Scid

		if id == 0 {
			has, err = engine.Get(c)
			id = c.(*models.ShiCi).Scid
		} else {
			has, err = engine.ID(id).Get(&models.ShiCi{})
		}
	case *models.ShiCiId:
		id = c.(*models.ShiCiId).Rid

		if id == 0 {
			has, err = engine.Get(c)
			id = c.(*models.ShiCiId).Rid
		} else {
			has, err = engine.ID(id).Get(&models.ShiCiId{})
		}
	case *models.CiYu:
		id = c.(*models.CiYu).Cid

		if id == 0 {
			has, err = engine.Get(c)
			id = c.(*models.CiYu).Cid
		} else {
			has, err = engine.ID(id).Get(&models.CiYu{})
		}
	case *models.CiYuId:
		id = c.(*models.CiYuId).Rid

		if id == 0 {
			has, err = engine.Get(c)
			id = c.(*models.CiYuId).Rid
		} else {
			has, err = engine.ID(id).Get(&models.CiYuId{})
		}
	case *models.HanChar:
		id = c.(*models.HanChar).Unid

		if id == 0 {
			has, err = engine.Get(c)
			id = c.(*models.HanChar).Unid
		} else {
			has, err = engine.ID(id).Get(&models.HanChar{})
		}
	case *fate.Character:
		id := c.(*fate.Character).Ch

		if id == "" {
			has, err = engine.Get(c)
			id = c.(*fate.Character).Ch
		} else {
			has, err = engine.ID(id).Get(&fate.Character{})
		}

		if err != nil {
			panic(err)
		}

		if has {
			_, err = engine.ID(id).Update(c)
		} else {
			_, err = engine.InsertOne(c)
		}

		if err != nil && err.Error() == "No content found to be updated" {
			err = nil
		}

		if err != nil {
			panic(err)
		}

		return nil
	default:
		panic(fmt.Sprintf("%v类型不支持", v))
	}

	if err != nil {
		panic(err)
	}

	if has {
		_, err = engine.ID(id).Update(c)
	} else {
		_, err = engine.InsertOne(c)
	}

	if err != nil && err.Error() == "No content found to be updated" {
		err = nil
	}

	if err != nil {
		panic(err)
	}

	return nil
}
