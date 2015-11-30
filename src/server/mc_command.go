// mc_command
package main

import (
	"cache"
	"proto"
	"time"
)

var (
	CmdFuncs = map[proto.CommandCode]func(*proto.MCRequest, *cache.Cache) *proto.MCResponse{}
)

func gets(req *proto.MCRequest, cache *cache.Cache) *proto.MCResponse {
	value, found := cache.Get(req.Key)
	res := proto.NewResFull(req.Opcode, proto.SUCCESS, req.Key, req.Flags, found)
	if found == false {
		res.Status = proto.END
	} else {
		valueb := value.([]byte)
		res.Value = make([]byte, len(valueb))
		copy(res.Value, valueb)
	}
	return res
}

func set(req *proto.MCRequest, cache *cache.Cache) *proto.MCResponse {
	cache.Set(req.Key, req.Value, time.Duration(req.Expires)*time.Second)
	res := proto.NewResStatus(req.Opcode, proto.SUCCESS)
	return res
}

func add(req *proto.MCRequest, cache *cache.Cache) *proto.MCResponse {
	res := proto.NewResStatus(req.Opcode, proto.STORED)

	if err := cache.Add(req.Key, req.Value,
		time.Duration(req.Expires)*time.Second); err != nil {
		res = proto.NewResStatus(req.Opcode, proto.NOT_STORED)
	}
	return res
}
func replace(req *proto.MCRequest, cache *cache.Cache) *proto.MCResponse {
	res := proto.NewResStatus(req.Opcode, proto.STORED)
	if err := cache.Replace(req.Key, req.Value,
		time.Duration(req.Expires)*time.Second); err != nil {
		res = proto.NewResStatus(req.Opcode, proto.NOT_STORED)
	}
	return res
}
func delete(req *proto.MCRequest, cache *cache.Cache) *proto.MCResponse {
	cache.Delete(req.Key)
	res := proto.NewResStatus(req.Opcode, proto.DELETED)
	return res
}

func init() {
	CmdFuncs[proto.GET] = gets
	CmdFuncs[proto.SET] = set
	CmdFuncs[proto.ADD] = add
	CmdFuncs[proto.REPLACE] = replace
	CmdFuncs[proto.DELETE] = delete
}
