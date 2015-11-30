// tcp_client
package main

import (
	"bufio"
	"cache"
	"net"
	"proto"
	"time"
)

type TcpClient struct {
	conn     net.Conn
	sendChan chan uint8
	cache    *cache.Cache
}

func (tc *TcpClient) conn_loop() {
	if tc.cache == nil {
		return
	}
	defer func() {
		if err := recover(); err != nil {
			mylog.Println("lost one client connect, error:%v", err)
		}
		tc.conn.Close()
	}()
	mylog.Println("client connected:", tc.conn.RemoteAddr().String())
	for {
		req := &proto.MCRequest{}
		if err := req.Receive(bufio.NewReader(tc.conn)); err != nil {
			goto exit
		}
		tc.req_process(req)
	}
exit:
	mylog.Println("lost one client connect")
}
func (tc *TcpClient) req_process(req *proto.MCRequest) {
	defer func() {
		if err := recover(); err != nil {
			mylog.Println(err)
			res := &proto.MCResponse{Status: proto.SERVER_ERROR}
			res.Transmit(tc.conn)
		}
	}()
	if req == nil {
		return
	}
	res := CmdFuncs[req.Opcode](req, tc.cache)
	if res.Timeout > 0 {
		time.Sleep(time.Duration(res.Timeout) * time.Millisecond)
	}
	if res == nil {
		return
	}
	if err := res.Transmit(tc.conn); err != nil {
		mylog.Println("response client error")
	}
}
func NewTcpClient(cn net.Conn, c *cache.Cache) *TcpClient {
	if c == nil {
		return nil
	}
	return &TcpClient{conn: cn, cache: c}
}
