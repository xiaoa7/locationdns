/**
 *使用SQLite memory 做数据库，查询请求ip归属
 */
package ipdb

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var mdb *sql.DB

//从文件加载内存数据库
func LoadDbByFile(filename string) {
	var err error
	mdb, err = sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	//建表
	_, err = mdb.Exec("create table iprs (id integer not null primary key, province text,carrieroperator text,start integer,end integer,comeintime integer);")
	if err != nil {
		log.Fatal(err)
	}
	//建索引
	_, err = mdb.Exec("create index iprsindex on iprs(start,end,comeintime);")
	if err != nil {
		log.Fatal(err)
	}
	//加载文件
	fi, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	bs, _ := ioutil.ReadAll(fi)
	fi.Close()
	iprs := make([]*IPRange, 0, 0)
	json.Unmarshal(bs, &iprs)
	//插入内存数据库
	tx, err := mdb.Begin()
	if err != nil {
		log.Fatal(err)
	}
	comeintime := time.Now().Unix()
	stmt, err := tx.Prepare("insert into iprs(province,carrieroperator,start,end,comeintime) values(?,?,?,?,?)")
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range iprs {
		_, err = stmt.Exec(v.Province, v.Carrieroperator, v.Start, v.End, comeintime)
		if err != nil {
			log.Fatal(err)
		}
	}
	tx.Commit()
	stmt.Close()
}

//根据ip查询省份，运营商
func FindIpInRange(ip string) (p string, c string, err error) {
	ipint := ip2int(ip)
	rs, err := mdb.Query("select province,carrieroperator from iprs where start<? and end >? order by comeintime desc limit 1", ipint, ipint)
	if err != nil {
		return "", "", err
	}
	defer rs.Close()
	if rs.Next() {
		err = rs.Scan(&p, &c)
	}
	return
}

//刷新数据库，以comeintime为准，不停服务
func UpdateDb(iprs []*IPRange) {
	if len(iprs) < 200 {
		return
	}
	//插入内存数据库
	tx, err := mdb.Begin()
	if err != nil {
		log.Fatal(err)
	}
	comeintime := time.Now().Unix()
	stmt, err := tx.Prepare("insert into iprs(province,carrieroperator,start,end,comeintime) values(?,?,?,?,?)")
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range iprs {
		_, err = stmt.Exec(v.Province, v.Carrieroperator, v.Start, v.End, comeintime)
		if err != nil {
			log.Fatal(err)
		}
	}
	tx.Commit()
	stmt.Close()
	//插入成功后，删除老数据（查询时，按comeintime倒排，也不会查老数据）
	mdb.Exec("delete from iprs where comeintime < ?", comeintime)
}

//释放内存数据库
func DestoryDb() {
	mdb.Close()
}
