package conditions

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapCheckCondition_IsMatch_Keys(t *testing.T) {
	cond := mapCheckCondition{
		keys:   []string{"test1"},
		values: nil,
		match:  nil,
	}

	assert.True(t, cond.isMatch(map[string]string{"test1": "yes", "test0": "no"}))
	assert.False(t, cond.isMatch(map[string]string{"test0": "no"}))
	assert.False(t, cond.isMatch(map[string]string{"test0": "test1"}))
	assert.False(t, cond.isMatch(map[string]string{}))
}

func TestMapCheckCondition_IsMatch_Values(t *testing.T) {
	cond := mapCheckCondition{
		keys:   nil,
		values: []string{"yes"},
		match:  nil,
	}

	assert.True(t, cond.isMatch(map[string]string{"test0": "no", "test1": "yes"}))
	assert.False(t, cond.isMatch(map[string]string{"yes": "no"}))
	assert.False(t, cond.isMatch(map[string]string{"test0": "no"}))
	assert.False(t, cond.isMatch(map[string]string{}))
}

func TestMapCheckCondition_IsMatch_Match(t *testing.T) {
	cond := mapCheckCondition{
		keys:   nil,
		values: nil,
		match:  map[string]string{"test1": "yes"},
	}

	assert.True(t, cond.isMatch(map[string]string{"test1": "yes"}))
	assert.False(t, cond.isMatch(map[string]string{"yes": "test1"}))
	assert.False(t, cond.isMatch(map[string]string{"test1": "no"}))
	assert.False(t, cond.isMatch(map[string]string{"test0": "yes"}))
	assert.False(t, cond.isMatch(map[string]string{"test0": "no"}))
	assert.False(t, cond.isMatch(map[string]string{}))
}
