package ipdb

import (
	"log"
	"testing"
)

//
func TestSpider(t *testing.T) {
	iprs := DownloadIpRangsByIcpn()
	SaveIprs2File("../iprs_ipcn.dat", iprs)
}

//
func TestLoad(t *testing.T) {
	LoadDbByFile("../iprs_ipcn.dat")
	p, c, err := FindIpInRange("1.192.62.67")
	log.Println(p, c, err)
}
