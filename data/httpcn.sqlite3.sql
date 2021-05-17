
-- 先由GB.txt生成，之后查出HanCheng的每个字的url，有就设置，没有就忽略
CREATE TABLE unihan_char ( unid INT (32) PRIMARY KEY,
-- COMMENT 'Unicode' ,
unicode_hex VARCHAR (6) NOT NULL UNIQUE
-- COMMENT 'Unicode十六进制表示' 
) WITHOUT ROWID;
-- COMMENT 'Unicode';

-- 只保留在汉程有unicode收录的字
CREATE TABLE han_cheng_char ( unid INT (32) PRIMARY KEY,
-- COMMENT 'Unicode',
url TEXT NOT NULL UNIQUE,
-- COMMENT '汉程链接',
FOREIGN KEY (unid) REFERENCES unihan_char (unid) ON UPDATE CASCADE ON DELETE CASCADE 
) WITHOUT ROWID;
-- COMMENT '汉程链接';

CREATE TABLE han_cheng ( unid INT (32) PRIMARY KEY,
-- COMMENT 'Unicode',
content TEXT,
-- COMMENT '汉程字典',
FOREIGN KEY (unid) REFERENCES unihan_char (unid) ON UPDATE CASCADE ON DELETE CASCADE 
) WITHOUT ROWID;
-- COMMENT '基本解释';

-- 繁体字关系，
CREATE TABLE traditional_id ( rid INTEGER PRIMARY KEY AUTOINCREMENT,
unid INT (32) NOT NULL UNIQUE,
-- COMMENT 'Unicode' ,
unid_s INT (32) NOT NULL,
-- COMMENT '最简字Unicode' ,
FOREIGN KEY (unid) REFERENCES unihan_char (unid) ON UPDATE CASCADE ON DELETE CASCADE ,
FOREIGN KEY (unid_s) REFERENCES unihan_char (unid) ON UPDATE CASCADE ON DELETE CASCADE 
);
-- COMMENT '繁体字映射';

-- 变体字关系，
CREATE TABLE variant_id ( rid INTEGER PRIMARY KEY AUTOINCREMENT,
unid INT (32) NOT NULL,
-- COMMENT 'Unicode' ,
unid_s INT (32) NOT NULL,
-- COMMENT '最简字Unicode' ,
FOREIGN KEY (unid) REFERENCES unihan_char (unid) ON UPDATE CASCADE ON DELETE CASCADE ,
FOREIGN KEY (unid_s) REFERENCES unihan_char (unid) ON UPDATE CASCADE ON DELETE CASCADE 
);
-- COMMENT '变体字映射';

-- 包含“古同”
CREATE TABLE variant_gu ( unid INT (32) PRIMARY KEY,
-- COMMENT 'Unicode' ,
ch VARCHAR (1) NOT NULL UNIQUE,
-- COMMENT '汉字' ,
FOREIGN KEY (unid) REFERENCES unihan_char (unid) ON UPDATE CASCADE ON DELETE CASCADE 
) WITHOUT ROWID;
-- COMMENT '古异体字';

-- 包含“古同”
CREATE TABLE variant_gu_id ( rid INTEGER PRIMARY KEY AUTOINCREMENT,
unid INT (32) NOT NULL,
-- COMMENT 'Unicode' ,
unid_s INT (32) NOT NULL,
-- COMMENT '最简字Unicode' ,
FOREIGN KEY (unid) REFERENCES unihan_char (unid) ON UPDATE CASCADE ON DELETE CASCADE ,
FOREIGN KEY (unid_s) REFERENCES unihan_char (unid) ON UPDATE CASCADE ON DELETE CASCADE 
);
-- COMMENT '古异体字映射';

CREATE TABLE bian_ma ( unid INT (32) PRIMARY KEY,
-- COMMENT 'Unicode',
wu_bi86 TEXT,
-- COMMENT '五笔86码',
wu_bi98 TEXT,
-- COMMENT '五笔98码',
cang_jie TEXT,
-- COMMENT '仓颉码',
si_jiao TEXT,
-- COMMENT '四角码',
gui_fan TEXT,
-- COMMENT '规范汉字码',
FOREIGN KEY (unid) REFERENCES unihan_char (unid) ON UPDATE CASCADE ON DELETE CASCADE 
) WITHOUT ROWID;
-- COMMENT '常用编码表';

CREATE TABLE min_su ( msid INTEGER PRIMARY KEY AUTOINCREMENT,
is_surname VARCHAR (1),
-- COMMENT '姓名学（姓氏）,~bool' ,
surname_gender VARCHAR (1),
-- COMMENT '姓氏性别' ,
wu_xing VARCHAR (1),
-- COMMENT '五行' ,
lucky VARCHAR (4),
-- COMMENT '幸运' ,
regular VARCHAR (1)
-- COMMENT '常用,~bool' 
);
-- COMMENT '民俗';

CREATE TABLE min_su_id ( rid INTEGER PRIMARY KEY AUTOINCREMENT,
msid INT (11) NOT NULL,
-- COMMENT '民俗id',
unid INT (32) NOT NULL,
-- COMMENT 'Unicode',
FOREIGN KEY (msid) REFERENCES min_su (msid) ON UPDATE CASCADE ON DELETE CASCADE,
FOREIGN KEY (unid) REFERENCES unihan_char (unid) ON UPDATE CASCADE ON DELETE CASCADE 
);
-- COMMENT '民俗关联';

CREATE TABLE glyph ( unid INT (32) PRIMARY KEY,
-- COMMENT 'Unicode',
as_radical VARCHAR (1),
-- COMMENT '用作偏旁,~bool' ,
radical VARCHAR (1),
-- COMMENT '部首' ,
radical_stroke TINYINT,
-- COMMENT '部首笔画' ,
stroke INTEGER,
-- COMMENT '总笔画' ,
simplified_radical VARCHAR (1),
-- COMMENT '简体部首' ,
simplified_radical_stroke TINYINT,
-- COMMENT '简体部首笔画' ,
simplified_total_stroke INTEGER,
-- COMMENT '简体总笔画' ,
traditional_radical VARCHAR (1),
-- COMMENT '繁体部首' ,
traditional_radical_stroke TINYINT,
-- COMMENT '繁体部首笔画' ,
traditional_total_stroke INTEGER,
-- COMMENT '总笔画' ,
shou_wei TEXT NOT NULL,
-- COMMENT '首尾分解查字',
bu_jian TEXT,
-- COMMENT '汉字部件构造',
bi_hao TEXT,
-- COMMENT '笔顺编号',
bi_du TEXT,
-- COMMENT '笔顺读写',
FOREIGN KEY (unid) REFERENCES unihan_char (unid) ON UPDATE CASCADE ON DELETE CASCADE 
) WITHOUT ROWID;
-- COMMENT '字形结构';

-- 对于新造简体字需要取对应的康熙笔画
CREATE TABLE science_stroke ( unid INTEGER PRIMARY KEY,
-- COMMENT 'Unicode',
science_stroke INTEGER DEFAULT (0),
-- COMMENT '姓名学笔画数',
FOREIGN KEY (unid) REFERENCES unihan_char (unid) ON UPDATE CASCADE ON DELETE CASCADE 
) WITHOUT ROWID;
-- COMMENT '姓名学笔画数';

CREATE TABLE yin_yun ( unid INT (32) PRIMARY KEY,
-- COMMENT 'Unicode',
shang_gu TEXT,
-- COMMENT '上古音',
guang_yun TEXT,
-- COMMENT '广韵',
ping_shui TEXT,
-- COMMENT '平水韵',
FOREIGN KEY (unid) REFERENCES unihan_char (unid) ON UPDATE CASCADE ON DELETE CASCADE 
) WITHOUT ROWID;
-- COMMENT '音韵参考';

CREATE TABLE pin_yin ( pid INTEGER PRIMARY KEY AUTOINCREMENT,
pin_yin TEXT NOT NULL UNIQUE
-- COMMENT '拼音' 
);
-- COMMENT '拼音';

CREATE TABLE pin_yin_id ( rid INTEGER PRIMARY KEY AUTOINCREMENT,
pid INT (32) NOT NULL,
-- COMMENT '拼音id',
unid INT (32) NOT NULL,
-- COMMENT 'Unicode',
FOREIGN KEY (pid) REFERENCES pin_yin (pid) ON UPDATE CASCADE ON DELETE CASCADE,
FOREIGN KEY (unid) REFERENCES unihan_char (unid) ON UPDATE CASCADE ON DELETE CASCADE 
);
-- COMMENT '拼音关联';

CREATE TABLE guo_yu_id ( rid INTEGER PRIMARY KEY AUTOINCREMENT,
pid INT (32) NOT NULL,
-- COMMENT '拼音id',
unid INT (32) NOT NULL,
-- COMMENT 'Unicode',
FOREIGN KEY (pid) REFERENCES pin_yin (pid) ON UPDATE CASCADE ON DELETE CASCADE,
FOREIGN KEY (unid) REFERENCES unihan_char (unid) ON UPDATE CASCADE ON DELETE CASCADE 
);
-- COMMENT '国语拼音关联';

CREATE TABLE zhu_yin ( zid INTEGER PRIMARY KEY AUTOINCREMENT,
zhu_yin TEXT NOT NULL UNIQUE
-- COMMENT '注音' 
);
-- COMMENT '注音';

CREATE TABLE zhu_yin_id ( rid INTEGER PRIMARY KEY AUTOINCREMENT,
zid INT (32) NOT NULL,
-- COMMENT '注音id',
unid INT (32) NOT NULL,
-- COMMENT 'Unicode',
FOREIGN KEY (zid) REFERENCES zhu_yin (zid) ON UPDATE CASCADE ON DELETE CASCADE,
FOREIGN KEY (unid) REFERENCES unihan_char (unid) ON UPDATE CASCADE ON DELETE CASCADE 
);
-- COMMENT '注音关联';

CREATE TABLE tang_yin ( tid INTEGER PRIMARY KEY AUTOINCREMENT,
tang_yin TEXT NOT NULL UNIQUE
-- COMMENT '唐音' 
);
-- COMMENT '唐音';

CREATE TABLE tang_yin_id ( rid INTEGER PRIMARY KEY AUTOINCREMENT,
tid INT (32) NOT NULL,
-- COMMENT '唐音id',
unid INT (32) NOT NULL,
-- COMMENT 'Unicode',
FOREIGN KEY (tid) REFERENCES tang_yin (tid) ON UPDATE CASCADE ON DELETE CASCADE,
FOREIGN KEY (unid) REFERENCES unihan_char (unid) ON UPDATE CASCADE ON DELETE CASCADE 
);
-- COMMENT '唐音关联';

CREATE TABLE yue_yin ( yid INTEGER PRIMARY KEY AUTOINCREMENT,
yue_yin TEXT NOT NULL UNIQUE
-- COMMENT '粤音' 
);
-- COMMENT '粤音';

CREATE TABLE yue_yin_id ( rid INTEGER PRIMARY KEY AUTOINCREMENT,
yid INT (32) NOT NULL,
-- COMMENT '粤音id',
unid INT (32) NOT NULL,
-- COMMENT 'Unicode',
FOREIGN KEY (yid) REFERENCES yue_yin (yid) ON UPDATE CASCADE ON DELETE CASCADE,
FOREIGN KEY (unid) REFERENCES unihan_char (unid) ON UPDATE CASCADE ON DELETE CASCADE 
);
-- COMMENT '粤音关联';

CREATE TABLE min_nan_yin ( mnid INTEGER PRIMARY KEY AUTOINCREMENT,
min_nan_yin TEXT NOT NULL UNIQUE
-- COMMENT '闽南音' 
);
-- COMMENT '闽南音';

CREATE TABLE min_nan_yin_id ( rid INTEGER PRIMARY KEY AUTOINCREMENT,
mnid INT (32) NOT NULL,
-- COMMENT 'Unicode',
unid INT (32) NOT NULL,
-- COMMENT '成语id',
FOREIGN KEY (mnid) REFERENCES min_nan_yin (mnid) ON UPDATE CASCADE ON DELETE CASCADE,
FOREIGN KEY (unid) REFERENCES unihan_char (unid) ON UPDATE CASCADE ON DELETE CASCADE 
);
-- COMMENT '闽南关联';

CREATE TABLE suo_yin ( unid INT (32) PRIMARY KEY,
-- COMMENT 'Unicode',
gu TEXT,
-- COMMENT '古文字诂林',
xun TEXT,
-- COMMENT '故训彙纂',
shuo TEXT,
-- COMMENT '说文解字',
kang TEXT,
-- COMMENT '康熙字典',
han TEXT,
-- COMMENT '汉语字典',
ci TEXT,
-- COMMENT '辞海',
FOREIGN KEY (unid) REFERENCES unihan_char (unid) ON UPDATE CASCADE ON DELETE CASCADE 
) WITHOUT ROWID;
-- COMMENT '索引参考';

CREATE TABLE xin_hua ( unid INT (32) PRIMARY KEY,
-- COMMENT 'Unicode',
content TEXT,
-- COMMENT '新华字典',
FOREIGN KEY (unid) REFERENCES unihan_char (unid) ON UPDATE CASCADE ON DELETE CASCADE 
) WITHOUT ROWID;
-- COMMENT '详细解释';

CREATE TABLE han_da ( unid INT (32) PRIMARY KEY,
-- COMMENT 'Unicode',
content TEXT,
-- COMMENT '汉语大字典',
FOREIGN KEY (unid) REFERENCES unihan_char (unid) ON UPDATE CASCADE ON DELETE CASCADE 
) WITHOUT ROWID;
-- COMMENT '汉语字典';

CREATE TABLE kang_xi ( unid INT (32) PRIMARY KEY,
-- COMMENT 'Unicode',
stroke INTEGER,
-- COMMENT '康熙笔画',
brief TEXT,
-- COMMENT '简要',
content TEXT,
-- COMMENT '解释',
url TEXT,
-- COMMENT '扫描资料链接',
FOREIGN KEY (unid) REFERENCES unihan_char (unid) ON UPDATE CASCADE ON DELETE CASCADE 
) WITHOUT ROWID;
-- COMMENT '康熙字典';

CREATE TABLE shuo_wen ( unid INT (32) PRIMARY KEY,
-- COMMENT 'Unicode',
brief TEXT,
-- COMMENT '简要',
content TEXT,
-- COMMENT '详解',
url TEXT,
-- COMMENT '图像链接',
FOREIGN KEY (unid) REFERENCES unihan_char (unid) ON UPDATE CASCADE ON DELETE CASCADE 
) WITHOUT ROWID;
-- COMMENT '说文解字';

CREATE TABLE yan_bian ( unid INT (32) PRIMARY KEY,
-- COMMENT 'Unicode',
jia_gu_url TEXT,
-- COMMENT '甲骨文图像链接',
jin_url TEXT,
-- COMMENT '金文图像链接',
xiao_zhuan_url TEXT,
-- COMMENT '小篆图像链接',
kai_url TEXT,
-- COMMENT '楷体图像链接',
FOREIGN KEY (unid) REFERENCES unihan_char (unid) ON UPDATE CASCADE ON DELETE CASCADE 
) WITHOUT ROWID;
-- COMMENT '字源演变';

CREATE TABLE cheng_yu ( cyid INTEGER PRIMARY KEY AUTOINCREMENT,
cheng_yu TEXT NOT NULL UNIQUE,
-- COMMENT '成语' ,
url TEXT NOT NULL UNIQUE
-- COMMENT '链接' 
);
-- COMMENT '成语';

CREATE TABLE cheng_yu_id ( rid INTEGER PRIMARY KEY AUTOINCREMENT,
cyid INT (32) NOT NULL,
-- COMMENT '成语id',
unid INT (32) NOT NULL,
-- COMMENT 'Unicode',
FOREIGN KEY (cyid) REFERENCES cheng_yu (cyid) ON UPDATE CASCADE ON DELETE CASCADE ,
FOREIGN KEY (unid) REFERENCES unihan_char (unid) ON UPDATE CASCADE ON DELETE CASCADE 
);
-- COMMENT '相关成语';

CREATE TABLE shi_ci ( scid INTEGER PRIMARY KEY AUTOINCREMENT,
shi_ci TEXT NOT NULL UNIQUE,
-- COMMENT '诗词' ,
url TEXT NOT NULL UNIQUE
-- COMMENT '链接' 
);
-- COMMENT '诗词';

CREATE TABLE shi_ci_id ( rid INTEGER PRIMARY KEY AUTOINCREMENT,
scid INT (32) NOT NULL,
-- COMMENT '诗词id',
unid INT (32) NOT NULL,
-- COMMENT 'Unicode',
FOREIGN KEY (scid) REFERENCES shi_ci (scid) ON UPDATE CASCADE ON DELETE CASCADE ,
FOREIGN KEY (unid) REFERENCES unihan_char (unid) ON UPDATE CASCADE ON DELETE CASCADE
);
-- COMMENT '相关诗词';

CREATE TABLE ci_yu ( cid INTEGER PRIMARY KEY AUTOINCREMENT,
ci_yu TEXT NOT NULL UNIQUE,
-- COMMENT '词语' ,
url TEXT NOT NULL UNIQUE
-- COMMENT '链接' 
);
-- COMMENT '词语';

CREATE TABLE ci_yu_id ( rid INTEGER PRIMARY KEY AUTOINCREMENT ,
cid INT (32) NOT NULL,
-- COMMENT '词语id',
unid INT (32) NOT NULL,
-- COMMENT 'Unicode',
FOREIGN KEY (cid) REFERENCES ci_yu (cid) ON UPDATE CASCADE ON DELETE CASCADE ,
FOREIGN KEY (unid) REFERENCES unihan_char (unid) ON UPDATE CASCADE ON DELETE CASCADE 
);
-- COMMENT '相关词语';

CREATE TABLE han_char ( unid INT (32) PRIMARY KEY,
-- COMMENT 'Unicode' ,
ch VARCHAR (1) NOT NULL UNIQUE,
-- COMMENT '汉字' ,
FOREIGN KEY (unid) REFERENCES unihan_char (unid) ON UPDATE CASCADE ON DELETE CASCADE 
) WITHOUT ROWID;
-- COMMENT '最简汉字表' ;
