package bitstructs

import (
	"fmt"
)

var BYTE_LENGTH int = 8

type BitSequence struct {
	data           *[]uint8
	numBits        int
	bytesAllocated int
	nextBitIdx     *int
	nextByteIdx    *int
}

func NewBitSequence(numBits int) BitSequence {

	bytesAllocated := numBits / BYTE_LENGTH
	if numBits%BYTE_LENGTH != 0 { // if there is a remainder we need to allocate an extra byte to accommodate the extra number of bit-structs
		bytesAllocated += 1
	}

	data := make([]uint8, bytesAllocated)
	nextBitIdx := 0
	nextByteIdx := 0
	return BitSequence{data: &data, numBits: numBits, bytesAllocated: bytesAllocated, nextBitIdx: &nextBitIdx, nextByteIdx: &nextByteIdx}
}

func NewBitSequenceFromByteArray(data *[]uint8, bitLen int) BitSequence {

	bitSeq := NewBitSequence(bitLen)

	bitSeq.SetNextByteStart(0)
	for i := 0; i < len(*data); i++ {
		bitSeq.SetNextByte((*data)[i])
	}
	return bitSeq
}

func (bseq BitSequence) SetBit(bitIdx int, set bool) {
	if bitIdx < 0 || bitIdx >= bseq.numBits {
		panic("set_bit idx out of bounds")
	}
	byteIdx := bitIdx / BYTE_LENGTH
	bitIdxInByte := bitIdx % BYTE_LENGTH

	if set {
		(*bseq.data)[byteIdx] |= 1 << bitIdxInByte
	} else {
		(*bseq.data)[byteIdx] &= ^(1 << bitIdxInByte)
	}
}

func (bseq BitSequence) SetBitsFromNum(startBitIdx int, number uint64) {

	if startBitIdx < 0 || startBitIdx >= bseq.numBits {
		panic("BitSequence.setBitsFromNum(..): idx out of bounds")
	}
	for number != 0 && startBitIdx < bseq.numBits {
		bseq.SetBit(startBitIdx, numToBool(uint8(number&0x1)))
		number >>= 1
		startBitIdx++
	}
}

func (bseq BitSequence) SetNextBitStart(bitIdx int) {
	if bitIdx < 0 || bitIdx >= bseq.numBits {
		panic("BitSequence.setNextBitStart(..): idx out of bounds")
	}
	*bseq.nextBitIdx = bitIdx

}

func (bseq BitSequence) SetNextBit(set bool) {
	if *bseq.nextBitIdx == -1 {
		panic("BitSequence.setNextBit(..) never called")
	}
	bseq.SetBit(*bseq.nextBitIdx, set)
	*bseq.nextBitIdx++
}

func (bseq BitSequence) SetNextByteStart(byteIdx int) {
	if byteIdx < 0 || byteIdx >= bseq.bytesAllocated {
		panic("BitSequence.setNextByteStart(): idx out of bounds")
	}
	*bseq.nextByteIdx = byteIdx
}

func (bseq BitSequence) SetNextByte(byte uint8) {
	if *bseq.nextByteIdx == -1 {
		panic("BitSequence.setNextByte() never called\n")

	}
	(*bseq.data)[*bseq.nextByteIdx] = byte
	*bseq.nextByteIdx++
}

func (bseq BitSequence) GetBit(bitIdx int) bool {
	if bitIdx < 0 || bitIdx >= bseq.numBits {
		panic("BitSequence.getBit() idx out of bounds")
		//return false
	}
	byteIdx := bitIdx / BYTE_LENGTH
	bitIdxInByte := bitIdx % BYTE_LENGTH

	zeroOrOne := ((*bseq.data)[byteIdx] >> bitIdxInByte) & 0x1

	return numToBool(zeroOrOne)

}

func (bseq BitSequence) GetByte(byteIdx int) uint8 {
	if byteIdx < 0 || byteIdx >= bseq.bytesAllocated {
		panic("BitSequence.getByte() idx out of bounds")
	}
	return (*bseq.data)[byteIdx]

}

func (bseq BitSequence) GetXBytes(numBytes int) uint64 {

	if numBytes > bseq.bytesAllocated {
		numBytes = bseq.bytesAllocated
	}
	res := uint64(0)

	for numBytes != 0 {
		res <<= BYTE_LENGTH
		res |= uint64(bseq.GetByte(numBytes - 1))
		numBytes--
	}
	return res
}

func (bseq BitSequence) GetNextBitStart(bitIdx int) {
	if bitIdx < 0 || bitIdx >= bseq.numBits {
		panic("BitSequence.getNextBitStart: idx out of bounds\n")
	}
	*bseq.nextBitIdx = bitIdx
}

func (bseq BitSequence) GetNextBit() bool {
	if *bseq.nextBitIdx == -1 {
		panic("BitSequence.getNextBitStart() never called")
	}
	if *bseq.nextBitIdx >= bseq.numBits {
		panic("BitSequence.GetNextBit() called one too many times. Index Out of bounds")
	}
	res := bseq.GetBit(*bseq.nextBitIdx)
	*bseq.nextBitIdx++
	return res
}

func (bseq BitSequence) GetNextByte() uint8 {
	if *bseq.nextByteIdx == -1 {
		panic("BitSequence.getNextByteStart() never called")
	}
	if *bseq.nextByteIdx >= bseq.bytesAllocated {
		panic("BitSequence.getNextByte() called one too many times. Index Out of bounds")
	}
	res := bseq.GetByte(*bseq.nextByteIdx)
	*bseq.nextByteIdx++
	return res
}

func (bseq BitSequence) GetBitSeq() []uint8 {
	return *bseq.data
}

func (bseq BitSequence) GetNumBits() int {
	return bseq.numBits
}
func (bseq BitSequence) GetBytesAllocated() int {
	return len(*bseq.data)
}
func (bseq BitSequence) GetNextBitIdx() int {
	return *bseq.nextBitIdx
}
func (bseq BitSequence) GetNextByteIdx() int {
	return *bseq.nextByteIdx
}

func (bseq BitSequence) Clear() {
	for i, _ := range *bseq.data {
		(*bseq.data)[i] = 0
	}
}

func (bseq BitSequence) String() string {
	result := ""

	//for i :=  0; i < bseq.bytesAllocated; i++ {
	//	result += fmt.Sprintf("%02X", (*bseq.data)[i])
	//}
	for i := bseq.numBits - 1; i >= 0; i-- {
		result += fmt.Sprintf("%v", boolToInt(bseq.GetBit(i)))
	}
	return result
}

func numToBool(n uint8) bool {
	if n == 0 {
		return false
	}
	return true
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
