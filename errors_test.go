package erx

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConditionErr(t *testing.T) {
	assert.Nil(t, ConditionErr(false, "condition not met"))
	assert.EqualError(t, ConditionErr(true, "condition met"), "condition met")
}

func TestAssertErr(t *testing.T) {
	assert.Nil(t, AssertErr(true, "assert ok"))
	assert.EqualError(t, AssertErr(false, "assert failed"), "assert failed")
}

func TestResultErr(t *testing.T) {
	assert.Nil(t, ResultErr("hello", nil))
	assert.EqualError(t, ResultErr("world", fmt.Errorf("error")), "error")
}
