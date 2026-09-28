package pkg

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

type HashConfig struct {
	memory  uint32
	time    uint32
	threads uint8
	keyLen  uint32
	saltLen uint32
}

// func NewHashConfig(memory, time, keyLen, saltLen uint32, threads uint8) *HashConfig {
// 	return &HashConfig{
// 		memory:  memory,
// 		time:    time,
// 		threads: threads,
// 		keyLen:  keyLen,
// 		saltLen: saltLen,
// 	}
// }

func NewRecommendHashConfig() *HashConfig {
	return &HashConfig{
		memory:  64 * 1024,
		time:    2,
		threads: 2,
		keyLen:  32,
		saltLen: 16,
	}
}

func (h *HashConfig) GenerateSalt() []byte {
	salt := make([]byte, h.saltLen)

	rand.Read(salt)
	return salt
}

func (h *HashConfig) GenerateHash(password string) string {
	salt := h.GenerateSalt()
	hash := argon2.IDKey([]byte(password), salt, h.time, h.memory, h.threads, h.keyLen)

	base64Hash := base64.RawStdEncoding.EncodeToString(hash)
	base64Salt := base64.RawStdEncoding.EncodeToString(salt)

	completeHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, h.memory, h.time, h.threads, base64Salt, base64Hash)

	return completeHash
}

func CompareHash(password string, hashedPass string) error {
	var version int
	var memory, time uint32
	var threads uint8

	result := strings.Split(hashedPass, "$")

	if len(result) != 6 {
		return errors.New("incorrect hash format")
	}

	if result[1] != "argon2id" {
		return errors.New("incorrect hash type")
	}

	if _, err := fmt.Sscanf(result[2], "v=%d", &version); err != nil {
		return err
	}

	if version != argon2.Version {
		return errors.New("incorrect hash version")
	}

	if _, err := fmt.Sscanf(result[3], "m=%d,t=%d,p=%d", &memory, &time, &threads); err != nil {
		return err
	}

	salt, err := base64.RawStdEncoding.DecodeString(result[4])
	if err != nil {
		return err
	}

	hash, err := base64.RawStdEncoding.DecodeString(result[5])
	if err != nil {
		return err
	}

	newHash := argon2.IDKey([]byte(password), salt, time, memory, threads, uint32(len(hash)))

	if subtle.ConstantTimeCompare(hash, newHash) == 0 {
		return errors.New("mismatch hash")
	}

	return nil
}
