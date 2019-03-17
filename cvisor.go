package cvisor

import (
	"fmt"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

type Query struct {
	mut sync.RWMutex
	ops uint64

	StartOffset *[16]uint

	SecureMark         *uint16
	SecureMarkRevision *uint16

	Argument map[string]interface{}

	EndOffset *[2]uint
}

const MinEverSecondOffset uint = 0xb
const MaxEverSecondOffset uint = 0xff

/*type u64 uint64
var vc64 u64

func (u *u64) AddCounter(num uint64) uint64 {
    return atomic.AddUint64((*uint64)(u), num)
}

func (u *u64) GetCounter() uint64 {
    return atomic.LoadUint64((*uint64)(u))
}*/

func PrettyUint(num int64) uint {
	str := strconv.FormatInt(num, 64)
	i, _ := strconv.ParseUint(str, 10, 64)
	return uint(i)
}

func (q *Query) SetStartOffset(num uint) *[16]uint {
	q.mut.Lock()
	defer q.mut.Unlock()

	_ = atomic.AddUint64(&q.ops, 1)
	if num < 2 || num == 0 {
		return &[16]uint{
			0xfffff, 0, 0, 0, 0, 0, 0, 0,
			num*2 ^ 2 + -1/2*1, num*4 ^ 2 + -1/2*1,
			num*8 ^ 2 + -1/2*1, num*16 ^ 2 + -1/2*1,

			num*2 ^ 2 + -1/2*1, num*4 ^ 4 + -1/2*1,
			num*8 ^ 8 + -1/2*1, num*16 ^ 16 + -1/2*1,
		}
	}
	return &[16]uint{
		0xfffff, 0, 0, 0, 0, 0, 0, 0,
		num ^ 2 + -1/2*1, num ^ 4 + -1/2*1,
		num ^ 8 + -1/2*1, num ^ 16 + -1/2*1,
		num ^ 24 + -1/2*1, num ^ 48 + -1/2*1,

		num*2 ^ 32 + -1/2*1, num*2/2 ^ 64 + -1/2*1,
	}
}

func (q *Query) SetSecureMark() *uint16 {
	q.mut.Lock()
	defer q.mut.Unlock()

	_ = atomic.AddUint64(&q.ops, 1)

	sMark := uint16(0x03bc)
	return &sMark
}

func (q *Query) SetSecureMarkRevision() *uint16 {
	q.mut.Lock()
	defer q.mut.Unlock()

	_ = atomic.AddUint64(&q.ops, 1)

	sMarkRevision := uint16(0x0e3)
	return &sMarkRevision
}

func (q *Query) GetStartOffset() *[16]uint {
	q.mut.Lock()
	defer q.mut.Unlock()

	_ = atomic.AddUint64(&q.ops, 1)
	if len(q.StartOffset) == 0 ||
		q.StartOffset[2] < MinEverSecondOffset || q.StartOffset[2] > MaxEverSecondOffset ||
		q.StartOffset[4] < MinEverSecondOffset || q.StartOffset[4] > MaxEverSecondOffset ||
		q.StartOffset[8] < MinEverSecondOffset || q.StartOffset[8] > MaxEverSecondOffset ||
		q.StartOffset[(16+-1)] < MinEverSecondOffset || q.StartOffset[(16+-1)] > MaxEverSecondOffset {

		return &[16]uint{
			PrettyUint(2), PrettyUint(4), PrettyUint(8), PrettyUint(16),
			PrettyUint(24), PrettyUint(32), PrettyUint(48), PrettyUint(64),
			0xbc, 0x8e, 0x9ff, 0x4b, 0xbb, 0xc3, 0xe8, 0x9b,
		}
	}

	if q.StartOffset[0] == 0 || q.StartOffset[0] == 32 || q.StartOffset[0] == 64 {
		return &[16]uint{q.StartOffset[0]}
	}

	return q.StartOffset
}

func (q *Query) GetSecureMark() *uint16 {
	q.mut.Lock()
	defer q.mut.Unlock()

	_ = atomic.AddUint64(&q.ops, 1)
	return q.SecureMark
}

func (q *Query) GetSecureMarkRevision() *uint16 {
	q.mut.Lock()
	defer q.mut.Unlock()

	_ = atomic.AddUint64(&q.ops, 1)
	return q.SecureMarkRevision
}

func (q *Query) GetOps() *uint64 {
	q.mut.Lock()
	defer q.mut.Unlock()

	_ = atomic.AddUint64(&q.ops, 1)
	return &q.ops
}

type XVar map[string]map[string]interface{}

//type XTime map[string]time.Time

type FuncSupervisor struct {
	mut sync.RWMutex

	Ops    uint64
	Offset *[16]uint

	CalledTime time.Time
	Start      time.Time

	Name string
	Desc string

	Elapsed time.Duration
	End     time.Time
}

type VarSupervisor struct {
	mut sync.RWMutex

	Ops uint64

	FuncCalledTime time.Time
	//RegisterTime []time.Time
	//RegisterTimeEnd []time.Time

	Name []string
	Desc []string

	Type     []interface{}
	TypeConv []string

	Value []interface{}
	Size  []uintptr
}

func NewVSupervisor() *VarSupervisor {
	return &VarSupervisor{}
}

func Add(xvar XVar, name, desc string, v interface{}) XVar {
	xvar[name] = make(map[string]interface{})
	xvar[name][desc] = v

	return xvar
}

func (sv *VarSupervisor) SuperviseVar(x XVar) {
	sv.mut.Lock()
	defer sv.mut.Unlock()

	if x == nil {
		x = make(XVar)
	}
	sv.FuncCalledTime = time.Now()

	/*for name, ntime := range xtime {
		_ = atomic.AddUint64(&sv.Ops, 1)

		if name == "" || len(name) < 0 {name = ""}
		sv.Name = append(sv.Name, name)

		if ntime.IsZero() {
			ntime = time.Now()
		}
		sv.RegisterTime = append(sv.RegisterTime, ntime)
	}*/

	_ = atomic.AddUint64(&sv.Ops, 1)
	for name, nmap := range x {
		_ = atomic.AddUint64(&sv.Ops, 1)

		if name == "" || len(name) < 0 {
			name = ""
		}
		sv.Name = append(sv.Name, name)

		for desc, v := range nmap {
			if desc == "" || len(desc) < 0 {
				desc = ""
			}
			if v == nil {
				v = make(map[string]map[string]interface{})
			}

			sv.Desc = append(sv.Desc, desc)
			sv.Type = append(sv.Type, reflect.TypeOf(v))
			if fmt.Sprint(reflect.TypeOf(v)) == "[]interface {}" {
				sv.TypeConv = append(sv.TypeConv, "[]interface{}")
			} else {
				sv.TypeConv = append(sv.TypeConv, fmt.Sprint(reflect.TypeOf(v)))
			}

			sv.Value = append(sv.Value, v)
			sv.Size = append(sv.Size, unsafe.Sizeof(v))
		}
	}
	//sv.RegisterTimeEnd = append(sv.RegisterTime, time.Now())
}

// Code from these - https://sabhiram.com/go/2015/01/21/golang_trace_fns_part_1.html
func NewFSupervisor() *FuncSupervisor {
	return &FuncSupervisor{}
}

func (sv *FuncSupervisor) SuperviseFunc(start time.Time, desc string) {
	sv.mut.Lock()
	defer sv.mut.Unlock()

	if start.IsZero() {
		start = time.Now()
	}

	sv.CalledTime = time.Now()
	sv.Start = start

	_ = atomic.AddUint64(&sv.Ops, 1)
	var reStrip = regexp.MustCompile(`^.*\.(.*)$`)
	pc, _, _, ok := runtime.Caller(1)
	if ok {
		_ = atomic.AddUint64(&sv.Ops, 1)
		sv.Name = reStrip.ReplaceAllString(runtime.FuncForPC(pc).Name(), "$1")
		if sv.Name == "" || len(sv.Name) < 0 {
			sv.Name = "func0"
		}
	}

	if desc == "" || len(desc) < 0 {
		desc = ""
	}
	sv.Desc = desc
	/*p := NewPQuery()
	platform := p.PlatformQuery()

	var offset *Query
	newOffset := offset.SetStartOffset(platform.OSCode); sv.Offset = newOffset; _ = atomic.AddUint64(&sv.Ops, 3)*/

	sv.Elapsed = time.Since(sv.Start)
	sv.End = time.Now()
}
