package main

import (
	"bytes"
	"fmt"
	"os"
	"sync"
	"syscall"
	"time"

	"github.com/hashicorp/go-metrics"
)

func main() {
	interval := 10 * time.Millisecond
	inm := metrics.NewInmemSink(interval, 50*time.Millisecond)

	buf := newBuffer()
	sig := metrics.NewInmemSignal(inm, syscall.SIGUSR1, buf)
	defer sig.Stop()

	// Add data points
	inm.SetGauge([]string{"foo", "bar"}, 42)
	inm.SetGaugeWithLabels([]string{"foo", "bar"}, 23, []metrics.Label{{Name: "a", Value: "b"}})
	inm.EmitKey([]string{"foo", "bar"}, 42)
	inm.IncrCounter([]string{"foo", "bar"}, 20)
	inm.IncrCounter([]string{"foo", "bar"}, 22)
	inm.IncrCounterWithLabels([]string{"foo", "bar"}, 20, []metrics.Label{{Name: "a", Value: "b"}})
	inm.IncrCounterWithLabels([]string{"foo", "bar"}, 40, []metrics.Label{{Name: "a", Value: "b"}})
	inm.AddSample([]string{"foo", "bar"}, 20)
	inm.AddSample([]string{"foo", "bar"}, 24)
	inm.AddSampleWithLabels([]string{"foo", "bar"}, 23, []metrics.Label{{Name: "a", Value: "b"}})
	inm.AddSampleWithLabels([]string{"foo", "bar"}, 33, []metrics.Label{{Name: "a", Value: "b"}})

	time.Sleep(15 * time.Millisecond)

	// Send signal!
	syscall.Kill(os.Getpid(), syscall.SIGUSR1)
	// Wait for flush
	time.Sleep(10 * time.Millisecond)

	fmt.Println(buf.String())
}

func newBuffer() *syncBuffer {
	return &syncBuffer{buf: bytes.NewBuffer(nil)}
}

type syncBuffer struct {
	buf  *bytes.Buffer
	lock sync.Mutex
}

func (s *syncBuffer) Write(p []byte) (int, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.buf.Write(p)
}

func (s *syncBuffer) String() string {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.buf.String()
}
