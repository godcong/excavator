package excavator

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/browser"
	"time"
)

func StartService() {
	engine := gin.Default()

	//todo:add handle

	go func() {
		<-time.After(5 * time.Second)
		err := browser.OpenURL("http://localhost:8080")
		if err != nil {
			Log.Error(err)
		}
	}()

	if err := engine.Run(); err != nil {
		Log.Fatal(err)
	}
}
