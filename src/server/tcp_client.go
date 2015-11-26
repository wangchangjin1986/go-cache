// tcp_client
package main

import (
	"bufio"
	"fmt"
	"net"
	"proto"
	//	"time"
)

type TcpClient struct {
	conn     net.Conn
	sendChan chan uint8
	req      *proto.MCRequest
	tcps     *TcpServer
}

func (tc *TcpClient) conn_loop() {
	if tc.tcps == nil {
		return
	}
	defer func() {
		if err := recover(); err != nil {
			mylog.Println("lost one client connect, error:%v", err)
		}
	}()
	mylog.Println("client %v connected", tc.conn.RemoteAddr().String())
	for {
		tc.req = &proto.MCRequest{}
		if err := tc.req.Receive(bufio.NewReader(tc.conn)); err != nil {
			goto exit
		}
		fmt.Printf("receive data:%v", tc.req)
		//time.Sleep(5 * time.Second)
	}
exit:
	mylog.Println("lost one client connect")
}

func NewTcpClient(cn net.Conn, s *TcpServer) *TcpClient {
	if s == nil {
		return nil
	}
	return &TcpClient{conn: cn, tcps: s}
}
