// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zip

import (
	"encoding/binary"
	"errors"
	"io"
)

var (
	ErrFormat       = errors.New("zip: not a valid zip file")
	ErrAlgorithm    = errors.New("zip: unsupported compression algorithm")
	ErrChecksum     = errors.New("zip: checksum error")
	ErrInsecurePath = errors.New("zip: insecure file path")
)

type FileContent struct {
	Filename []byte
	Content  []byte
}

func ReadFile(r io.ReadSeeker) (*FileContent, error) {
	var buf [fileHeaderLen]byte
	if _, err := io.ReadFull(r, buf[:]); err != nil {
		return nil, err
	}
	b := readBuf(buf[:])
	if sig := b.uint32(); sig != fileHeaderSignature {
		return nil, ErrFormat
	}
	b = b[18:] // skip over most of the header
	fileContentLen := int(b.uint32())
	filenameLen := int(b.uint16())
	extraLen := int64(b.uint16())
	filename := make([]byte, filenameLen)
	if _, err := io.ReadFull(r, filename[:]); err != nil {
		return nil, err
	}
	if _, err := r.Seek(extraLen, io.SeekCurrent); err != nil {
		return nil, err
	}
	// pos, err := r.Seek(0, io.SeekCurrent)
	// if err != nil {
	// 	return nil, err
	// }
	// fmt.Printf("pso: %d", pos)
	content := make([]byte, fileContentLen)
	if _, err := io.ReadFull(r, content[:]); err != nil {
		return nil, err
	}
	// pos, err = r.Seek(0, io.SeekCurrent)
	// if err != nil {
	// 	return nil, err
	// }
	// fmt.Printf("pso: %d", pos)
	return &FileContent{
		Filename: filename,
		Content:  content,
	}, nil
}


type readBuf []byte

func (b *readBuf) uint8() uint8 {
	v := (*b)[0]
	*b = (*b)[1:]
	return v
}

func (b *readBuf) uint16() uint16 {
	v := binary.LittleEndian.Uint16(*b)
	*b = (*b)[2:]
	return v
}

func (b *readBuf) uint32() uint32 {
	v := binary.LittleEndian.Uint32(*b)
	*b = (*b)[4:]
	return v
}

func (b *readBuf) uint64() uint64 {
	v := binary.LittleEndian.Uint64(*b)
	*b = (*b)[8:]
	return v
}

func (b *readBuf) sub(n int) readBuf {
	b2 := (*b)[:n]
	*b = (*b)[n:]
	return b2
}
