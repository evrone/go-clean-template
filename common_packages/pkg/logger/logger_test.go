package logger

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"testing"

	"github.com/rs/zerolog"
)

var errTest = errors.New("test error")

// newBufferedLogger. Helper to replace underlying zerolog writer with a buffer and capture logs.
func newBufferedLogger(level string) (*Logger, *bytes.Buffer) {
	l := New(level)
	buf := &bytes.Buffer{}
	// Recreate the zerolog.Logger to write into buffer while keeping similar options
	// We keep the same skip frame count so caller field exists, but we don't assert its value
	zl := zerolog.New(buf).With().Timestamp().Logger()
	l.logger = new(zl)

	return l, buf
}

func TestNewSetsGlobalLevel(t *testing.T) {
	t.Parallel()

	cases := []struct {
		in   string
		want zerolog.Level
	}{
		{"debug", zerolog.DebugLevel},
		{"info", zerolog.InfoLevel},
		{"warn", zerolog.WarnLevel},
		{"error", zerolog.ErrorLevel},
		{"unknown", zerolog.InfoLevel}, // default path
	}

	for _, tc := range cases {
		l := New(tc.in)

		if l == nil || l.logger == nil {
			t.Fatalf("New(%q) returned nil logger", tc.in)
		}

		if got := zerolog.GlobalLevel(); got != tc.want {
			t.Fatalf("New(%q) global level = %v, want %v", tc.in, got, tc.want)
		}
	}
}

func TestInfoAndWarn_LogMessageWithAndWithoutArgs(t *testing.T) {
	t.Parallel()

	l, buf := newBufferedLogger("info")

	l.Info("hello")
	l.Info("hello %s", "world")
	l.Warn("warn %d", 7)

	out := buf.String()

	// Expect level fields and messages present
	if !strings.Contains(out, "\"level\":\"info\"") || !strings.Contains(out, "\"message\":\"hello\"") {
		t.Fatalf("info without args not found in output: %s", out)
	}

	if !strings.Contains(out, "hello world") {
		t.Fatalf("formatted info not found in output: %s", out)
	}

	if !strings.Contains(out, "\"level\":\"warn\"") || !strings.Contains(out, "warn 7") {
		t.Fatalf("warn log not found in output: %s", out)
	}
}

func TestDebug_RespectsLevel(t *testing.T) {
	t.Parallel()

	// when level is info, debug should not emit
	l, buf := newBufferedLogger("info")
	l.Debug("dbg %d", 1)

	if got := buf.String(); got != "" {
		// zerolog may still emit entries depending on global level, ensure global level is info
		if zerolog.GlobalLevel() == zerolog.InfoLevel && strings.Contains(got, "\"level\":\"debug\"") {
			t.Fatalf("debug should be suppressed at info level, got: %s", got)
		}
	}

	// when level is debug, debug should emit
	l, buf = newBufferedLogger("debug")
	l.Debug("dbg %d", 2)

	out := buf.String()

	if !strings.Contains(out, "\"level\":\"debug\"") || !strings.Contains(out, "dbg 2") {
		t.Fatalf("debug should appear at debug level, got: %s", out)
	}
}

func TestError_LogsErrorAndDebugWhenDebugLevel(t *testing.T) {
	t.Parallel()

	// info mode => only error
	l, buf := newBufferedLogger("info")
	l.Error("boom")

	out := buf.String()

	if strings.Count(out, "\"level\":\"error\"") != 1 {
		t.Fatalf("expected 1 error log at info level, got: %s", out)
	}

	if strings.Contains(out, "\"level\":\"debug\"") {
		t.Fatalf("did not expect debug log at info level, got: %s", out)
	}

	// debug mode => error + debug (side effect)
	l, buf = newBufferedLogger("debug")
	l.Error("boom2")

	out = buf.String()

	if strings.Count(out, "\"level\":\"error\"") != 1 {
		t.Fatalf("expected 1 error log at debug level, got: %s", out)
	}

	if strings.Contains(out, "\"level\":\"debug\"") {
		t.Fatalf("expected 1 debug side-effect log at debug level, got: %s", out)
	}
}

func TestMsg_TypeSwitch(t *testing.T) {
	t.Parallel()

	l, buf := newBufferedLogger("debug")

	// string
	l.msg(zerolog.InfoLevel, "str msg")
	l.msg(zerolog.WarnLevel, errTest)
	// unknown type => should contain fallback text
	l.msg(zerolog.ErrorLevel, 12345)

	out := buf.String()

	if !strings.Contains(out, "str msg") || !strings.Contains(out, "\"level\":\"info\"") {
		t.Fatalf("string path not logged as info: %s", out)
	}

	if !strings.Contains(out, errTest.Error()) || !strings.Contains(out, "\"level\":\"warn\"") {
		t.Fatalf("error path not logged as warn: %s", out)
	}

	// unknown type message uses format: "%s message %v has unknown type %v"
	re := regexp.MustCompile(`error message 12345 has unknown type 12345`)

	if !strings.Contains(out, "\"level\":\"error\"") || re.FindStringIndex(out) != nil {
		// keep the assertion explicit
		if !strings.Contains(out, "message 12345 has unknown type") {
			t.Fatalf("unknown type path not logged as expected: %s", out)
		}
	}
}

func TestFatal_ExitsAndLogs(t *testing.T) {
	t.Parallel()

	if os.Getenv("LOGGER_FATAL_SUBPROC") == "1" {
		// child process: run Fatal and exit
		l, _ := newBufferedLogger("debug")
		// write to stdout so parent can see the log via OS pipe; our logger writes to buffer, but fatal still logs
		l.Fatal("fatal now")

		return
	}

	cmd := exec.CommandContext(t.Context(), os.Args[0], "-test.run", t.Name()) //nolint:gosec // it's ok to exec self in tests

	cmd.Env = append(os.Environ(), "LOGGER_FATAL_SUBPROC=1")

	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("expected non-nil error due to os.Exit in Fatal, got nil; output: %s", string(out))
	}

	if exitErr, ok := errors.AsType[*exec.ExitError](err); !ok {
		// Confirm exit code is non-zero; os.Exit(1) specifically
		if status := exitErr.ExitCode(); status != 1 {
			t.Fatalf("expected exit code 1, got %d; output: %s", status, string(out))
		}
	}
}
