/**
从ipcn网站爬取ip归属
*/
package ipdb

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	gq "github.com/PuerkitoBio/goquery"
)

const (
	carrieroperator_url_prefix = "http://ipcn.chacuo.net/view/"
)

//ip段
type IPRange struct {
	Province        string `json:"province"`
	Carrieroperator string `json:"cop"`
	Start           int64  `json:"start"`
	End             int64  `json:"end"`
}

//运营商清单
var carrieroperator = []string{
	"i_UNICOM",
	"i_CHINANET",
	"i_CMNET",
	"i_CERNET",
	"i_CRCT",
	"i_CNCGROUP",
	"i_GWBN",
	"i_CSTN",
	"i_BCN",
	"i_GeHua",
	"i_Topway",
	"i_ZHONG-BANG-YA-TONG",
	"i_FOUNDERBN",
	"i_WASU",
	"i_GZPRBNET",
	"i_HTXX",
	"i_eTrunk",
	"i_WSN",
	"i_CHINAGBN",
	"i_EASTERNFIBERNET",
	"i_LiaoHe-HuaYu",
	"i_CTN",
}

//省份清理
var province_name_clear *regexp.Regexp

//
func init() {
	province_name_clear, _ = regexp.Compile("\\(.[^\\)]*\\)")
}

//从icpn网站下载运营商分省份的ip段
func DownloadIpRangsByIcpn() (iprs []*IPRange) {
	log.Println("开始下载icpn ip段数据")
	iprs = make([]*IPRange, 0, 0)
	for _, v := range carrieroperator {
		log.Println("下载", v+"的ip地址段")
		doc, err := gq.NewDocument(carrieroperator_url_prefix + v)
		if err != nil {
			fmt.Println("err:", err.Error())
			return
		}
		current_province := ""
		doc.Find("dl.list").Each(func(i int, q *gq.Selection) {
			q.Children().Each(func(i int, cq *gq.Selection) {
				if tag := strings.ToLower(cq.Nodes[0].Data); tag == "dt" {
					current_province = province_name_clear.ReplaceAllString(cq.Text(), "")
				} else if tag == "dd" {
					cqc := cq.Children()
					start, end := ip2int(cqc.First().Text()), ip2int(cqc.Next().Text())
					if start == 0 || end == 0 {
						html, err := cq.Html()
						fmt.Println("ip转uint64出错了", err, html)
					} else {
						ipr := IPRange{
							Province:        current_province,
							Carrieroperator: v,
							Start:           start,
							End:             end,
						}
						iprs = append(iprs, &ipr)
					}
				}
			})
		})
	}
	return
}

//
func ip2int(ip string) int64 {
	tmp := strings.Split(ip, ".")
	if len(tmp) != 4 {
		fmt.Println("长度不够")
		return 0
	}
	b0, _ := strconv.Atoi(tmp[0])
	b1, _ := strconv.Atoi(tmp[1])
	b2, _ := strconv.Atoi(tmp[2])
	b3, _ := strconv.Atoi(tmp[3])

	var sum int64

	sum += int64(b0) << 24
	sum += int64(b1) << 16
	sum += int64(b2) << 8
	sum += int64(b3)
	return sum
}

//保存文件
func SaveIprs2File(filename string, iprs []*IPRange) {
	bs, err := json.Marshal(iprs)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fi, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_SYNC|os.O_TRUNC, 0777)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer fi.Close()
	_, err = fi.Write(bs)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
