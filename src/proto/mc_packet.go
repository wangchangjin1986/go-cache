package proto

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
)

//request command type
type CommandCode string

//response status type
type Status int

//status result map
var StatusRes map[Status]string

//response 状态
const (
	SUCCESS Status = 0

	STORED          Status = 1
	NOT_STORED      Status = 2
	END             Status = 3
	DELETED         Status = 4
	NOT_FOUND       Status = 5
	UNKNOWN_COMMAND Status = 6

	ERROR        Status = 101
	CLIENT_ERROR Status = 102
	SERVER_ERROR Status = 103
)

//request command
const (
	GET     CommandCode = "gets"
	SET     CommandCode = "set"
	ADD     CommandCode = "add"
	REPLACE CommandCode = "replace"
	DELETE  CommandCode = "delete"
	STATS   CommandCode = "stats"
	QUIT    CommandCode = "quit"
)

//status to string
func (s *Status) ToString() string {
	rv := StatusRes[*s]
	if rv == "" {
		rv = fmt.Sprintf("%s\r\n", StatusRes[NOT_FOUND])
	} else {
		rv = fmt.Sprintf("%s\r\n", rv)
	}
	return rv
}

//mc请求产生一个request对象
type MCRequest struct {
	//请求命令
	Opcode CommandCode
	//key
	Key string
	//请求内容
	Value []byte
	//请求标识
	Flags int32
	//请求内容长度
	Length int
	//过期时间
	Expires int64
}
type MCResponse struct {
	//命令
	Opcoed CommandCode
	//返回状态
	Status Status
	//key
	Key string
	//返回内容
	Value []byte
	//返回标识
	Flags int32
	//错误
	Fatal bool
	//延时(ms)
	Timeout int
	//丢弃
	NoReply bool
}

//request to string
func (req *MCRequest) String() string {
	return fmt.Sprintf("{MCRequest opcode=%s, key='%s', value=%d }",
		req.Opcode, req.Key, req.Value)
}

//将socket请求内容 解析为一个MCRequest对象
func (req *MCRequest) Receive(r *bufio.Reader) error {
	line, _, err := r.ReadLine()
	if err != nil {
		return err
	}
	var command CommandCode
	params := strings.Fields(string(line))
	if len(params) == 0 {
		command = ""
	} else {
		command = CommandCode(params[0])
	}
	switch command {

	case SET, ADD, REPLACE:
		req.Opcode = command
		if len(params) != 5 {
			req.Key = ""
			req.Length = 0
		} else {
			req.Key = params[1]
			req.Length, _ = strconv.Atoi(params[4])

			value := make([]byte, req.Length+2)
			io.ReadFull(r, value)

			req.Value = make([]byte, req.Length)
			copy(req.Value, value)
		}
	case GET:
		req.Opcode = command
		req.Key = params[1]
		// RunStats["cmd_get"].(*CounterStat).Increment(1)
	case STATS:
		req.Opcode = command
		req.Key = ""
	case DELETE:
		req.Opcode = command
		req.Key = params[1]
	case QUIT:
		req.Opcode = command
	}
	return err
}

//解析response 并把返回结果写入socket链接
func (res *MCResponse) Transmit(w net.Conn) (err error) {
	switch res.Status {
	case ERROR, CLIENT_ERROR, SERVER_ERROR:
		_, err = w.Write([]byte(res.Status.ToString()))
	default:
		switch res.Opcoed {
		case STATS:
			_, err = w.Write(res.Value)
		case GET:
			if res.Status == SUCCESS {
				rs := fmt.Sprintf("VALUE %s %d %d\r\n%s\r\nEND\r\n", res.Key, res.Flags, len(res.Value), res.Value)
				if _, err = w.Write([]byte(rs)); err != nil {
					fmt.Println("response error")
					return err
				} else {
					fmt.Println("response success")
				}
			} else {
				_, err = w.Write([]byte(res.Status.ToString()))
			}
		case SET, REPLACE, ADD:
			_, err = w.Write([]byte(res.Status.ToString()))
		case DELETE:
			_, err = w.Write([]byte(res.Status.ToString()))
		}
	}
	return
}
func (res *MCResponse) GenerateRes(req *MCRequest) {

}
func NewResStatus(opcoed CommandCode, status Status) *MCResponse {
	return &MCResponse{Opcoed: opcoed, Status: status}
}
func NewResFull(opcoed CommandCode, status Status, key string, flags int32, fatal bool) *MCResponse {
	return &MCResponse{Opcoed: opcoed, Status: status, Key: key,
		Flags: flags, Fatal: fatal}
}

//初始化response状态的返回结果
func init() {
	StatusRes = make(map[Status]string)
	StatusRes[ERROR] = "ERROR"
	StatusRes[CLIENT_ERROR] = "CLIENT_ERROR client_error"
	StatusRes[SERVER_ERROR] = "SERVER_ERROR server_error"
	StatusRes[STORED] = "STORED"
	StatusRes[NOT_STORED] = "NOT_STORED"
	StatusRes[END] = "END"
	StatusRes[DELETED] = "DELETED"
	StatusRes[NOT_FOUND] = "NOT_FOUND"
	StatusRes[UNKNOWN_COMMAND] = "UNKNOWN_COMMAND"
}
