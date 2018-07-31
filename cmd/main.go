package main

import (
"fmt"


"github.com/godcong/excavator"
"github.com/godcong/excavator/db"
)

func main() {
	chars := excavator.CommonlyTop("http://www.zdic.net/z/zb/cc1.htm")
	fmt.Println("start total:", len(chars))
	for idx, v := range chars {
		bc := excavator.CommonlyBase("http://www.zdic.net", v)
		db.DB("base").Insert(bc)
		fmt.Println("current is :", idx, v.Character)
	}
	fmt.Println("end")
}
