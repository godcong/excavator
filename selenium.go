package excavator

import (
	"fmt"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

// Selenium ...
type Selenium struct {
	webDriver selenium.WebDriver
	path      string
	port      int
}

// NewSelenium ...
func NewSelenium(path string, port int) *Selenium {
	return &Selenium{
		path: path,
		port: port,
	}
}

// Start ...
func (s *Selenium) Start() {
	//opts := []selenium.ServiceOption{}
	caps := selenium.Capabilities{
		"browserName": "chrome",
	}

	// 禁止加载图片，加快渲染速度
	imagCaps := map[string]interface{}{
		"profile.managed_default_content_settings.images": 2,
	}

	chromeCaps := chrome.Capabilities{
		Prefs: imagCaps,
		Path:  "",
		Args: []string{
			"--headless", // 设置Chrome无头模式
			"--no-sandbox",
			"--user-agent=Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_2) AppleWebKit/604.4.7 (KHTML, like Gecko) Version/11.0.2 Safari/604.4.7", // 模拟user-agent，防反爬
		},
	}
	caps.AddChrome(chromeCaps)
	path := "chromedriver"
	if s.path != "" {
		path = s.path
	}
	// 启动chromedriver，端口号可自定义
	_, err := selenium.NewChromeDriverService(path, s.port)
	if err != nil {
		panic(fmt.Sprintf("Error starting the ChromeDriver server: %v", err))
	}
	// 调起chrome浏览器
	s.webDriver, err = selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", s.port))
	if err != nil {
		panic(err)
	}
}

// Get ...
func (s *Selenium) Get(url string) (w selenium.WebDriver, e error) {

	// 这是目标网站留下的坑，不加这个在linux系统中会显示手机网页，每个网站的策略不一样，需要区别处理。
	//webDriver.AddCookie(&selenium.Cookie{
	//	Name:  "defaultJumpDomain",
	//	Value: "www",
	//})
	// 导航到目标网站
	e = s.webDriver.Get(url)
	if e != nil {
		return nil, fmt.Errorf("failed to load page: %s", e)
	}
	return s.webDriver, nil
}
