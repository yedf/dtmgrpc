/*
 * Copyright (c) 2021 yedf. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package dtmgrpc

import (
	"database/sql"

	"github.com/dtm-labs/dtmcli"
	"github.com/dtm-labs/dtmcli/dtmimp"
	"github.com/dtm-labs/dtmgrpc/dtmgimp"
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

// PrepareAndSubmit one method for the entire busi->prepare->submit
func (s *MsgGrpc) PrepareAndSubmit(queryPrepared string, db *sql.DB, busiCall dtmcli.BarrierBusiFunc) error {
	bb, err := dtmcli.BarrierFrom(s.TransType, s.Gid, "00", "msg") // a special barrier for msg QueryPrepared
	if err == nil {
		err = bb.CallWithDB(db, func(tx *sql.Tx) error {
			err := busiCall(tx)
			if err == nil {
				err = s.Prepare(queryPrepared)
			}
			return err
		})
	}
	if err == nil {
		err = s.Submit()
	}
	return err
}
