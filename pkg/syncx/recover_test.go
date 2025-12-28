/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-12-28 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-28 00:00:00 10:00:00
 * @FilePath: \go-toolbox\pkg\syncx\recover_test.go
 * @Description: recover 函数测试
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package syncx

import (
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSafeGo(t *testing.T) {
	t.Run("normal execution", func(t *testing.T) {
		executed := int32(0)
		SafeGo(func() {
			atomic.AddInt32(&executed, 1)
		}, nil)

		time.Sleep(50 * time.Millisecond)
		assert.Equal(t, int32(1), atomic.LoadInt32(&executed))
	})

	t.Run("with panic", func(t *testing.T) {
		panicCaught := int32(0)
		SafeGo(func() {
			panic("test panic")
		}, func(r interface{}) {
			atomic.AddInt32(&panicCaught, 1)
		})

		time.Sleep(50 * time.Millisecond)
		assert.Equal(t, int32(1), atomic.LoadInt32(&panicCaught))
	})
}

func TestRecoverWithHandler(t *testing.T) {
	t.Run("with panic", func(t *testing.T) {
		handled := false
		func() {
			defer RecoverWithHandler(func(r interface{}) {
				handled = true
			})
			panic("test")
		}()
		assert.True(t, handled)
	})

	t.Run("without panic", func(t *testing.T) {
		handled := false
		func() {
			defer RecoverWithHandler(func(r interface{}) {
				handled = true
			})
		}()
		assert.False(t, handled)
	})

	t.Run("nil handler", func(t *testing.T) {
		assert.NotPanics(t, func() {
			defer RecoverWithHandler(nil)
			panic("test")
		})
	})
}

func TestRecover(t *testing.T) {
	assert.NotPanics(t, func() {
		defer Recover()
		panic("test")
	})
}

func TestMustRecover(t *testing.T) {
	t.Run("panic with handler", func(t *testing.T) {
		handled := false
		assert.Panics(t, func() {
			defer MustRecover(func(r interface{}) {
				handled = true
			})
			panic("test")
		})
		assert.True(t, handled)
	})

	t.Run("no panic", func(t *testing.T) {
		handled := false
		assert.NotPanics(t, func() {
			defer MustRecover(func(r interface{}) {
				handled = true
			})
		})
		assert.False(t, handled)
	})
}

func TestRecoverToError(t *testing.T) {
	t.Run("panic with error", func(t *testing.T) {
		testErr := errors.New("test error")
		fn := func() (err error) {
			defer RecoverToError(&err, nil)
			panic(testErr)
		}
		err := fn()
		assert.Error(t, err)
		assert.Equal(t, testErr, err)
	})

	t.Run("panic with string", func(t *testing.T) {
		fn := func() (err error) {
			defer RecoverToError(&err, nil)
			panic("test panic")
		}
		err := fn()
		assert.Error(t, err)
		assert.Equal(t, "test panic", err.Error())
	})

	t.Run("panic with other type", func(t *testing.T) {
		fn := func() (err error) {
			defer RecoverToError(&err, nil)
			panic(123)
		}
		err := fn()
		assert.Error(t, err)
	})

	t.Run("no panic", func(t *testing.T) {
		fn := func() (err error) {
			defer RecoverToError(&err, nil)
			return nil
		}
		err := fn()
		assert.NoError(t, err)
	})

	t.Run("with handler", func(t *testing.T) {
		handled := false
		fn := func() (err error) {
			defer RecoverToError(&err, func(r interface{}) {
				handled = true
			})
			panic("test")
		}
		err := fn()
		assert.Error(t, err)
		assert.True(t, handled)
	})
}

func TestFormatPanic(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{"error type", errors.New("test error"), "test error"},
		{"string type", "test string", "test string"},
		{"other type", 123, "panic occurred"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatPanic(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
