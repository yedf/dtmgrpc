package dtmgrpc

import (
	"github.com/yedf/dtmcli"
	"github.com/yedf/dtmcli/dtmimp"
	"github.com/yedf/dtmgrpc/dtmgimp"
	"google.golang.org/protobuf/proto"
)

// MsgGrpc reliable msg type
type MsgGrpc struct {
	dtmcli.Msg
}

// NewMsgGrpc create new msg
func NewMsgGrpc(server string, gid string) *MsgGrpc {
	return &MsgGrpc{Msg: *dtmcli.NewMsg(server, gid)}
}

// Add add a new step
func (s *MsgGrpc) Add(action string, msg proto.Message) *MsgGrpc {
	s.Steps = append(s.Steps, map[string]string{"action": action})
	s.BinPayloads = append(s.BinPayloads, dtmgimp.MustProtoMarshal(msg))
	return s
}

// Prepare prepare the msg, msg will later be submitted
func (s *MsgGrpc) Prepare(queryPrepared string) error {
	s.QueryPrepared = dtmimp.OrString(queryPrepared, s.QueryPrepared)
	return dtmgimp.DtmGrpcCall(&s.TransBase, "Prepare")
}

// Submit submit the msg
func (s *MsgGrpc) Submit() error {
	return dtmgimp.DtmGrpcCall(&s.TransBase, "Submit")
}
