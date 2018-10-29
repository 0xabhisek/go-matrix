// Copyright 2018 The MATRIX Authors as well as Copyright 2014-2017 The go-ethereum Authors
// This file is consisted of the MATRIX library and part of the go-ethereum library.
//
// The MATRIX-ethereum library is free software: you can redistribute it and/or modify it under the terms of the MIT License.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"),
// to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, 
//and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject tothe following conditions:
//
//The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
//
//THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
//FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, 
//WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISINGFROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
//OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
// Code generated by github.com/fjl/gencodec. DO NOT EDIT.

package whisperv5

import (
	"encoding/json"

	"github.com/matrix/go-matrix/common/hexutil"
)

var _ = (*messageOverride)(nil)

func (m Message) MarshalJSON() ([]byte, error) {
	type Message struct {
		Sig       hexutil.Bytes `json:"sig,omitempty"`
		TTL       uint32        `json:"ttl"`
		Timestamp uint32        `json:"timestamp"`
		Topic     TopicType     `json:"topic"`
		Payload   hexutil.Bytes `json:"payload"`
		Padding   hexutil.Bytes `json:"padding"`
		PoW       float64       `json:"pow"`
		Hash      hexutil.Bytes `json:"hash"`
		Dst       hexutil.Bytes `json:"recipientPublicKey,omitempty"`
	}
	var enc Message
	enc.Sig = m.Sig
	enc.TTL = m.TTL
	enc.Timestamp = m.Timestamp
	enc.Topic = m.Topic
	enc.Payload = m.Payload
	enc.Padding = m.Padding
	enc.PoW = m.PoW
	enc.Hash = m.Hash
	enc.Dst = m.Dst
	return json.Marshal(&enc)
}

func (m *Message) UnmarshalJSON(input []byte) error {
	type Message struct {
		Sig       *hexutil.Bytes `json:"sig,omitempty"`
		TTL       *uint32        `json:"ttl"`
		Timestamp *uint32        `json:"timestamp"`
		Topic     *TopicType     `json:"topic"`
		Payload   *hexutil.Bytes `json:"payload"`
		Padding   *hexutil.Bytes `json:"padding"`
		PoW       *float64       `json:"pow"`
		Hash      *hexutil.Bytes `json:"hash"`
		Dst       *hexutil.Bytes `json:"recipientPublicKey,omitempty"`
	}
	var dec Message
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}
	if dec.Sig != nil {
		m.Sig = *dec.Sig
	}
	if dec.TTL != nil {
		m.TTL = *dec.TTL
	}
	if dec.Timestamp != nil {
		m.Timestamp = *dec.Timestamp
	}
	if dec.Topic != nil {
		m.Topic = *dec.Topic
	}
	if dec.Payload != nil {
		m.Payload = *dec.Payload
	}
	if dec.Padding != nil {
		m.Padding = *dec.Padding
	}
	if dec.PoW != nil {
		m.PoW = *dec.PoW
	}
	if dec.Hash != nil {
		m.Hash = *dec.Hash
	}
	if dec.Dst != nil {
		m.Dst = *dec.Dst
	}
	return nil
}
