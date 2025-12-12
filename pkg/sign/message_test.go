/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-06-05 13:35:59
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-09-16 18:15:50
 * @FilePath: \go-toolbox\pkg\sign\message_test.go
 * @Description: 签名客户端测试，公共参数提取，结合自定义 WaitGroup 并发测试
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package sign

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type TestPayload struct {
	UserID   int
	Username string
	Email    string
}

// 公共测试参数辅助函数

func newTestClient(t *testing.T) *SignerClient[TestPayload] {
	secretKey := []byte("my_secret_key_123")
	serializer := JSONSerializer{}

	client := NewSignerClient[TestPayload]().
		WithSecretKey(secretKey).
		WithSerializer(serializer)

	client, err := client.WithAlgorithm(AlgorithmSHA256)
	assert.NoError(t, err)

	return client
}

func newBenchmarkClient(b *testing.B) *SignerClient[TestPayload] {
	secretKey := []byte("my_secret_key_123")
	serializer := JSONSerializer{}

	client := NewSignerClient[TestPayload]().
		WithSecretKey(secretKey).
		WithSerializer(serializer)

	client, err := client.WithAlgorithm(AlgorithmSHA256)
	assert.NoError(b, err)
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

func TestSignerClientChainUsage(t *testing.T) {
	issuer := "issuer string"
	client := newTestClient(t).WithIssuer(issuer)
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
	assert.Equal(t, issuer, gotPayload.Header.Issuer)
	assert.NotEmpty(t, gotPayload.Header.Send)
	assert.WithinDuration(t, time.UnixMicro(gotPayload.Header.IssuedAt), time.Now(), time.Minute)

	// 测试错误算法设置
	_, err = client.WithAlgorithm("UNSUPPORTED-ALG")
	assert.Error(t, err)
}

func TestValidateValidSignature(t *testing.T) {
	client := newTestClient(t)
	payload := newTestPayload()
	signedMessage, err := client.WithExpiration(1 * time.Second).Create(payload)
	assert.NoError(t, err, "应该成功创建签名消息")

	// 等待一段时间使签名过期
	time.Sleep(2 * time.Second) // 暂停2秒以确保过期

	// 验证签名消息
	validatedMessage, valid, err := client.Validate(signedMessage)
	assert.False(t, valid, "签名消息应该过期")
	assert.Error(t, err, "应该返回过期错误")
	assert.Nil(t, validatedMessage, "过期的签名消息应该返回 nil")
}

func TestSignerClientValidateErrors(t *testing.T) {
	secretKey := []byte("my_secret_key_123")
	client := NewSignerClient[TestPayload]().
		WithSecretKey(secretKey)

	// 不设置算法，调用Create应报错
	_, err := client.Create(TestPayload{})
	assert.Error(t, err)

	// 设置算法后创建成功
	client, err = client.WithAlgorithm(AlgorithmSHA256)
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

func BenchmarkSignerClientCreate(b *testing.B) {
	client := newBenchmarkClient(b)
	payload := newTestPayload()
	b.ResetTimer() // 重置计时器，以确保不测量设置时间
	for i := 0; i < b.N; i++ {
		_, err := client.Create(payload)
		if err != nil {
			b.Fatalf("failed to create signed message: %v", err)
		}
	}
}

func BenchmarkSignerClientValidate(b *testing.B) {
	client := newBenchmarkClient(b)
	payload := newTestPayload()
	signedStr, err := client.Create(payload)
	if err != nil {
		b.Fatalf("failed to create signed message: %v", err)
	}

	b.ResetTimer() // 重置计时器
	for i := 0; i < b.N; i++ {
		_, valid, err := client.Validate(signedStr)
		if err != nil || !valid {
			b.Fatalf("failed to validate signed message: %v, valid=%v", err, valid)
		}
	}
}

func BenchmarkSignerClientConcurrentCreate(b *testing.B) {
	client := newBenchmarkClient(b)
	payload := newTestPayload()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := client.Create(payload)
			if err != nil {
				b.Fatalf("failed to create signed message: %v", err)
			}
		}
	})
}

func BenchmarkSignerClientConcurrentValidate(b *testing.B) {
	client := newBenchmarkClient(b)
	payload := newTestPayload()
	signedStr, err := client.Create(payload)
	if err != nil {
		b.Fatalf("failed to create signed message: %v", err)
	}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, valid, err := client.Validate(signedStr)
			if err != nil || !valid {
				b.Fatalf("failed to validate signed message: %v, valid=%v", err, valid)
			}
		}
	})
}
