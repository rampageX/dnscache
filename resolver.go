package main

import (
	"fmt"
	//"sync"
	"time"

	"github.com/miekg/dns"
	"github.com/pkg/errors"
	"gopkg.in/fatih/pool.v2"
)

//Resolver the resolver struct
type Resolver struct {
	NameserversPool []pool.Pool
}

//Lookup do Lookup with the resolver
func (r *Resolver) Lookup(net string, req *dns.Msg) (message *dns.Msg, err error) {
	qname := req.Question[0].Name

	res := make(chan *dns.Msg, 1)

	L := func(nsPool pool.Pool) {

		for i := 0; i < 10; i++ {
			//r, rtt, err := c.Exchange(req, nameserver)
			fmt.Println(fmt.Sprintln("The conn pool size: ", nsPool.Len()))

			c, err := nsPool.Get()
			defer c.Close()
			if err != nil {
				if pc, ok := c.(*pool.PoolConn); ok == true {
					pc.MarkUnusable()
					pc.Close()
				}
				return
			}

			co := &dns.Conn{Conn: c.(*pool.PoolConn).Conn} // c is your net.Conn

			_ = co.WriteMsg(req) // No handler write error
			r, err := co.ReadMsg()
			defer co.Close()

			if err != nil {
				fmt.Println("socket error on ", qname)
				fmt.Println("error:", err.Error())

				if pc, ok := c.(*pool.PoolConn); ok {
					pc.MarkUnusable()
					pc.Close()
				}
			}
			// If SERVFAIL happen, should return immediately and try another upstream resolver.
			// However, other Error code like NXDOMAIN is an clear response stating
			// that it has been verified no such domain existas and ask other resolvers
			// would make no sense. See more about #20
			if r == nil {
				fmt.Println("Got nil response")
			}
			if r != nil && r.Rcode != dns.RcodeSuccess {
				fmt.Println("Failed to get an valid answer ", qname)
			}
			if r != nil && r.Rcode == dns.RcodeSuccess {
				fmt.Println("resolv ", UnFqdn(qname), " on ", r.String(), r.Len())
				res <- r
				break
			}
		}

	}
	// Start lookup on each nameserver top-down, in every second
	for _, nspool := range r.NameserversPool {
		go L(nspool)
	}
	timeout := time.After(time.Second * Timeout)
	select {
	case r := <-res:
		if r == nil {
			return nil, errors.New(" failed to got valid answer")
		}
		return r, nil
	case <-timeout:
		return nil, errors.New("Time out on dns query")
	}
}

//Timeout set the timeout
func (r *Resolver) Timeout() time.Duration {
	return time.Duration(Timeout) * time.Second
}
