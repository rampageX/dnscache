package main

import (
	"fmt"
	"net"
	"time"

	"github.com/golang/groupcache/lru"
	"github.com/miekg/dns"
	"gopkg.in/fatih/pool.v2"
)

const (
	notIPQuery = 0
	_IP4Query  = 4
	_IP6Query  = 6
)

type Question struct {
	qname  string
	qtype  string
	qclass string
}

func (q *Question) String() string {
	return q.qname + " " + q.qclass + " " + q.qtype
}

type GODNSHandler struct {
	resolver *Resolver
	Cache    *lru.Cache
}

func NewHandler() *GODNSHandler {

	var (
		resolver *Resolver
		Cache    *lru.Cache
	)
	resolver = &Resolver{}
	Cache = lru.New(MAX_CACHES)
	return &GODNSHandler{resolver, Cache}
}

func (h *GODNSHandler) GetHour() string {
	return time.Now().Format("2006010215")
}

func (h *GODNSHandler) DoInitPool(nsaddr string) {
	//fmt.Println("try to connect to ", nsaddr)
	p, err := pool.NewChannelPool(1, 10, func() (net.Conn, error) { return net.Dial("tcp", nsaddr) })
	if err == nil {
		h.resolver.NameserversPool = append(h.resolver.NameserversPool, p)
	}
}
func (h *GODNSHandler) PreparePool() {
	for _, nsaddr := range NSADDRS {
		go h.DoInitPool(nsaddr)
	}
}

func (h *GODNSHandler) do(Net string, w dns.ResponseWriter, req *dns.Msg) {
	q := req.Question[0]
	Q := Question{UnFqdn(q.Name), dns.TypeToString[q.Qtype], dns.ClassToString[q.Qclass]}

	//fmt.Println("DNS Lookup ", Q.String())

	IPQuery := h.isIPQuery(q)
	key := fmt.Sprintf("%s-%s", h.GetHour(), Q.String())
	//fmt.Println("Cache key: ", key)
	if IPQuery > 0 {
		mesg, ok := h.Cache.Get(key)
		if ok == true {
			//fmt.Println("Hit cache", Q.String())
			rmesg := mesg.(*dns.Msg)
			rmesg.Id = req.Id
			w.WriteMsg(BuildDNSMsg(rmesg))
			return
		}
	}

	mesg, err := h.resolver.Lookup(Net, req)

	if err != nil {
		mesg, err = h.resolver.Lookup(Net, req) // try to lookup again
		if err != nil {
			//		fmt.Println("Resolve query error ", err)
			dns.HandleFailed(w, req)
		}
		return
	}

	w.WriteMsg(BuildDNSMsg(mesg))

	if IPQuery > 0 && len(mesg.Answer) > 0 {
		h.Cache.Add(key, mesg)
		//	fmt.Println("Insert into cache", Q.String())
	}
}

func BuildDNSMsg(msg *dns.Msg) *dns.Msg {
	msg.Compress = true
	//fmt.Println(msg)
	return msg
}
func (h *GODNSHandler) DoTCP(w dns.ResponseWriter, req *dns.Msg) {
	h.do("tcp", w, req)
}

func (h *GODNSHandler) DoUDP(w dns.ResponseWriter, req *dns.Msg) {
	h.do("udp", w, req)
}

func (h *GODNSHandler) isIPQuery(q dns.Question) int {
	if q.Qclass != dns.ClassINET {
		return notIPQuery
	}

	switch q.Qtype {
	case dns.TypeA:
		return _IP4Query
	case dns.TypeAAAA:
		return _IP6Query
	default:
		return notIPQuery
	}
}

func UnFqdn(s string) string {
	if dns.IsFqdn(s) {
		return s[:len(s)-1]
	}
	return s
}
