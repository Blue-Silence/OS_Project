package File

import (
	"LSF/DiskLayer"
	"LSF/Setting"
	"bytes"
	"encoding/gob"
	"log"
)

func (s FolderBlock) ToBlock() DiskLayer.RealBlock {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(s)
	if err != nil {
		log.Fatal("encode error:", err)
	}
	if buf.Len() > Setting.BlockSize {
		log.Fatal("FolderBlock is too big to be a block.Need filesystem adjustment.")
	}
	return DiskLayer.BytesToBlock(buf.Bytes())
}

func (s FolderBlock) FromBlock(d DiskLayer.RealBlock) DiskLayer.Block {
	bufP := bytes.NewBuffer(DiskLayer.BlockToBytes(d))
	dec := gob.NewDecoder(bufP)
	err := dec.Decode(&s)
	if err != nil {
		log.Fatal("decode error 1:", err)
	}
	return s
}
