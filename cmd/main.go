package main

import (
	"fmt"
	"net/http"
	"os"
	"sort"

	"github.com/godcong/excavator"
	"github.com/urfave/cli"
)

func main() {
	app := cli.App{
		Version: "v0.0.1",
		Name:    "excavator",
		Usage:   "excavator a dictionary",
		Action: func(c *cli.Context) error {
			url := ""
			if c.NArg() > 0 {
				url = c.Args().Get(0)
			}
			exc := excavator.New(url, "")
			header := make(http.Header)
			header.Set("Cookie", "hy_so_4=%255B%257B%2522zi%2522%253A%2522%25E8%2592%258B%2522%252C%2522url%2522%253A%252234%252FKOKORNKOCQXVILXVB%252F%2522%252C%2522py%2522%253A%2522ji%25C7%258Eng%252C%2522%252C%2522bushou%2522%253A%2522%25E8%2589%25B9%2522%252C%2522num%2522%253A%252217%2522%257D%255D; ASP.NET_SessionId=zilmx52mwtr3xsq5i212pd5a; UM_distinctid=16c2efb5e9e134-0cfc801ee6ae06-353166-1fa400-16c2efb5e9f3c8; CNZZDATA1267010321=1299014713-1564151968-%7C1564151968; Hm_lvt_cd7ed86134e0e3138a7cf1994e6966c8=1564156322; Hm_lpvt_cd7ed86134e0e3138a7cf1994e6966c8=1564156322")
			header.Set("Origin", "http://hy.httpcn.com")
			header.Set("Accept-Encoding", "gzip, deflate")
			header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,ja;q=0.7,zh-TW;q=0.6")
			header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.142 Mobile Safari/537.36")
			header.Set("Content-Type", "application/x-www-form-urlencoded")
			header.Set("Accept", "application/json")
			header.Set("Referer", "http://hy.httpcn.com/bushou/kangxi/")
			header.Set("X-Requested-With", "XMLHttpRequest")
			header.Set("Connection", "keep-alive")
			exc.SetHeader(header)

			exc.SetStep(excavator.Step(c.String("step")))

			exc.Run()

			r := exc.Radical()
			cr := exc.Character()
			if r != nil {
				for {
					select {
					case rr := <-r:
						if rr == nil {
							goto END
						}
					}
				}
			}

			if cr != nil {
				for {
					select {
					case ccr := <-cr:
						fmt.Println(ccr.Character, "inserted")
						if ccr == nil {
							goto END
						}
					}
				}
			}
		END:
			return nil
		},
		Flags: mainFlags(),
	}
	app.Commands = []cli.Command{}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))
	err := app.Run(os.Args)
	if err != nil {
		return
	}
}

func mainFlags() (flags []cli.Flag) {
	flags = []cli.Flag{
		cli.StringFlag{
			Name:  "workspace",
			Usage: "set workspace to storage temp file",
		},
		cli.StringFlag{
			Name:  "step",
			Usage: "set the run step",
		},
	}
	return flags
}
