package utils

import (
	"github.com/stretchr/testify/assert"
	"os"
	"sync"
	"syscall"
	"testing"
	"time"
)

func Test_ShutdownHookShouldExecuteAllHookWhenSIGTERMIsReceived(t *testing.T) {
	testWithSignal(t, syscall.SIGTERM)
}

func Test_ShutdownHookShouldExecuteAllHookWhenSIGKILLIsReceived(t *testing.T) {
	testWithSignal(t, syscall.SIGKILL)
}

func testWithSignal(t *testing.T, sig os.Signal) {
	defer clearHooks()
	defer startListening() // since the first test makes the listener go routine exit, we need to restart it when we're done

	timeout := time.Second * 10
	sigDone := make(chan struct{})

	executions := sync.WaitGroup{}
	executions.Add(2)

	hook := NewSignalHook(sig, func() {
		executions.Done()
	})

	RegisterShutdownHook(hook)
	RegisterShutdownHook(hook)

	signalChannel <- sig

	go func() {
		executions.Wait()
		sigDone <- struct{}{}
	}()

	select {
	case <-sigDone:
	case <-time.After(timeout):
		assert.FailNow(t, "Timeout waiting for hooks to execute")
	}
}
