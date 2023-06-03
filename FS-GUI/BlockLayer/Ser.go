package BlockLayer

import (
	"LSF/DiskLayer"
	"LSF/Setting"
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"runtime/debug"
)

func (s SuperBlock) ToBlocks() []DiskLayer.RealBlock {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(s)
	if err != nil {
		log.Fatal("encode error:", err)
	}
	return DiskLayer.BytesToBlocks(buf.Bytes())
}

func (s SuperBlock) FromBlocks(ds []DiskLayer.RealBlock) SuperBlock {
	d := []byte{}
	for _, v := range ds {
		d = append(d, DiskLayer.BlockToBytes(v)...)
	}

	bufP := bytes.NewBuffer(d)
	dec := gob.NewDecoder(bufP)
	err := dec.Decode(&s)
	if err != nil {
		debug.PrintStack()
		log.Fatal("decode error 1:", err)
	}
	return s
}

func (s SegHead) ToBlock() DiskLayer.RealBlock {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(s)
	if err != nil {
		log.Fatal("encode error:", err)
	}
	if buf.Len() > Setting.BlockSize {
		log.Fatal("SegHead is too big to be a block.Need filesystem adjustment.Len:", buf.Len())
	}
	return DiskLayer.BytesToBlock(buf.Bytes())
}

func (s SegHead) FromBlock(d DiskLayer.RealBlock) DiskLayer.Block {
	bufP := bytes.NewBuffer(DiskLayer.BlockToBytes(d))
	dec := gob.NewDecoder(bufP)
	err := dec.Decode(&s)
	if err != nil {
		fmt.Println(string(debug.Stack()))
		log.Fatal("decode error 1:", err)
	}
	return s
}

func (s INodeMap) ToBlock() DiskLayer.RealBlock {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(s)
	if err != nil {
		log.Fatal("encode error:", err)
	}
	if buf.Len() > Setting.BlockSize {
		log.Fatal("INodeMap is too big to be a block.Need filesystem adjustment.")
	}
	return DiskLayer.BytesToBlock(buf.Bytes())
}

func (s INodeMap) FromBlock(d DiskLayer.RealBlock) DiskLayer.Block {
	bufP := bytes.NewBuffer(DiskLayer.BlockToBytes(d))
	dec := gob.NewDecoder(bufP)
	err := dec.Decode(&s)
	if err != nil {
		debug.PrintStack()
		log.Fatal("decode error 1:", err)
	}
	return s
}

func (s INodeBlock) ToBlock() DiskLayer.RealBlock {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(s)
	if err != nil {
		log.Fatal("encode error:", err)
	}
	if buf.Len() > Setting.BlockSize {
		log.Fatal("INodeBlock is too big to be a block.Need filesystem adjustment.")
	}
	return DiskLayer.BytesToBlock(buf.Bytes())
}

func (s INodeBlock) FromBlock(d DiskLayer.RealBlock) DiskLayer.Block {
	bufP := bytes.NewBuffer(DiskLayer.BlockToBytes(d))
	dec := gob.NewDecoder(bufP)
	err := dec.Decode(&s)
	if err != nil {
		debug.PrintStack()
		log.Fatal("decode error 1:", err)
	}
	return s
}

func (s DataBlock) ToBlock() DiskLayer.RealBlock {
	return DiskLayer.BytesToBlock(s.Data[:])
}

func (s DataBlock) FromBlock(d DiskLayer.RealBlock) DiskLayer.Block {
	copy(s.Data[:], DiskLayer.BlockToBytes(d))
	return s
}
