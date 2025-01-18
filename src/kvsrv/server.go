package kvsrv

import (
	"log"
	"sync"
)

const Debug = false

func DPrintf(format string, a ...interface{}) (n int, err error) {
	if Debug {
		log.Printf(format, a...)
	}
	return
}

type KVServer struct {
	mu sync.RWMutex

	data map[string]string
	// FIXME: This implement will not release `versions`, which will cause memory leak when the client will be closed.
	cache    map[int64]string
	versions map[int64]bool
}

func (kv *KVServer) Get(args *GetArgs, reply *GetReply) {
	kv.mu.RLock()
	defer kv.mu.RUnlock()
	value, exists := kv.data[args.Key]
	if exists {
		reply.Value = value
	} else {
		reply.Value = ""
	}
}

func (kv *KVServer) Put(args *PutAppendArgs, reply *PutAppendReply) {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	if version, exists := kv.versions[args.ClkID]; !(exists && (version == args.Version)) { // new operation
		kv.data[args.Key] = args.Value
		kv.cache[args.ClkID] = "" // FIXME: delete instead of setting to empty string
		kv.versions[args.ClkID] = args.Version
	}
}

func (kv *KVServer) Append(args *PutAppendArgs, reply *PutAppendReply) {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	if version, exists := kv.versions[args.ClkID]; !(exists && version == args.Version) { // new operation
		old, exists := kv.data[args.Key]
		if exists {
			kv.data[args.Key] = old + args.Value
		} else {
			kv.data[args.Key] = args.Value
		}
		kv.cache[args.ClkID] = old
		kv.versions[args.ClkID] = args.Version
	}
	reply.Value = kv.cache[args.ClkID]
}

func StartKVServer() *KVServer {
	kv := new(KVServer)

	kv.data = make(map[string]string)
	kv.cache = make(map[int64]string)
	kv.versions = make(map[int64]bool)
	return kv
}
