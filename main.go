package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"syscall"
)

type procStatus struct {
	pid  int
	exit bool
}

func main() {
	fmt.Printf("Run %v\n", os.Args[1:])

	cmd := exec.Command(os.Args[1], os.Args[2:]...)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Ptrace: true,
	}

	cmd.Start()
	err := cmd.Wait()
	if err != nil {
		fmt.Printf("Wait returned: %v\n", err)
	}

	sc, err := traceSyscalls(cmd.Process.Pid)
	sc.print()
}

func traceSyscalls(pid int) (syscallCounter, error) {
	var sc syscallCounter
	var regs syscall.PtraceRegs
	wpid := pid
	sc = sc.init()
	pids := []procStatus{
		procStatus{wpid, true},
	}
	curIdx := 0
	exitPid := false
	newPid := false

	// Trace child processes
	err := syscall.PtraceSetOptions(pid, syscall.PTRACE_O_TRACECLONE)
	if err != nil {
		return nil, err
	}
	err = syscall.PtraceSetOptions(pid, syscall.PTRACE_O_TRACEFORK)
	if err != nil {
		return nil, err
	}

	for {
		if newPid == false {
			err := syscall.PtraceGetRegs(wpid, &regs)
			if err != nil {
				// Remove PID
				i := indexOf(pids, wpid)
				if i == 0 {
					pids = pids[i+1:]
				} else {
					pids = append(pids[:i], pids[i+1:]...)
				}

				if len(pids) == 0 {
					break
				} else {
					exitPid = true
				}
			}

			if exitPid == false {
				// Count only one "occurrence"
				if pids[curIdx].exit {
					// Uncomment to print syscalls
					//name := ss.getName(regs.Orig_rax)
					//fmt.Printf("[PID %d]: %s\n", wpid, name)

					sc.inc(regs.Orig_rax)
				}

				err = syscall.PtraceSyscall(wpid, 0)
				if err != nil {
					return nil, err
				}
			} else {
				exitPid = false
			}
		} else {
			newPid = false
		}

		wpidTmp, err := syscall.Wait4(-1, nil, 0, nil)
		wpid = wpidTmp
		if err != nil {
			return nil, err
		}

		// Check for new PID
		curIdx = indexOf(pids, wpid)
		if curIdx < 0 {
			pids = append(pids, procStatus{wpid, false})
			curIdx = len(pids) - 1
			newPid = true

			err = syscall.PtraceSyscall(wpid, 0)
			if err != nil {
				return nil, err
			}

			fmt.Println("Attach to " + strconv.Itoa(wpid))
		}

		pids[curIdx].exit = !pids[curIdx].exit
	}

	return sc, nil
}

func indexOf(s []procStatus, e int) int {
	for i, a := range s {
		if a.pid == e {
			return i
		}
	}
	return -1
}
