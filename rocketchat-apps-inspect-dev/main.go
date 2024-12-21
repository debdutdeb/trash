package main

import (
	"flag"
	"fmt"
	"os"
	"syscall"
	"time"
)

type processStatus struct {
	pid       int
	registers syscall.PtraceRegs
}

func (p processStatus) IsWriting() bool {
	return p.registers.Orig_rax == syscall.SYS_WRITE
}

func (p processStatus) IsReading() bool {
	return p.registers.Orig_rax == syscall.SYS_READ
}

func (p processStatus) WriteLength() uint64 {
	return p.registers.Rdx // 3rd argument to write
}

func (p processStatus) WriteArgument() (string, error) {
	data := make([]byte, p.WriteLength())

	_, err := syscall.PtracePeekText(p.pid, uintptr(p.registers.Rsi), data) // rsi 2nd argument to write
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (p processStatus) WriteFileDescriptor() int {
	return int(p.registers.Rdi)
}

func panicf(format string, args ...interface{}) {
	panic(fmt.Sprintf(format, args...))
}

func main() {
	pid := flag.Int("pid", -1, "--pid")

	task := flag.Bool("task", false, "--task")

	flag.Parse()

	if task != nil && *task {
		fmt.Println(os.Getpid())
		for {
			fmt.Println("noop")
			time.Sleep(time.Second * 1)
		}

		return
	}

	if pid != nil && *pid == -1 {
		panic("no pid passed")
	}

	err := syscall.PtraceAttach(*pid)
	if err != nil {
		panicf("error attaching to process: %v", err)
	}

	var waitStatus syscall.WaitStatus
	_, err = syscall.Wait4(*pid, &waitStatus, 0, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to wait for process pid: %v", *pid)
		syscall.PtraceDetach(*pid)
		return
	}

	for {
		err = syscall.PtraceSyscall(*pid, 0)
		if err != nil {
			panicf("failed to continue process: %v", err)
		}

		_, err = syscall.Wait4(*pid, &waitStatus, 0, nil)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to wait for process pid: %v", *pid)
			syscall.PtraceDetach(*pid)
			return
		}

		if waitStatus.Exited() {
			panicf("process already exited")
		}

		var registers syscall.PtraceRegs
		err = syscall.PtraceGetRegs(*pid, &registers)
		if err != nil {
			panic(err)
		}

		p := processStatus{*pid, registers}

		if !p.IsWriting() {
			continue
		}

		arg, err := p.WriteArgument()
		if err != nil {
			fmt.Fprint(os.Stderr, err)
		}

		fmt.Println(p.WriteFileDescriptor(), arg)
	}
}
