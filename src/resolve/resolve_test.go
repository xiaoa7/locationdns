package resolve

import (
	"log"
	"testing"
)

func TestResolve(t *testing.T) {
	LoadByFile("../config.json")
	ip, err := ResolveDns("www.baidu.com", "河南")
	log.Println(ip, err)
}
