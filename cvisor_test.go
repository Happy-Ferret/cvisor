package cvisor 

import (
	"fmt"
	"time"
	"testing"
)

type MyStruct struct {
	A string 
	B int 
}

func TestSuperviseFunc(t *testing.T) {
	_func := func() (f *FuncSupervisor) {
		f = NewFSupervisor() // or v = new(cvisor.FuncSupervisor)
		defer f.SuperviseFunc(time.Now(), "Simple description for func1")
		time.Sleep(100 *time.Millisecond)

		return
	}

	f := _func()
	fmt.Printf("Operations per seconds: %d\nCalled time: %s\nStart time: %s\nName of function: %s\nDescription for function: %s\nElapsed time: %s\nEnd time: %s\n\n\r",
		f.Ops, f.CalledTime, f.Start, f.Name, f.Desc, f.Elapsed, f.End,
	)
}

func TestSuperviseFuncName(t *testing.T) {
	_func := func() (f *FuncSupervisor) {
		f = NewFSupervisor() // or v = new(cvisor.FuncSupervisor)
		defer f.SuperviseFunc(time.Now(), "Simple description for func1")
		time.Sleep(100 *time.Millisecond)

		return
	}

	f := _func()
	if f.Name == "" || len(f.Name) < 0 {
		t.Fatal("Length of 'f.Name' is nil o_O")
	}
}

func TestSuperviseFuncTime(t *testing.T) {
	_func := func() (f *FuncSupervisor) {
		f = NewFSupervisor() // or v = new(cvisor.FuncSupervisor)
		defer f.SuperviseFunc(time.Time{}, "Simple description for func1")
		time.Sleep(100 *time.Millisecond)

		return
	}

	f := _func()
	if f.Start.IsZero() {
		t.Fatal("Start time is zero")
	}
}

func TestSuperviseVar(t *testing.T) {
	_func := func() (v *VarSupervisor) {
		v = NewVSupervisor() // or v = new(cvisor.VarSupervisor)
		var register = make(XVar)

		testString := "test test test 1"; 
		register = Add(register, "testString", "Simple string", testString)

		time.Sleep(3 * time.Second)
		var testPointerString *string
		temp := "test test test 2"; testPointerString = &temp
		register = Add(register, "testPointerString", "Simple test string pointer", testPointerString)

		testInt := int(128)
		register = Add(register, "testInt", "Simple test integer", testInt)

		testUint := uint(0x3bc4f)
		register = Add(register, "testUint", "Simple test uint", testUint)

		testFloat := float32(0.128)
		register = Add(register, "testFloat", "Simple test float", testFloat)

		var testInterfaceStruct MyStruct
		testInterfaceStruct.A = "test"
		testInterfaceStruct.B = 100
		register = Add(register, "testInterfaceStruct", "Simple test interface structure", testInterfaceStruct)

		var test_v []interface{}
		test_v = append(test_v, "test string")
		test_v = append(test_v, 125)
		test_v = append(test_v, uint8(0xff))
		test_v = append(test_v, uint16(0x32b))
		test_v = append(test_v, uint32(0xfc452))
		test_v = append(test_v, uint64(0x8563bc2))
		test_v = append(test_v, uint(0x38553bbc14ccf3))
		test_v = append(test_v, uintptr(0x862b2736509375c))
		test_v = append(test_v, float32(0.1))
		test_v = append(test_v, float64(1726442.34753421))
		test_v = append(test_v, testInterfaceStruct)
		register = Add(register, "test_v", "Global test", test_v)
		defer v.SuperviseVar(register)

		return
	}
	
	v := _func()
	for i := range v.Name {
		fmt.Printf("[%d]Name of var: %s\n[%d]Type of var: %v\n[%d]Converted type of var: %s\n[%d]Variable value: %v\n[%d]Size of variable: %d\n\n\r",
			i, v.Name[i], i, v.Type[i], i, v.TypeConv[i], i, v.Value[i],i, v.Size[i],
		)
	}
}

func TestSuperviseVarNil(t *testing.T) {
	_func := func() (v *VarSupervisor) {
		v = NewVSupervisor()
		var register XVar
		defer v.SuperviseVar(register)

		return
	}

	// v:XVar can not be nil
	v := _func()
	if v == nil {
		t.Fatal("'v' is nil o_O")
	}
}
