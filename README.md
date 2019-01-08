# cvisor - Golang code supervisor

Library for function, variable track/supervising in high-level mode. 
Than is convenient and what is it for? 
* **Function:** Capture detailed data which transform to output structure, runtime.
* **Variables:** Capture start call time of function, register variables: name of var, description of var, type of var, converted type of var, value of var. 

# What is this for:
This is necessary to support timely, fast, reliable, registration of these functions and variables.

# How to use:
**Example of function testing**:
```go
func TestFunc() (f *cvisor.FuncSupervisor){
	f = cvisor.NewFSupervisor() // or f = new(cvisor.FuncSupervisor)
	defer f.SuperviseFunc(time.Now(), "Simple description for Test function")
	time.Sleep(100 *time.Millisecond)
	return
}

func main() {
	f := TestFunc()
	fmt.Printf("Name: %s\nDescription: %s\n\nCalled time: %s\nStart time: %s\nEnd time: %s\n\nElapsed time: %s\n\nOPS: %d\n\r",
		f.Name, f.Desc, f.CalledTime, f.Start, f.End, f.Elapsed, f.Ops,
	)
}
```
**Example of variable registration testing**:
```go
func TestVar() (v *cvisor.VarSupervisor){
	v = cvisor.NewVSupervisor() 
	register := make(cvisor.XVar)

	testVar := "hello"
	register = cvisor.Add(register, "testVar", "Simple test description for testVar", testVar)

	testVar2 := 122
	register = cvisor.Add(register, "testVar2", "Simple test description for testVar2", testVar2)

	defer v.SuperviseVar(register)
	return
}

func main() {
	v := TestVar()

	for i := range v.Name {
		fmt.Printf("[%d]Name: %s\n[%d]Description: %s\n\n[%d]Type: %v\n[%d]Converted type: %s\n[%d]Value: %v\n\n[%d]Size: %d\n\n",
			i, v.Name[i], i, v.Desc[i], i, v.Type[i], i, v.TypeConv[i], i, v.Value[i], i, v.Size[i], 
		)
	}
}
```
