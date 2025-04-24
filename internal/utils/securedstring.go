//nolint:exhaustruct,ireturn
package utils

import (
	"crypto/rand"
	"encoding/binary"

	obs "github.com/Dentrax/obscure-go/observer"
)

const KEY int = 54343

type ISecureString interface {
	Apply() ISecureString
	AddWatcher(obs obs.Observer)
	SetKey(key int)
	Set(str string) ISecureString
	Get() string
	GetSelf() *SecureString
	Decrypt() []rune
	RandomizeKey()
	IsEquals(secureString ISecureString) bool
}

type SecureString struct {
	obs.Observable
	Key           int
	RealValue     []rune
	FakeValue     string
	Initialized   bool
	HackDetecting bool
}

func NewString(value string) ISecureString {
	sec := &SecureString{
		Key:           KEY,
		RealValue:     []rune(value),
		FakeValue:     value,
		Initialized:   false,
		HackDetecting: false,
	}

	sec.Apply()

	return sec
}

func (i *SecureString) Apply() ISecureString {
	if !i.Initialized {
		i.RealValue = i.XOR(i.RealValue, i.Key)
		i.Initialized = true
	}

	return i
}

func (i *SecureString) AddWatcher(obs obs.Observer) {
	i.AddObserver(obs)
	i.HackDetecting = true
}

func (i *SecureString) SetKey(key int) {
	i.Key = key
}

func (i *SecureString) RandomizeKey() {
	i.RealValue = i.Decrypt()

	// Generate a secure random 4-byte key
	var keyBytes [4]byte
	_, _ = rand.Read(keyBytes[:])

	// Convert bytes to an int
	i.Key = int(binary.BigEndian.Uint32(keyBytes[:]))

	i.RealValue = i.XOR(i.RealValue, i.Key)
}

func (i *SecureString) XOR(value []rune, key int) []rune {
	res := make([]rune, len(value))

	for i, v := range value {
		//nolint:gosec
		res[i] = v ^ int32(key)
	}

	return res
}

func (i *SecureString) Get() string {
	return string(i.Decrypt())
}

func (i *SecureString) GetSelf() *SecureString {
	return i
}

func (i *SecureString) Set(value string) ISecureString {
	i.RealValue = i.XOR([]rune(value), i.Key)

	if i.HackDetecting {
		i.FakeValue = value
	}

	return i
}

func (i *SecureString) Decrypt() []rune {
	if !i.Initialized {
		i.Key = KEY
		i.FakeValue = ""
		i.RealValue = i.XOR(nil, 0)
		i.Initialized = false

		return nil
	}

	res := i.XOR(i.RealValue, i.Key)

	if i.HackDetecting && string(res) != i.FakeValue {
		i.NotifyAll("hack")
	}

	return res
}

func (i *SecureString) IsEquals(o ISecureString) bool {
	if i.Key != o.GetSelf().Key {
		return string(i.XOR(i.RealValue, i.Key)) == string(i.XOR(o.GetSelf().RealValue, o.GetSelf().Key))
	}

	return string(i.RealValue) == string(o.GetSelf().RealValue)
}
