package main

import (
	"flag"
	"ipdb"
	"log"
	"net"
	"os"
	"os/signal"
	"resolve"
	"syscall"

	"github.com/miekg/dns"
	"github.com/robfig/cron"
)

//初始化
func init() {
	resolve.LoadByFile("./config.json")
	ipdb.LoadDbByFile("./iprs_ipcn.dat")
}

//
type handler struct{}

//
func (h *handler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	msg := dns.Msg{}
	msg.SetReply(r)
	switch r.Question[0].Qtype {
	case dns.TypeA: //只处理A记录
		msg.Authoritative = true
		domain := msg.Question[0].Name //
		rip := w.RemoteAddr().String()
		if province, _, err := ipdb.FindIpInRange(rip); err == nil {
			if ip, err := resolve.ResolveDns(domain, province); err != nil {
				msg.Answer = append(msg.Answer, &dns.A{
					Hdr: dns.RR_Header{Name: domain, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60},
					A:   net.ParseIP(ip),
				})
			}
		}
	}
	w.WriteMsg(&msg)
}

//
func main() {
	//
	addr := flag.String("addr", ":53", "监听地址")
	nettype := flag.String("nettype", "udp", "协议")
	flushinterval := flag.String("interval", "@weekly", "ipcn 数据刷新间隔 cron表达式")
	flag.Parse()
	//更新IP段库
	c := cron.New()
	c.Start()
	defer c.Stop()
	c.AddFunc(*flushinterval, func() {
		ipdb.UpdateDb(ipdb.DownloadIpRangsByIcpn())
	})

	//
	srv := &dns.Server{Addr: *addr, Net: *nettype}
	srv.Handler = &handler{}
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		c.Stop()

		srv.Shutdown()
		ipdb.DestoryDb()
	}()
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Failed to set udp listener %s\n", err.Error())
	}
}
