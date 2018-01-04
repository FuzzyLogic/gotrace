# gotrace

A simple system call tracer written in Go. It also supports tracing any processes/threads that are forked/cloned from the original tracee.
This is particularly useful when tracing bash scripts or more complex processes.

This code is based on the strace-from-scratch code by Liz Rice. However, strace-from-scratch does not support tracing any forked/cloned processes, making it less useful.

## The original strace-from-scratch repo
The strace-from-scratch code was shown at Gophercon 2017. [Here's a walkthrough of this code](https://medium.com/@lizrice/strace-in-60-lines-of-go-b4b76e3ecd64) and [here's the slide deck](https://speakerdeck.com/lizrice/a-go-programmers-guide-to-syscalls). 

The repository can be found [here](https://github.com/lizrice/strace-from-scratch)
