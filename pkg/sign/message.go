/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-06-05 13:35:59
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-06-05 17:22:17
 * @FilePath: \go-toolbox\pkg\sign\message.go
 * @Description:
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package sign

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/kamalyes/go-toolbox/pkg/random"
	"github.com/kamalyes/go-toolbox/pkg/syncx"
)

// SignerClient 是带状态的签名客户端，支持复用 secretKey、算法和签名器，且并发安全。
type SignerClient[T any] struct {
	secretKey  []byte
	alg        HashCryptoFunc
	serializer Serializer
	signer     Signer
	mu         sync.RWMutex
}

// ----------- 签名器注册与获取 -----------

// signerRegistry 是一个全局映射表，
// 用于存储算法名称到对应签名器实例的映射，方便根据算法名称快速获取签名器
var signerRegistry = make(map[HashCryptoFunc]Signer)

// RegisterSigner 注册一个签名器，将其算法名称作为键，签名器实例作为值
// 方便后续通过算法名称获取对应的签名器
func RegisterSigner(s Signer) {
	signerRegistry[s.Algorithm()] = s
}

// GetSigner 根据算法名称获取对应的签名器实例，如果不存在则返回错误
func GetSigner(alg HashCryptoFunc) (Signer, error) {
	if s, ok := signerRegistry[alg]; ok {
		return s, nil
	}
	return nil, errors.New("unsupported algorithm: " + string(alg))
}

// ----------- 序列化接口定义 -----------

// Serializer 定义了序列化和反序列化的接口，方便支持多种数据格式
// Marshal 将对象序列化为字节切片
// Unmarshal 将字节切片反序列化为对象
type Serializer interface {
	Marshal(v any) ([]byte, error)
	Unmarshal(data []byte, v any) error
}

// JSONSerializer 是 Serializer 接口的 JSON 实现，使用标准库 encoding/json
type JSONSerializer struct{}

// Marshal 使用 json.Marshal 将对象序列化为 JSON 格式字节切片
func (JSONSerializer) Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

// Unmarshal 使用 json.Unmarshal 将 JSON 格式字节切片反序列化为对象
func (JSONSerializer) Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

// ----------- 合并后的签名消息结构 -----------

// SignedMessage 是泛型结构，包含签名算法名称、随机字符串、时间戳和额外数据
// 其中：
// - Alg 字段表示本消息使用的签名算法名称，便于验证时动态选择签名器
// - Send 是随机字符串，用于防止重放攻击，保证每条消息唯一
// - GenUnixMicro 是消息生成时的时间戳，单位为微秒，方便做时效校验
// - ExtraData 是用户自定义的泛型数据，支持任意结构
type SignedMessage[T any] struct {
	Alg          HashCryptoFunc // 签名算法名称（如 "HMAC-SHA256"）
	Send         string         // 随机字符串，用于防重放
	GenUnixMicro int64          // 生成时间戳，微秒级
	ExtraData    T              // 用户自定义额外数据
}

// NewSignerClient 创建默认客户端实例
func NewSignerClient[T any]() *SignerClient[T] {
	return &SignerClient[T]{
		serializer: JSONSerializer{},
	}
}

// WithSecretKey 设置密钥，链式调用
func (c *SignerClient[T]) WithSecretKey(key []byte) *SignerClient[T] {
	syncx.WithLock(&c.mu, func() {
		c.secretKey = key
	})
	return c
}

// WithAlgorithm 设置算法及对应签名器，链式调用
func (c *SignerClient[T]) WithAlgorithm(alg HashCryptoFunc) (*SignerClient[T], error) {
	signer, err := GetSigner(alg)
	if err != nil {
		return c, err
	}
	syncx.WithLock(&c.mu, func() {
		c.alg = alg
		c.signer = signer
	})
	return c, nil
}

// WithSerializer 设置序列化器，链式调用
func (c *SignerClient[T]) WithSerializer(s Serializer) *SignerClient[T] {
	syncx.WithLock(&c.mu, func() {
		c.serializer = s
	})
	return c
}

// Create 创建签名消息字符串，内部使用客户端状态
func (c *SignerClient[T]) Create(extraData T) (string, error) {
	// 读锁获取状态
	state, err := syncx.WithRLockReturn(&c.mu, func() (*SignerClient[T], error) {
		if c.signer == nil {
			return nil, errors.New("signer not set, call WithAlgorithm first")
		}
		if len(c.secretKey) == 0 {
			return nil, errors.New("secretKey not set, call WithSecretKey first")
		}
		return c, nil
	})
	if err != nil {
		return "", err
	}
	// 构造负载结构体，包含算法名、随机字符串、生成时间戳和用户数据
	payload := SignedMessage[T]{
		Alg:          state.alg,
		Send:         random.RandString(16, random.LOWERCASE|random.NUMBER|random.CAPITAL),
		GenUnixMicro: time.Now().UnixMicro(),
		ExtraData:    extraData,
	}
	// 使用传入的序列化器将负载结构体序列化为字节切片
	payloadBytes, err := state.serializer.Marshal(payload)
	if err != nil {
		return "", err
	}

	// 对序列化后的负载字节进行 base64 URL 编码，不带填充字符
	payloadB64 := base64.RawURLEncoding.EncodeToString(payloadBytes)
	// 使用签名器和密钥对签名输入进行签名，得到签名字节
	signature, err := state.signer.Sign([]byte(payloadB64), state.secretKey)
	if err != nil {
		return "", err
	}

	// 对签名字节进行 base64 URL 编码，不带填充
	signatureB64 := base64.RawURLEncoding.EncodeToString(signature)
	// 拼接返回字符串，格式为：base64(payload) + "." + base64(signature)
	return payloadB64 + "." + signatureB64, nil
}

// Validate 验证签名消息字符串，返回负载、是否有效和错误
func (c *SignerClient[T]) Validate(signedStr string) (*SignedMessage[T], bool, error) {
	// 将签名字符串以 '.' 分割，应该分割成两部分：payload 和 signature
	parts := strings.Split(signedStr, ".")
	if len(parts) != 2 {
		return nil, false, errors.New("invalid signed string format: must contain exactly one '.' separator")
	}

	// 定义辅助函数，用于 base64 URL 解码，不带填充
	decode := func(s string) ([]byte, error) {
		return base64.RawURLEncoding.DecodeString(s)
	}

	// 解码负载部分的 base64 字符串
	payloadBytes, err := decode(parts[0])
	if err != nil {
		return nil, false, errors.New("payload base64 decode failed: " + err.Error())
	}
	// 解码签名部分的 base64 字符串
	signature, err := decode(parts[1])
	if err != nil {
		return nil, false, errors.New("signature base64 decode failed: " + err.Error())
	}

	// 反序列化负载字节，得到 SignedMessage 结构体实例
	var payload SignedMessage[T]
	if err := c.serializer.Unmarshal(payloadBytes, &payload); err != nil {
		return nil, false, errors.New("payload unmarshal failed: " + err.Error())
	}

	// 先读锁获取当前 signer 和 alg
	state, _ := syncx.WithRLockReturn(&c.mu, func() (*SignerClient[T], error) {
		return c, nil
	})

	signer := state.signer
	alg := state.alg
	secretKey := state.secretKey

	// 如果签名器为空或算法不匹配，则尝试重新获取签名器并写锁更新
	if signer == nil || alg != payload.Alg {
		newSigner, err := GetSigner(payload.Alg)
		if err != nil {
			return nil, false, errors.New("unsupported algorithm in payload: " + err.Error())
		}
		syncx.WithLock(&c.mu, func() {
			c.signer = newSigner
			c.alg = payload.Alg
			signer = newSigner
		})
	}

	if len(secretKey) == 0 {
		return nil, false, errors.New("secretKey not set, call WithSecretKey first")
	}

	ok, err := signer.Verify([]byte(parts[0]), secretKey, signature)
	if err != nil {
		return nil, false, errors.New("signature verify error: " + err.Error())
	}
	if !ok {
		return nil, false, errors.New("signature verification failed")
	}
	// 签名验证通过，返回负载数据和成功标志
	return &payload, true, nil
}

// ----------- 初始化注册所有支持的 HMAC 签名算法 -----------

// init 函数会在程序启动时自动执行，
// 遍历所有支持的 HMAC 算法，创建对应签名器并注册到全局注册表
func init() {
	for alg := range SupportHMACCryptoFunc {
		if signer, err := NewHMACSigner(alg); err == nil {
			RegisterSigner(signer)
		}
	}
}
