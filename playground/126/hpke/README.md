# crypto/hpke

Go 1.26 新增标准库包 **`crypto/hpke`**, 实现 RFC 9180 **Hybrid Public Key Encryption (HPKE)**.

## HPKE 是什么

**HPKE** 是一套标准化的「**用对方的公钥加密数据**」的构造方法, 把容易写错、互操作性差的 ad-hoc 方案 (各自拼 ECDH + HKDF + AES-GCM) 收成**固定流程 + 可协商算法**.

典型流程可以粗分为三步 (RFC 9180 术语):

| 阶段 | 作用 | 本 demo 选用 |
|------|------|----------------|
| **KEM** (Key Encapsulation) | 用接收方公钥协商出共享秘密, 并产生一段 encapsulation 传给解密方 | P-256 **DHKEM** |
| **KDF** | 从共享秘密派生 AEAD 密钥等 | **HKDF-SHA256** |
| **AEAD** | 认证加密明文 | **AES-128-GCM** |

发送方 **`Seal`**: 输出「封装数据 + 密文」. 接收方 **`Open`**: 用自己的私钥解封装, 再解密. 可选的 **`info`** 等上下文绑定用途/协议, 防止密文被挪用到别的场景.

Go 1.26 的 `crypto/hpke` 还提供 **ML-KEM 及与传统 ECDH 的混合 KEM**, 用于逐步引入**后量子**密钥交换; 本 demo 仅用经典 P-256, 便于理解 API.

## 应用场景

- **TLS 1.3 之外的应用层加密**: 消息/文件只给持有公钥的一方读 (类似 ECIES, 但有标准套件 ID).
- **协议与中间件**: OAuth, VPN, 消息队列 payload 加密等, 库之间按 RFC 9180 对齐算法名与 wire 格式.
- **后量子迁移**: 同一 API 下从 `DHKEM(P-256)` 换到 `MLKEM768X25519` 等混合套件, 而不重写整套拼密逻辑.

和「自己 `ecdh` + `hkdf` + `cipher`」相比, HPKE 的价值主要是**互操作、可审计的算法列表、以及 RFC 级别的测试向量**.

## 本 demo

使用 **P-256 DHKEM** + **HKDF-SHA256** + **AES-128-GCM**, 调用:

- **`hpke.Seal`**: 加密明文.
- **`hpke.Open`**: 解密, 校验 `info` 等上下文一致.

不涉及网络传输, 仅验证 seal/open 往返.

## 运行

```bash
GOWORK=off go test ./hpke/...
```

## 参考

- [Go 1.26 - crypto/hpke](https://go.dev/doc/go1.26#crypto_hpke)
- [RFC 9180](https://www.rfc-editor.org/rfc/rfc9180.html)
- [`crypto/hpke` 包文档](https://pkg.go.dev/crypto/hpke)
