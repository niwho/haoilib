package haoilib

import (
	"unicode/utf16"
	//"unicode/utf8"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	//Haoib Haoilib
	Hb HB
)

func GbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

func Utf8ToGbk(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewEncoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

type Haoilib struct {
	//dll        *windows.DLL
	lazydll *windows.LazyDLL
	mdll    *windows.LazyDLL
	//getpoint   *windows.Proc
	getPointW   *windows.LazyProc
	getPoint    *windows.LazyProc
	sendByteEx  *windows.LazyProc
	sendByteExW *windows.LazyProc
	getBusyW    *windows.LazyProc
	setRebateW  *windows.LazyProc
}

func sumFileMd5(name string) string {
	f, err := os.Open(name)
	if err != nil {
		panic(err)
	}
	md5 := md5.New()
	io.Copy(md5, f)
	return hex.EncodeToString(md5.Sum(nil))
}
func init() {
	if sumFileMd5("haoi.dll") != "9a4f6df9d212b739cdf07f5d9ea644ae" {
		log.Println("fatal")
		panic("[haoi.dll] broken")
	}
	Haoib := Haoilib{}
	Haoib.lazydll = windows.NewLazyDLL("haoi.dll")
	Haoib.getPoint = Haoib.lazydll.NewProc("GetPoint")
	Haoib.getPointW = Haoib.lazydll.NewProc("GetPointW")
	Haoib.sendByteEx = Haoib.lazydll.NewProc("SendByteEx")
	Haoib.sendByteExW = Haoib.lazydll.NewProc("SendByteExW")
	Haoib.getBusyW = Haoib.lazydll.NewProc("GetBusyW")
	Haoib.setRebateW = Haoib.lazydll.NewProc("SetRebateW")
	Hb = &Haoib

}

func (h *Haoilib) SetRebateW(softkey string) {
	softkey1 := utf16.Encode([]rune(softkey + "\x00"))
	_, _, errno := h.setRebateW.Call(uintptr(unsafe.Pointer(&softkey1[0])))

	if errno.(syscall.Errno) != 0 {
		panic(errno)
	}
}

func (h *Haoilib) SendByteEx(MyUserStr, GameID string, PicBuffer []byte, Size, TimeOut, LostPoint int64,
	BeiZhu string) (Result, Reply string, err error) {

	MyUserStr1, _ := Utf8ToGbk([]byte(MyUserStr + "\x00"))
	GameID1, _ := Utf8ToGbk([]byte(GameID + "\x00"))
	BeiZhu1, _ := Utf8ToGbk([]byte(BeiZhu + "\x00"))
	result := make([]byte, 512)
	reply := make([]byte, 512)
	log.Println("PicBuffer", len(PicBuffer))
	log.Println(len(MyUserStr1))
	log.Println(len(GameID1))
	log.Println(len(BeiZhu1))
	_, _, errno := h.sendByteEx.Call(
		uintptr(unsafe.Pointer(&(MyUserStr1[0]))),
		uintptr(unsafe.Pointer(&(GameID1[0]))),
		uintptr(unsafe.Pointer(&(PicBuffer[0]))),
		uintptr(Size),
		uintptr(TimeOut),
		uintptr(LostPoint),
		uintptr(unsafe.Pointer(&(BeiZhu1[0]))),
		uintptr(unsafe.Pointer(&(result[0]))),
		uintptr(unsafe.Pointer(&(reply[0]))))
	if errno.(syscall.Errno) != 0 {
		err = errno
		return
	}

	Result = string((result[:h.validByteLen(result)]))
	Reply = string((reply[:h.validByteLen(reply)]))
	return
}

func (h *Haoilib) SendByteExW(MyUserStr, GameID string, PicBuffer []byte, Size, TimeOut, LostPoint int64,
	BeiZhu string) (Result, Reply string, err error) {

	//MyUserStr1 := utf16.Encode([]rune(MyUserStr + "\x00"))
	//GameID1 := utf16.Encode([]rune(GameID + "\x00"))
	//BeiZhu1 := utf16.Encode([]rune(BeiZhu + "\x00"))
	MyUserStr1 := utf16.Encode([]rune(MyUserStr))
	GameID1 := utf16.Encode([]rune(GameID))
	BeiZhu1 := utf16.Encode([]rune(BeiZhu))
	result := make([]uint16, 512)
	reply := make([]uint16, 512)
	log.Println("PicBuffer", len(PicBuffer))
	log.Println(len(MyUserStr1))
	log.Println(len(GameID1))
	log.Println(len(BeiZhu1))

	_, _, errno := h.sendByteExW.Call(
		uintptr(unsafe.Pointer(&MyUserStr1[0])),
		uintptr(unsafe.Pointer(&GameID1[0])),
		uintptr(unsafe.Pointer(&(PicBuffer[0]))),
		uintptr(Size),
		uintptr(TimeOut),
		uintptr(LostPoint),
		uintptr(unsafe.Pointer(&BeiZhu1[0])),
		uintptr(unsafe.Pointer(&(result[0]))),
		uintptr(unsafe.Pointer(&(reply[0]))))
	if errno.(syscall.Errno) != 0 {
		err = errno
		return
	}

	Result = string((utf16.Decode(result[:h.validLen(result)])))
	Reply = string((utf16.Decode(reply[:h.validLen(reply)])))
	return
}

func (h *Haoilib) validByteLen(ar []byte) (valid_len int) {
	for _, c := range ar {
		if c > 0 {
			valid_len += 1
			continue
		}
		break
	}
	return
}

func (h *Haoilib) validLen(ar []uint16) (valid_len int) {
	for _, c := range ar {
		if c > 0 {
			valid_len += 1
			continue
		}
		break
	}
	return
}
func (h *Haoilib) GetPoint(MyUserStr string) (int, error) {
	a, _ := Utf8ToGbk([]byte(MyUserStr))
	rtxx := make([]byte, 512)
	r1, _, errno := h.getPoint.Call(uintptr(unsafe.Pointer(&a[0])), uintptr(unsafe.Pointer(&((rtxx)[0]))))
	if errno.(syscall.Errno) != 0 {
		return -1, errno
	}
	valid_len := 0
	for _, c := range rtxx {
		if c > 0 {
			valid_len += 1
			continue
		}
		break
	}
	astr := string((rtxx[:valid_len]))
	if int(r1) != 4 {
		return -1, errors.New(astr)
	}
	val, err := strconv.Atoi(astr)
	if err != nil {
		return -1, err
	}
	return val, nil
}
func (h *Haoilib) GetPointW(MyUserStr string) (int, error) {
	a := utf16.Encode([]rune(MyUserStr))
	rtxx := make([]uint16, 512)
	r1, _, errno := h.getPointW.Call(uintptr(unsafe.Pointer(&a[0])), uintptr(unsafe.Pointer(&((rtxx)[0]))))
	if errno.(syscall.Errno) != 0 {
		return -1, errno
	}
	valid_len := 0
	for _, c := range rtxx {
		if c > 0 {
			valid_len += 1
			continue
		}
		break
	}
	astr := string((utf16.Decode(rtxx[:valid_len])))
	if int(r1) != 4 {
		return -1, errors.New(astr)
	}
	val, err := strconv.Atoi(astr)
	if err != nil {
		return -1, err
	}
	return val, nil
}

func (h *Haoilib) GetBusyW() (string, error) {
	result := make([]uint16, 512)
	_, _, errno := h.getBusyW.Call(uintptr(unsafe.Pointer(&result[0])))

	if errno.(syscall.Errno) != 0 {
		return "", errno
	}
	return string((utf16.Decode(result[:h.validLen(result)]))), nil
}
