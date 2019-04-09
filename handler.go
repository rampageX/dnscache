package main

import (
	"fmt"
	"net"
	"time"

	"github.com/miekg/dns"
	lru "github.com/netroby/fastlru"
	"gopkg.in/fatih/pool.v2"
)

const (
	notIPQuery = 0
	_IP4Query  = 4
	_IP6Query  = 6
)

// Question : dns query struct
type Question struct {
	qname  string
	qtype  string
	qclass string
}

func (q *Question) String() string {
	return q.qname + " " + q.qclass + " " + q.qtype
}

// GODNSHandler : dns handler struct
type GODNSHandler struct {
	resolver *Resolver
	Cache    *lru.Cache
}

// NewHandler : create new DNS handler
func NewHandler() *GODNSHandler {

	var (
		resolver *Resolver
		Cache    *lru.Cache
	)
	resolver = &Resolver{}
	Cache = lru.New(MaxCaches)
	return &GODNSHandler{resolver, Cache}
}

// GetHour : get hour format string
func (h *GODNSHandler) GetHour() string {
	return time.Now().Format("2006010215")
}

// DoInitPool : Do Initialize Pool
func (h *GODNSHandler) DoInitPool(nsaddr string) {
	fmt.Println("DoInitPool, try to connect to ", nsaddr)
	p, err := pool.NewChannelPool(5, 30, func() (net.Conn, error) {
		var d = net.Dialer{
			KeepAlive: time.Duration(Timeout),
		}
		return d.Dial("tcp", nsaddr)
	})
	if err == nil {
		h.resolver.NameserversPool = append(h.resolver.NameserversPool, p)
	}
	conn, err := p.Get()
	defer conn.Close()
	if err != nil {
		fmt.Println("error when get conn from pool")
		fmt.Println(err)
		if pc, ok := conn.(*pool.PoolConn); ok == true {
			pc.MarkUnusable()
			pc.Close()
		}
	}
	fmt.Println(fmt.Sprintln("the pool size: ", p.Len()))
}

// PreparePool : To prepare pool for use
func (h *GODNSHandler) PreparePool() {
	for _, nsaddr := range NsAddrs {
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

// BuildDNSMsg : to build dns msg for response
func BuildDNSMsg(msg *dns.Msg) *dns.Msg {
	msg.Compress = true
	//fmt.Println(msg)
	return msg
}

// DoTCP : do tcp dns resolve
func (h *GODNSHandler) DoTCP(w dns.ResponseWriter, req *dns.Msg) {
	h.do("tcp", w, req)
}

// DoUDP : do udp dns resolve
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

// UnFqdn : process and return dns msg
func UnFqdn(s string) string {
	if dns.IsFqdn(s) {
		return s[:len(s)-1]
	}
	return s
}
