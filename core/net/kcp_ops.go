package net

import (
	"crypto/sha1"
	"errors"
	"time"

	kcp "github.com/xtaci/kcp-go"
	"golang.org/x/crypto/pbkdf2"
)

var (
	ErrInvalidCrypt = errors.New("ErrInvalidCrypt")
)

// KCPOptions 配置文件
type KCPOptions struct {
	Addr string

	// Key 通信加密
	Key string

	// Crypt aes, aes-128, aes-192, salsa20, blowfish, twofish, cast5, 3des, tea, xtea, xor, sm4, none
	Crypt string

	// Mode fast3, fast2, fast, normal, manual
	Mode string

	// MTU set maximum transmission unit for UDP packets
	MTU int

	// SndWnd set send window size(num of packets)
	SndWnd int

	// RcvWnd set receive window size(num of packets)
	RcvWnd int

	// DataShard set reed-solomon erasure coding - datashard
	DataShard int

	// ParityShard set reed-solomon erasure coding - parityshard
	ParityShard int

	// DSCP (6bit)
	DSCP int

	// NoComp disable compression
	NoComp bool

	// AckNodelay flush ack immediately when a packet is received
	AckNodelay bool

	// SockBuf per-socket buffer in bytes
	SockBuf int

	// KeepAlive seconds between heartbeats
	KeepAlive int

	// ReadDeadline 读取超时
	ReadDeadline time.Duration
}

func (kops *KCPOptions) Delay() (noDelay int, interval int, resend int, noCongestion int) {
	switch kops.Mode {
	case "normal":
		noDelay, interval, resend, noCongestion = 0, 40, 2, 1
	case "fast":
		noDelay, interval, resend, noCongestion = 0, 40, 2, 1
	case "fast2":
		noDelay, interval, resend, noCongestion = 0, 40, 2, 1
	case "fast3":
		noDelay, interval, resend, noCongestion = 0, 40, 2, 1
	default:
		noDelay, interval, resend, noCongestion = 0, 40, 2, 1
	}

	return
}

func (kops *KCPOptions) BlockCrypt() (block kcp.BlockCrypt, err error) {
	pass := pbkdf2.Key([]byte(kops.Key), []byte("foxchat-kcp-go"), 4096, 32, sha1.New)
	crypt := kops.Crypt
	switch crypt {
	case "sm4":
		block, err = kcp.NewSM4BlockCrypt(pass[:16])
	case "tea":
		block, err = kcp.NewTEABlockCrypt(pass[:16])
	case "xor":
		block, err = kcp.NewSimpleXORBlockCrypt(pass)
	case "none":
		block, err = kcp.NewNoneBlockCrypt(pass)
	case "aes-128":
		block, err = kcp.NewAESBlockCrypt(pass[:16])
	case "aes-192":
		block, err = kcp.NewAESBlockCrypt(pass[:24])
	case "blowfish":
		block, err = kcp.NewBlowfishBlockCrypt(pass)
	case "twofish":
		block, err = kcp.NewTwofishBlockCrypt(pass)
	case "cast5":
		block, err = kcp.NewCast5BlockCrypt(pass[:16])
	case "3des":
		block, err = kcp.NewTripleDESBlockCrypt(pass[:24])
	case "xtea":
		block, err = kcp.NewXTEABlockCrypt(pass[:16])
	case "salsa20":
		block, err = kcp.NewSalsa20BlockCrypt(pass)
	case "aes":
		block, err = kcp.NewAESBlockCrypt(pass)
	default:
		block, err = nil, ErrInvalidCrypt
	}

	return
}

func NewDefaultKCPOptions() *KCPOptions {
	return &KCPOptions{
		Addr:         ":19099",
		Crypt:        "aes",
		Mode:         "fast",
		MTU:          1350,
		SndWnd:       1024,
		RcvWnd:       1024,
		DataShard:    10,
		ParityShard:  3,
		DSCP:         0,
		NoComp:       false,
		AckNodelay:   false,
		SockBuf:      4194304,
		KeepAlive:    10,
		ReadDeadline: time.Second * 10,
	}
}
