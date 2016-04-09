# haoilib
好爱答题golang接口包装

#使用示例

```
package main

import (
	//"fmt"
	"github.com/niwho/haoi/haoilib"
	"io/ioutil"
	"log"
	"os"
)

func main() {

	haoilib.Hb.SetRebateW("1863|EAB7978B27325246")
	rt, err := haoilib.Hb.GetPointW("xxxxx|5Axxxxxxxxx7838B")
	if err != nil {
		log.Println("err:", err)
		return
	}
	log.Println(rt)
	busy, err := haoilib.Hb.GetBusyW()

	if err != nil {
		log.Println("err:", err)
		return
	}
	log.Println(busy)
	return
	f, err := os.Open(`C:\niwho\workspace\haoi\tter.571220160403205539.jpg`)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	fa, err := ioutil.ReadAll(f)
	log.Println(int(len(fa)))
	ret, rep, err := haoilib.Hb.SendByteExW(
		"adinfo|5AD075939B77838B",
		"3004",
		fa,
		int64(len(fa)),
		20,
		0,
		"golang",
	)
	log.Println(ret, rep, err)

}
```
