/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-06-05 13:35:59
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-08-05 13:15:55
 * @FilePath: \go-toolbox\tests\sign_message_test.go
 * @Description: 签名客户端测试，公共参数提取，结合自定义 WaitGroup 并发测试
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/kamalyes/go-toolbox/pkg/sign"
	"github.com/stretchr/testify/assert"
)

type TestPayload struct {
	UserID   int
	Username string
	Email    string
}

// 公共测试参数辅助函数

func newTestClient(t *testing.T) *sign.SignerClient[TestPayload] {
	secretKey := []byte("my_secret_key_123")
	serializer := sign.JSONSerializer{}

	client := sign.NewSignerClient[TestPayload]().
		WithSecretKey(secretKey).
		WithSerializer(serializer)

	client, err := client.WithAlgorithm(sign.AlgorithmSHA256)
	assert.NoError(t, err)

	return client
}

func newTestPayload() TestPayload {
	return TestPayload{
		UserID:   42,
		Username: "kamalyes",
		Email:    "kamalyes@example.com",
	}
}

// 测试用例

func TestSignerClient_ChainUsage(t *testing.T) {
	client := newTestClient(t)
	payload := newTestPayload()

	signedStr, err := client.Create(payload)
	assert.NoError(t, err)
	assert.NotEmpty(t, signedStr)

	gotPayload, valid, err := client.Validate(signedStr)
	assert.NoError(t, err)
	assert.True(t, valid)
	assert.NotNil(t, gotPayload)
	assert.Equal(t, payload.UserID, gotPayload.ExtraData.UserID)
	assert.Equal(t, payload.Username, gotPayload.ExtraData.Username)
	assert.Equal(t, payload.Email, gotPayload.ExtraData.Email)
	assert.NotEmpty(t, gotPayload.Send)
	assert.WithinDuration(t, time.UnixMicro(gotPayload.GenUnixMicro), time.Now(), time.Minute)

	// 测试错误算法设置
	_, err = client.WithAlgorithm("UNSUPPORTED-ALG")
	assert.Error(t, err)
}

func TestSignerClient_ValidateErrors(t *testing.T) {
	secretKey := []byte("my_secret_key_123")
	client := sign.NewSignerClient[TestPayload]().
		WithSecretKey(secretKey)

	// 不设置算法，调用Create应报错
	_, err := client.Create(TestPayload{})
	assert.Error(t, err)

	// 设置算法后创建成功
	client, err = client.WithAlgorithm(sign.AlgorithmSHA256)
	assert.NoError(t, err)
	signedStr, err := client.Create(TestPayload{UserID: 1})
	assert.NoError(t, err)

	// 篡改签名，验证失败
	tampered := signedStr + "tamper"
	_, valid, err := client.Validate(tampered)
	assert.Error(t, err)
	assert.False(t, valid)

	// 空字符串验证失败
	_, valid, err = client.Validate("")
	assert.Error(t, err)
	assert.False(t, valid)

	// 格式错误，缺少点分割
	_, valid, err = client.Validate("invalidstringwithoutdot")
	assert.Error(t, err)
	assert.False(t, valid)
}

func TestSignerClient_ConcurrentUsage_WithSync(t *testing.T) {
	client := newTestClient(t)
	payload := newTestPayload()

	const concurrency = 50
	const iterations = 100

	var wg sync.WaitGroup
	var mu sync.Mutex
	var firstErr error // 用于存储第一个发生的错误

	for i := 0; i < concurrency; i++ {
		wg.Add(1) // 增加等待计数
		go func() {
			defer wg.Done() // 确保在 goroutine 完成时调用 Done

			for j := 0; j < iterations; j++ {
				signedStr, err := client.Create(payload)
				if err != nil {
					mu.Lock()
					if firstErr == nil { // 只记录第一个错误
						firstErr = err
					}
					mu.Unlock()
					return
				}

				gotPayload, valid, err := client.Validate(signedStr)
				if err != nil || !valid {
					mu.Lock()
					if firstErr == nil { // 只记录第一个错误
						firstErr = fmt.Errorf("validate failed: %v, valid=%v", err, valid)
					}
					mu.Unlock()
					return
				}

				if gotPayload.ExtraData.UserID != payload.UserID {
					mu.Lock()
					if firstErr == nil { // 只记录第一个错误
						firstErr = fmt.Errorf("payload mismatch")
					}
					mu.Unlock()
					return
				}
			}
		}()
	}

	wg.Wait() // 等待所有 goroutine 完成
	assert.NoError(t, firstErr)
}
