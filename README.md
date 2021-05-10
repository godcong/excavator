# excavator

从ipfs下载：[数据包](https://ipfs.io/ipfs/QmZXYMee9TfSYsd7AUHfXa6vCuzr7fJVFBgMmJZ6czwT8Y?filename=excavator_unzip.exe)

此链接不能用浏览器打开，可以找一个ipfs工具（如ipfs-desktop），获取链接中的CID，检查文件，然后下载

推荐使用数据包，因为重新捕获数据对网站压力很大。

数据包在win10命名为exe可以自解压，linux需要使用带zstd的codec的7zip工具解压。

数据包内容：
```
.
└── Temp
    ├── GB.txt （Unicode列表）
    ├── config.json （配置文件）
    ├── exc.db （excavator的sqlite文件）
    ├── ft.db （fate的sqlite文件）
    └── tool.httpcn.com （页面缓存）
```

查看sqlite内容可以用sqlitestudio：

```powershell
choco install sqlitestudio -y
```

## How to use

由于使用sqlite，请确保PATH中能搜索到gcc

程序运行需要使config.json存储余当前路径。（把数据包放在NVME或者RAMDISK可以明显加速，实测40分钟可以完成）

```bash
go test -timeout 3600s -run ^TestExcavator_Run$ excavator -v
```

In path `regular/`, there is another patch tool for regular.

```bash
go test -timeout 30s -run ^TestNew$ excavator/regular -v
```

In path `strokefix/`, there is another patch tool for regular.

```bash
go test -timeout 30s -run ^TestNumberChar$ excavator/strokefix -v
```


### this tool used `exc.db` and `ft.db` as sqlite database storage.

### if you want to change the database

change the definition in `config.json`
when created them with sql file in path `data/`.
