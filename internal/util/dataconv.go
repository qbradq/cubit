package util

import (
	"encoding/binary"
	"io"
	"math"
	"time"
)

// PutBool writes a boolean value
func PutBool(w io.Writer, v bool) {
	var b [1]byte
	if v {
		b[0] = 1
	} else {
		b[0] = 0
	}
	w.Write(b[:])
}

// PutString writes a null-terminated string
func PutString(w io.Writer, s string) {
	var b [1]byte
	w.Write([]byte(s))
	w.Write(b[:])
}

// Pad writes zero-padding
func Pad(w io.Writer, l int) {
	var buf [1024]byte
	w.Write(buf[:l])
}

// PutByte writes a single byte
func PutByte(w io.Writer, v byte) {
	var b [1]byte
	b[0] = v
	w.Write(b[:])
}

// PutUint16 writes a 16-bit numeric value
func PutUint16(w io.Writer, v uint16) {
	var b [2]byte
	binary.LittleEndian.PutUint16(b[:], v)
	w.Write(b[:])
}

// PutUint32 writes a 32-bit numeric value
func PutUint32(w io.Writer, v uint32) {
	var b [4]byte
	binary.LittleEndian.PutUint32(b[:], v)
	w.Write(b[:])
}

// PutUint64 writes a 64-bit numeric value
func PutUint64(w io.Writer, v uint64) {
	var b [8]byte
	binary.LittleEndian.PutUint64(b[:], v)
	w.Write(b[:])
}

// PutFloat writes a 64-bit floating-point value.
func PutFloat(w io.Writer, v float64) {
	PutUint64(w, math.Float64bits(v))
}

// PutTime writes a time.Time value to the writer.
func PutTime(w io.Writer, t time.Time) {
	PutUint64(w, uint64(t.UnixMilli()))
}

// PutBytes puts a slice of bytes to the writer.
func PutBytes(w io.Writer, d []byte) {
	PutUint32(w, uint32(len(d)))
	w.Write(d)
}

// GetBool returns the next boolean value in the data buffer.
func GetBool(r io.Reader) bool {
	return GetByte(r) != 0
}

// GetString returns the next null-terminated string in the data buffer.
func GetString(r io.Reader) string {
	var buf = []byte{0}
	var ret []byte
	for {
		r.Read(buf)
		if buf[0] == 0 {
			return string(ret)
		}
		ret = append(ret, buf[0])
	}
}

// GetByte is a convenience function that returns the next byte in the buffer.
func GetByte(r io.Reader) byte {
	var buf = []byte{0}
	r.Read(buf)
	return buf[0]
}

// GetUint16 returns the next unsigned 16-bit integer in the data buffer.
func GetUint16(r io.Reader) uint16 {
	var buf = []byte{0, 0}
	r.Read(buf)
	return binary.LittleEndian.Uint16(buf)
}

// GetUint32 returns the next unsigned 32-bit integer in the data buffer.
func GetUint32(r io.Reader) uint32 {
	var buf = []byte{0, 0, 0, 0}
	r.Read(buf)
	return binary.LittleEndian.Uint32(buf)
}

// GetUint64 returns the next unsigned 64-bit integer in the data buffer.
func GetUint64(r io.Reader) uint64 {
	var buf = []byte{0, 0, 0, 0, 0, 0, 0, 0}
	r.Read(buf)
	return binary.LittleEndian.Uint64(buf)
}

// GetFloat returns the next 64-bit floating-point value.
func GetFloat(r io.Reader) float64 {
	return math.Float64frombits(GetUint64(r))
}

// GetTime returns the next time.Time value in the data buffer.
func GetTime(r io.Reader) time.Time {
	return time.UnixMilli(int64(GetUint64(r)))
}

// GetBytes returns the next byte slice in the data buffer.
func GetBytes(r io.Reader) []byte {
	n := GetUint32(r)
	ret := make([]byte, n)
	r.Read(ret)
	return ret
}
