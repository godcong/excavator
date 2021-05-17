package excavator

import (
	"embed"
)

//下面的注释不能删，是用于向二进制文件中嵌入静态文件的
//go:embed data
var DataFiles embed.FS
