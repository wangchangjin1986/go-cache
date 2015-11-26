// tcp_server
package main

import (
	"cache"
	"container/list"
	"errors"
	"net"
	"runtime"
	"util"
)

type TcpServer struct {
	tcp_listener net.Listener
	haddr        string
	waitGroup    util.WaitGroupWrapper
	exitChan     chan int
	clientlist   *list.List
}

func (cs *TcpServer) waitforconn(ln net.Listener, c *cache.Cache) {
	for {
		tconn, err := ln.Accept()
		if err != nil {
			if nerr, ok := err.(net.Error); ok && nerr.Temporary() {
				runtime.Gosched()
				continue
			} else {
				goto exit
			}
		}
		tc := NewTcpClient(tconn, cs)
		if cs.clientlist == nil {
			runtime.Gosched()
			continue
		}
		go tc.conn_loop()
		cs.clientlist.PushBack(tc)
	}
exit:
	mylog.Println("server waitforconn exit")
}

func (cs *TcpServer) Start(c *cache.Cache) (int, error) {
	if cs == nil || len(cs.haddr) < 0 {
		return -1, errors.New("server address is not set")
	}
	if tcplistener, err := net.Listen("tcp", cs.haddr); err == nil {
		cs.tcp_listener = tcplistener
	} else {
		return -1, errors.New("could not listen")
	}
	cs.waitGroup.Wrap(func() { cs.waitforconn(cs.tcp_listener, c) })
	return 0, nil
}
func (cs *TcpServer) clear() {
	cs.tcp_listener.Close()
	for e := cs.clientlist.Front(); e != nil; e = e.Next() {
		client := e.Value
		if _, found := client.(TcpClient); found {
			client.(TcpClient).conn.Close()
		}
	}
}
func (cs *TcpServer) Stop() {
	//cs.exitChan <- 1
	cs.clear()
	cs.waitGroup.Wait()
}
func NewTcpServer(addr string) *TcpServer {
	if len(addr) <= 0 {
		return &TcpServer{haddr: "0.0.0.0:9891", exitChan: make(chan int), clientlist: list.New()}
	} else {
		return &TcpServer{haddr: addr, exitChan: make(chan int), clientlist: list.New()}
	}
}
