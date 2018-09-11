//dns 解析
package resolve

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
)

type Dns map[string]string

var dsn_resolve map[string]Dns

//从文件中加载
func LoadByFile(filename string) {
	//加载文件
	fi, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	bs, _ := ioutil.ReadAll(fi)
	fi.Close()
	dsn_resolve = make(map[string]Dns)
	if err = json.Unmarshal(bs, &dsn_resolve); err != nil {
		log.Fatal(err)
	}
}

//dns解析，暂时只根据省份，因为基本上所有的机房都是双线路或多线路
func ResolveDns(host, province string) (string, error) {
	if v, ok := dsn_resolve[host]; ok {
		if v1, ok := v[province]; ok {
			return v1, nil
		} else {
			return v["default"], nil
		}
	} else {
		return "0.0.0.0", errors.New("找不到host=" + host + "的配置")
	}
}
