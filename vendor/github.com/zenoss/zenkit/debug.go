package zenkit

import (
	"context"
	"fmt"
	"path"
	"regexp"
	"runtime"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	stripFnPreamble = regexp.MustCompile(`^.*\.(.*)$`)
)

// Logs entry/exit of function and func duration
// Usage: defer Trace(ctx)()
func Trace(ctx context.Context) func() {
	return newTracerWithDepth(2).Trace(ctx)
}

// NewTracer returns a Tracer that can be used to log entry/exit of a function
// as well as its duration.  Example uage:
//
// defer zenkit.GetTracer().WithField("my", "field").Trace(ctx)()
func NewTracer() *Tracer {
	return &Tracer{depth: 1}
}

func newTracerWithDepth(d int) *Tracer {
	return &Tracer{depth: d}
}

type Tracer struct {
	fields logrus.Fields
	depth  int
}

func (t *Tracer) WithField(key string, value interface{}) *Tracer {
	if t.fields == nil {
		t.fields = logrus.Fields{}
	}

	t.fields[key] = value

	return t
}

func (t *Tracer) Trace(ctx context.Context) func() {
	fnName := "<unknown>"
	// Skip this function and fetch the PC and file for the parent
	pc, file, line, ok := runtime.Caller(t.depth)
	log := ContextLogger(ctx).WithFields(t.fields).WithField("loc", fmt.Sprintf("%s:%d", path.Base(file), line))
	if ok {
		fnName = stripFnPreamble.ReplaceAllString(runtime.FuncForPC(pc).Name(), "$1")
	}

	// Enter
	begin := time.Now()
	log.Debugf("ENTER %s()", fnName)
	exit := func() {
		log.WithField("dur", time.Now().Sub(begin)).Debugf("EXIT %s()", fnName)

	}
	return exit
}
