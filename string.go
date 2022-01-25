package tool

import (
	"bytes"
	"errors"
	"github.com/iEvan-lhr/string/evan"

	"unicode/utf8"
	"unsafe"
)

type String struct {
	addr  *String
	runes []rune
	buf   []byte
}

// EString 根据字符串来构建一个String
func EString(str string) *String {
	s := String{}
	_, err := s.writeString(str)
	if err != nil {
		evan.ErrorLog(err)
		return nil
	}
	return &s
}

// ToString 字符串转型输出
func (s *String) ToString() string {
	return s.string()

}

// JoinStrString 拼接字符串
func (s *String) JoinStrString(str *String) {
	_, err := s.Write(str.buf)
	evan.ErrorLog(err)
}

// JoinString 拼接字符串
func (s *String) JoinString(str string) {
	_, err := s.writeString(str)
	evan.ErrorLog(err)
}

// IndexString 返回数据中含有字串的下标 没有返回-1
func (s *String) IndexString(str *String) int {
	return bytes.Index(s.buf, str.buf)
}

// Index 返回数据中含有字串的下标 没有返回-1
func (s *String) Index(str string) int {
	return bytes.Index(s.buf, []byte(str))

}

// Split 按照string片段来分割字符串 返回[]string
func (s *String) Split(str string) []string {
	var order []string
	for _, v := range bytes.Split(s.buf, []byte(str)) {
		order = append(order, string(v))
	}
	return order
}

// SplitString 按照*String来分割字符串 返回[]*String
func (s *String) SplitString(str String) []*String {
	byt := bytes.Split(s.buf, str.buf)
	var order []*String
	for i := range byt {
		order = append(order, &String{buf: byt[i]})
	}
	return order
}

func noescape(p unsafe.Pointer) unsafe.Pointer {
	x := uintptr(p)
	return unsafe.Pointer(x ^ 0)
}

func (s *String) copyCheck() {
	if s.addr == nil {
		s.addr = (*String)(noescape(unsafe.Pointer(s)))
	} else if s.addr != s {
		panic("strings: illegal use of non-zero String copied by value")
	}
}

func (s *String) string() string {
	return *(*string)(unsafe.Pointer(&s.buf))
}

// Len 返回字符串长度
func (s *String) Len() int { return len(s.buf) }

// LenByRune 返回含有中文的字符串长度
func (s *String) LenByRune() int { return len(bytes.Runes(s.buf)) }

func (s *String) cap() int { return cap(s.buf) }

func (s *String) reset() {
	s.addr = nil
	s.buf = nil
}

func (s *String) grow(n int) {
	buf := make([]byte, len(s.buf), 2*cap(s.buf)+n)
	copy(buf, s.buf)
	s.buf = buf
}

// Grow  扩充大小
func (s *String) Grow(n int) {
	s.copyCheck()
	if n < 0 {
		panic("strings.String.Grow: negative count")
	}
	if cap(s.buf)-len(s.buf) < n {
		s.grow(n)
	}
}

// WriteByte 写入[]Byte的数据
func (s *String) Write(p []byte) (int, error) {
	s.copyCheck()
	s.buf = append(s.buf, p...)
	return len(p), nil
}

// WriteByte 写入Byte字符格式的数据
func (s *String) WriteByte(c byte) error {
	s.copyCheck()
	s.buf = append(s.buf, c)
	return nil
}

// WriteRune 写入Rune字符格式的数据
func (s *String) WriteRune(r rune) (int, error) {
	s.copyCheck()
	if r < utf8.RuneSelf {
		s.buf = append(s.buf, byte(r))
		return 1, nil
	}
	l := len(s.buf)
	if cap(s.buf)-l < utf8.UTFMax {
		s.grow(utf8.UTFMax)
	}
	n := utf8.EncodeRune(s.buf[l:l+utf8.UTFMax], r)
	s.buf = s.buf[:l+n]
	return n, nil
}

func (s *String) writeString(str string) (int, error) {
	s.copyCheck()
	s.buf = append(s.buf, str...)
	return len(str), nil
}

// RemoveLastStr 从尾部移除固定长度的字符
func (s *String) RemoveLastStr(lens int) {
	if lens > s.Len() {
		evan.ErrorLog(errors.New("RemoveLens>stringLens Please Check"))
		return
	}
	s.buf = s.buf[:s.Len()-lens]
}

// RemoveLastStrByRune 从尾部移除固定长度的字符 并且支持中文字符的移除
func (s *String) RemoveLastStrByRune(lens int) {
	runes := bytes.Runes(s.buf)
	if lens > len(runes) {
		evan.ErrorLog(errors.New("RemoveLens>stringLens Please Check"))
		return
	}
	s.buf = RunesToBytes(runes[:len(runes)-lens])
}

// GetByte 获取字符串的单个字符值
func (s *String) GetByte(index int) byte {
	return s.buf[index]
}

// GetStr 获取字符串的某个片段 返回String
func (s *String) GetStr(index, end int) string {
	return string(s.buf[index:end])
}

// GetStrString 获取字符串的某个片段 返回String结构
func (s *String) GetStrString(index, end int) *String {
	return &String{buf: s.buf[index:end]}
}

func (s *String) RemoveIndexStr(lens int) {
	if lens > s.Len() {
		evan.ErrorLog(errors.New("RemoveLens>stringLens Please Check"))
		return
	}
	s.buf = s.buf[lens:]
}

func (s *String) RemoveIndexStrByRune(lens int) {
	runes := s.runes
	if lens > len(runes) {
		evan.ErrorLog(errors.New("RemoveLens>stringLens Please Check"))
		return
	}
	s.buf = RunesToBytes(runes[lens:])
}

// CheckIsNull 检查字符串是否为空 只包含' '与'\t'与'\n'都会被视为不合法的值
func (s *String) CheckIsNull() bool {
	for _, b := range s.buf {
		if b != 32 && b != 9 && b != 10 {
			return false
		}
	}
	return true
}

func RunesToBytes(rune []rune) []byte {
	size := 0
	for _, r := range rune {
		size += utf8.RuneLen(r)
	}

	bs := make([]byte, size)

	count := 0
	for _, r := range rune {
		count += utf8.EncodeRune(bs[count:], r)
	}
	return bs
}
