---
name: sign-crypto
description: 加密签名工具，提供RSA密钥对生成/导入、HMAC签名、XOR加密、TOTP验证码。当需要RSA加解密、HMAC签名验证、XOR混淆加密、或生成/验证TOTP验证码时使用。
---

# sign - 加密签名

提供RSA密钥对管理、HMAC签名、XOR加密与TOTP动态验证码。

## 快速开始

```go
import "github.com/kamalyes/go-toolbox/pkg/sign"
```

RSA密钥对：

```go
keyPair, err := sign.GenerateRsaKeyPair(2048)
pem := sign.ExportRsaPrivateKeyToPEM(keyPair.PrivateKey)
```

HMAC签名：

```go
signer := sign.NewHMACSigner(crypto.SHA256)
sig := signer.Sign([]byte("data"))
```

## 完整API索引

### 函数

#### RSA

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `GenerateRsaKeyPair` | `func(keySize RsaKeySize) (*RsaKeyPair, error)` | 生成RSA密钥对 |
| `ExportRsaPrivateKeyToPEM` | `func(key *rsa.PrivateKey) []byte` | 导出RSA私钥为PEM |
| `ExportRsaPublicKeyToPEM` | `func(key *rsa.PublicKey) []byte` | 导出RSA公钥为PEM |
| `NewRsaCryptoFromKeys` | `func(priv, pub []byte) (RsaCrypto, error)` | 从密钥字节创建RSA加密器 |
| `NewRsaCryptoFromPrivateFile` | `func(path string) (RsaCrypto, error)` | 从私钥文件创建RSA加密器 |
| `NewRsaCryptoFromPublicPEM` | `func(pem []byte) (RsaCrypto, error)` | 从公钥PEM创建RSA加密器 |
| `ParsePrivateKey` | `func(data []byte) (*rsa.PrivateKey, error)` | 解析RSA私钥 |
| `ParsePublicKey` | `func(data []byte) (*rsa.PublicKey, error)` | 解析RSA公钥 |
| `DecryptOAEPWithPrivateKey` | `func(privateKey *rsa.PrivateKey, ciphertext []byte, hashFunc func() hash.Hash) ([]byte, error)` | RSA私钥OAEP解密（简易版，hashFunc传nil默认SHA256） |
| `EncryptOAEPWithPublicKey` | `func(publicKey *rsa.PublicKey, plaintext []byte, hashFunc func() hash.Hash) ([]byte, error)` | RSA公钥OAEP加密（简易版，hashFunc传nil默认SHA256） |
| `RSAPublicKeyToJWK` | `func(publicKey *rsa.PublicKey) (n, e string)` | RSA公钥转JWK格式（Base64URL编码的模数和指数） |

#### HMAC

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `NewGenericHMACSigner` | `func(hash crypto.Hash) *GenericHMACSigner` | 创建通用HMAC签名器 |
| `NewHMACSigner` | `func(hash crypto.Hash, key []byte) *SignerClient[[]byte]` | 创建HMAC签名客户端 |
| `RegisterSigner` | `func(signer Signer)` | 注册签名器 |

#### XOR

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `NewXORCipher` | `func(key []byte) *XORCipher` | 创建XOR加密器 |

#### TOTP

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `DefaultTOTPConfig` | `func() *TOTPConfig` | 获取默认TOTP配置 |
| `GenerateTOTPSecret` | `func(secretLength int) string` | 生成TOTP密钥（Base32编码） |
| `GenerateTOTPURI` | `func(secret, account, issuer string, config *TOTPConfig) string` | 生成TOTP URI（otpauth://） |
| `ValidateTOTPCode` | `func(secret, code string, config *TOTPConfig) bool` | 验证TOTP验证码 |
| `GenerateTOTPCode` | `func(secret string, config *TOTPConfig) (string, error)` | 生成TOTP验证码 |
| `GenerateBackupCodes` | `func(count int) []string` | 生成备份码 |
| `ConsumeBackupCode` | `func(backupCodesJSON, code string) (bool, string)` | 从JSON备份码数组中消耗一个码 |

#### Bcrypt

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `GenerateFromPassword` | `func(data []byte, cost ...int) ([]byte, error)` | bcrypt哈希，cost可选（默认10） |
| `CompareHashAndPassword` | `func(hashed, data []byte) error` | 校验数据与bcrypt哈希是否匹配 |

### 类型

| 导出名称 | 说明 |
|---|---|
| `RsaKeySize` | RSA密钥大小类型 |
| `RsaKeyPair` | RSA密钥对类型 |
| `RsaCrypto` | RSA加密器接口 |
| `Signer` | 签名器接口 |
| `GenericHMACSigner` | 通用HMAC签名器类型 |
| `HashCryptoFunc` | 哈希加密函数类型 |
| `XORCipher` | XOR加密器类型 |
| `SignerClient[T]` | 签名客户端泛型类型 |
| `TOTPConfig` | TOTP配置类型 |

### 常量/变量

| 导出名称 | 值/类型 | 说明 |
|---|---|---|
| `SupportHMACCryptoFunc` | map[crypto.Hash]HashCryptoFunc | 支持的HMAC加密函数映射 |
| `ErrUnsupportedAlgorithmHMAC` | error | 不支持的HMAC算法错误 |

## 注意事项

- RSA密钥长度建议至少2048位
- `NewXORCipher` 仅用于简单混淆，非安全加密
- TOTP验证码有时钟窗口容忍，服务端与客户端时间需同步
