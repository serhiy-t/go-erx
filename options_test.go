package erx

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func testLogger() *ErrLogger {
	var out error
	e := ErrorLogger(&out)

	out = fmt.Errorf("error-returned")
	e.suppressed = append(e.suppressed, fmt.Errorf("error-suppressed"))
	e.ignored = append(e.ignored, fmt.Errorf("error-ignored"))

	return e
}

func testLogFn(logs *[]string) LogFn {
	return func(logMessage *LogMessage) {
		result := ""
		if len(logMessage.Stack()) == 0 {
			result += "[no-stack]"
		}
		for _, tag := range logMessage.Tags {
			result += fmt.Sprintf("[%s]", tag)
		}
		result += " "
		result += fmt.Sprintf(logMessage.Format, logMessage.A...)
		*logs = append(*logs, result)
	}
}

func TestLogReturned(t *testing.T) {
	var logs []string
	defer SetLogFn(testLogFn(&logs)).ThenRestore()
	e := testLogger()

	e.LogReturnedError()
	assert.Equal(t, []string{
		"[erx][returned] error: error-returned",
	}, logs)
}

func TestLogReturned_NoError(t *testing.T) {
	var logs []string
	defer SetLogFn(testLogFn(&logs)).ThenRestore()
	var out error
	e := ErrorLogger(&out)

	e.LogReturnedError()
	assert.Empty(t, logs)
}

func TestLogSuppressed(t *testing.T) {
	var logs []string
	defer SetLogFn(testLogFn(&logs)).ThenRestore()
	e := testLogger()

	e.LogSuppressedErrors()
	assert.Equal(t, []string{
		"[erx][suppressed] error: error-suppressed",
	}, logs)
}

func TestLogIgnored(t *testing.T) {
	var logs []string
	defer SetLogFn(testLogFn(&logs)).ThenRestore()
	e := testLogger()

	e.LogIgnoredErrors()
	assert.Equal(t, []string{
		"[erx][ignored] error: error-ignored",
	}, logs)
}

func TestLogAll(t *testing.T) {
	var logs []string
	defer SetLogFn(testLogFn(&logs)).ThenRestore()
	e := testLogger()

	e.LogAllErrors()
	assert.Equal(t, []string{
		"[erx][returned] error: error-returned",
		"[erx][suppressed] error: error-suppressed",
		"[erx][ignored] error: error-ignored",
	}, logs)
}

func TestLogSilent(t *testing.T) {
	var logs []string
	defer SetLogFn(testLogFn(&logs)).ThenRestore()
	e := testLogger()

	e.LogSilentErrors()
	assert.Equal(t, []string{
		"[erx][suppressed] error: error-suppressed",
		"[erx][ignored] error: error-ignored",
	}, logs)
}
