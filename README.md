#  DNS Cache

This is a DNS Cache server for local network. For MacOSX, Linux User. 
Not test for Windows

You need golang 1.5+ to compile dnscache

## Features

1. Tcp conn to upstream
2. Cache dns query 
3. parallel request upstream got dns result
4. easy to use (at least on linux/macosx)

## Install


### Fast 

```
go get -u -v github.com/netroby/dnscache
```

Then dnscache will be your go/bin dir. like $HOME/go/bin/dnscache

You now can run with it 

```
$HOME/go/bin/dnscache 0.0.0.0 53
```

try to modify your /etc/resolv.conf or using drill dig to check your dns
cache

```
[huzhifeng@localhost ~]$ drill www.google.com @127.0.0.1
;; ->>HEADER<<- opcode: QUERY, rcode: NOERROR, id: 55816
;; flags: qr rd ra ; QUERY: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 0 
;; QUESTION SECTION:
;; www.google.com.	IN	A

;; ANSWER SECTION:
www.google.com.	197	IN	A	216.58.221.164

;; AUTHORITY SECTION:

;; ADDITIONAL SECTION:

;; Query time: 0 msec
;; SERVER: 127.0.0.1
;; WHEN: Fri Nov  4 17:47:10 2016
;; MSG SIZE  rcvd: 48

```

### Build your own

```
go get -u -v github.com/tools/godep
go get ./...
godep restore
godep go build
sudo ./dnscache
```
It will listen on 127.0.0.1:53 , both UDP/TCP port


## Donate me please

### Bitcoin donate

```
136MYemy5QmmBPLBLr1GHZfkES7CsoG4Qh
```
### Alipay donate
![Scan QRCode donate me via Alipay](https://www.netroby.com/assets/images/alipayme.jpg)

**Scan QRCode donate me via Alipay**
