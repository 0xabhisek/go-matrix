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

package filters

import (
	"bytes"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/matrix/go-matrix/common"
	"github.com/matrix/go-matrix/common/bitutil"
	"github.com/matrix/go-matrix/core/bloombits"
	"github.com/matrix/go-matrix/core/rawdb"
	"github.com/matrix/go-matrix/core/types"
	"github.com/matrix/go-matrix/mandb"
	"github.com/matrix/go-matrix/event"
	"github.com/matrix/go-matrix/node"
)

func BenchmarkBloomBits512(b *testing.B) {
	benchmarkBloomBits(b, 512)
}

func BenchmarkBloomBits1k(b *testing.B) {
	benchmarkBloomBits(b, 1024)
}

func BenchmarkBloomBits2k(b *testing.B) {
	benchmarkBloomBits(b, 2048)
}

func BenchmarkBloomBits4k(b *testing.B) {
	benchmarkBloomBits(b, 4096)
}

func BenchmarkBloomBits8k(b *testing.B) {
	benchmarkBloomBits(b, 8192)
}

func BenchmarkBloomBits16k(b *testing.B) {
	benchmarkBloomBits(b, 16384)
}

func BenchmarkBloomBits32k(b *testing.B) {
	benchmarkBloomBits(b, 32768)
}

const benchFilterCnt = 2000

func benchmarkBloomBits(b *testing.B, sectionSize uint64) {
	benchDataDir := node.DefaultDataDir() + "/gman/chaindata"
	fmt.Println("Running bloombits benchmark   section size:", sectionSize)

	db, err := mandb.NewLDBDatabase(benchDataDir, 128, 1024)
	if err != nil {
		b.Fatalf("error opening database at %v: %v", benchDataDir, err)
	}
	head := rawdb.ReadHeadBlockHash(db)
	if head == (common.Hash{}) {
		b.Fatalf("chain data not found at %v", benchDataDir)
	}

	clearBloomBits(db)
	fmt.Println("Generating bloombits data...")
	headNum := rawdb.ReadHeaderNumber(db, head)
	if headNum == nil || *headNum < sectionSize+512 {
		b.Fatalf("not enough blocks for running a benchmark")
	}

	start := time.Now()
	cnt := (*headNum - 512) / sectionSize
	var dataSize, compSize uint64
	for sectionIdx := uint64(0); sectionIdx < cnt; sectionIdx++ {
		bc, err := bloombits.NewGenerator(uint(sectionSize))
		if err != nil {
			b.Fatalf("failed to create generator: %v", err)
		}
		var header *types.Header
		for i := sectionIdx * sectionSize; i < (sectionIdx+1)*sectionSize; i++ {
			hash := rawdb.ReadCanonicalHash(db, i)
			header = rawdb.ReadHeader(db, hash, i)
			if header == nil {
				b.Fatalf("Error creating bloomBits data")
			}
			bc.AddBloom(uint(i-sectionIdx*sectionSize), header.Bloom)
		}
		sectionHead := rawdb.ReadCanonicalHash(db, (sectionIdx+1)*sectionSize-1)
		for i := 0; i < types.BloomBitLength; i++ {
			data, err := bc.Bitset(uint(i))
			if err != nil {
				b.Fatalf("failed to retrieve bitset: %v", err)
			}
			comp := bitutil.CompressBytes(data)
			dataSize += uint64(len(data))
			compSize += uint64(len(comp))
			rawdb.WriteBloomBits(db, uint(i), sectionIdx, sectionHead, comp)
		}
		//if sectionIdx%50 == 0 {
		//	fmt.Println(" section", sectionIdx, "/", cnt)
		//}
	}

	d := time.Since(start)
	fmt.Println("Finished generating bloombits data")
	fmt.Println(" ", d, "total  ", d/time.Duration(cnt*sectionSize), "per block")
	fmt.Println(" data size:", dataSize, "  compressed size:", compSize, "  compression ratio:", float64(compSize)/float64(dataSize))

	fmt.Println("Running filter benchmarks...")
	start = time.Now()
	mux := new(event.TypeMux)
	var backend *testBackend

	for i := 0; i < benchFilterCnt; i++ {
		if i%20 == 0 {
			db.Close()
			db, _ = mandb.NewLDBDatabase(benchDataDir, 128, 1024)
			backend = &testBackend{mux, db, cnt, new(event.Feed), new(event.Feed), new(event.Feed), new(event.Feed)}
		}
		var addr common.Address
		addr[0] = byte(i)
		addr[1] = byte(i / 256)
		filter := New(backend, 0, int64(cnt*sectionSize-1), []common.Address{addr}, nil)
		if _, err := filter.Logs(context.Background()); err != nil {
			b.Error("filter.Find error:", err)
		}
	}
	d = time.Since(start)
	fmt.Println("Finished running filter benchmarks")
	fmt.Println(" ", d, "total  ", d/time.Duration(benchFilterCnt), "per address", d*time.Duration(1000000)/time.Duration(benchFilterCnt*cnt*sectionSize), "per million blocks")
	db.Close()
}

func forEachKey(db mandb.Database, startPrefix, endPrefix []byte, fn func(key []byte)) {
	it := db.(*mandb.LDBDatabase).NewIterator()
	it.Seek(startPrefix)
	for it.Valid() {
		key := it.Key()
		cmpLen := len(key)
		if len(endPrefix) < cmpLen {
			cmpLen = len(endPrefix)
		}
		if bytes.Compare(key[:cmpLen], endPrefix) == 1 {
			break
		}
		fn(common.CopyBytes(key))
		it.Next()
	}
	it.Release()
}

var bloomBitsPrefix = []byte("bloomBits-")

func clearBloomBits(db mandb.Database) {
	fmt.Println("Clearing bloombits data...")
	forEachKey(db, bloomBitsPrefix, bloomBitsPrefix, func(key []byte) {
		db.Delete(key)
	})
}

func BenchmarkNoBloomBits(b *testing.B) {
	benchDataDir := node.DefaultDataDir() + "/gman/chaindata"
	fmt.Println("Running benchmark without bloombits")
	db, err := mandb.NewLDBDatabase(benchDataDir, 128, 1024)
	if err != nil {
		b.Fatalf("error opening database at %v: %v", benchDataDir, err)
	}
	head := rawdb.ReadHeadBlockHash(db)
	if head == (common.Hash{}) {
		b.Fatalf("chain data not found at %v", benchDataDir)
	}
	headNum := rawdb.ReadHeaderNumber(db, head)

	clearBloomBits(db)

	fmt.Println("Running filter benchmarks...")
	start := time.Now()
	mux := new(event.TypeMux)
	backend := &testBackend{mux, db, 0, new(event.Feed), new(event.Feed), new(event.Feed), new(event.Feed)}
	filter := New(backend, 0, int64(*headNum), []common.Address{{}}, nil)
	filter.Logs(context.Background())
	d := time.Since(start)
	fmt.Println("Finished running filter benchmarks")
	fmt.Println(" ", d, "total  ", d*time.Duration(1000000)/time.Duration(*headNum+1), "per million blocks")
	db.Close()
}
