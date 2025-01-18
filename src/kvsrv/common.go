package kvsrv

// Put or Append
type PutAppendArgs struct {
	Key     string
	Value   string
	ClkID   int64
	Version bool
}

type PutAppendReply struct {
	Value   string
	ClkID   int64
	Version bool
}

type GetArgs struct {
	Key     string
	ClkID   int64
	Version bool
}

type GetReply struct {
	Value   string
	ClkID   int64
	Version bool
}
