package kvsrv

import (
	"crypto/rand"
	"math/big"
	"time"

	"6.5840/labrpc"
)

type Clerk struct {
	server  *labrpc.ClientEnd
	clkID   int64
	version bool
}

func nrand() int64 {
	max := big.NewInt(int64(1) << 62)
	bigx, _ := rand.Int(rand.Reader, max)
	x := bigx.Int64()
	return x
}

func MakeClerk(server *labrpc.ClientEnd) *Clerk {
	ck := new(Clerk)
	ck.server = server
	ck.clkID = nrand()
	return ck
}

func (ck *Clerk) Get(key string) string {
	// ck.version = !ck.version
	args := &GetArgs{Key: key, ClkID: ck.clkID, Version: ck.version}
	reply := &GetReply{}

	for {
		// Try the current server
		if ok := ck.server.Call("KVServer.Get", args, reply); ok {
			return reply.Value
		}

		// If RPC failed, sleep briefly to avoid overwhelming the network
		time.Sleep(100 * time.Millisecond)
	}
}

func (ck *Clerk) PutAppend(key string, value string, op string) string {
	ck.version = !ck.version
	args := &PutAppendArgs{Key: key, Value: value, ClkID: ck.clkID, Version: ck.version}
	reply := &PutAppendReply{}
	for {
		if ok := ck.server.Call("KVServer."+op, args, reply); ok {
			return reply.Value
		}

		time.Sleep(100 * time.Millisecond)
	}
}

func (ck *Clerk) Put(key string, value string) {
	ck.PutAppend(key, value, "Put")
}

// Append value to key's value and return that value
func (ck *Clerk) Append(key string, value string) string {
	return ck.PutAppend(key, value, "Append")
}
