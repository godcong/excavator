
CREATE TABLE character ( ch VARCHAR (1) PRIMARY KEY,
-- COMMENT '汉字' ,
is_regular BOOLEAN NOT NULL，
-- COMMENT '常用' ,
pin_yin TEXT NOT NULL,
-- COMMENT '拼音' ,
is_duo_yin BOOLEAN NOT NULL,
-- COMMENT '多音字' ,
is_surname BOOLEAN,
-- COMMENT '姓氏' ,
surname_gender VARCHAR (1),
-- COMMENT '性别' ,
wu_xing VARCHAR (1) NOT NULL,
-- COMMENT '五行' ,
lucky TEXT,
-- COMMENT '吉凶寓意' ,
radical VARCHAR (1) NOT NULL,
-- COMMENT '部首' ,
stroke INTEGER DEFAULT (0) ,
-- COMMENT '笔画数' ,
science_stroke INTEGER DEFAULT (0)
-- COMMENT '姓名学笔画数'
) WITHOUT ROWID;
-- COMMENT '取名字表' ;

