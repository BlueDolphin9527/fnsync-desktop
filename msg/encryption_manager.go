package msg

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"io"

	"github.com/cxfksword/fnsync-desktop/entity"
	"github.com/cxfksword/fnsync-desktop/utils"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/pbkdf2"
)

const (
	salt                = "1234qweasdzxc"
	security_iterations = 100
	security_key_len    = 256 / 8
	iv_size_bytes       = 12
)

type EncryptionManager struct {
	code string
	key  []byte
}

func NewEncryptionManager(code []byte) *EncryptionManager {
	return &EncryptionManager{
		code: string(code),
		key:  pbkdf2.Key(code, []byte(salt), security_iterations, security_key_len, sha1.New),
	}
}

func (e *EncryptionManager) Encrypt(data []byte) ([]byte, error) {
	encryptMsg := entity.NewEncryptMsg(data)
	data = utils.ToJSON(encryptMsg)
	log.Debug().Msgf("Encrypt text: %s", string(data))

	block, err := aes.NewCipher(e.key)
	if err != nil {
		return nil, err
	}

	// Never use more than 2^32 random nonces with a given key because of the risk of a repeat.
	nonce := make([]byte, iv_size_bytes)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	ciphertext := aesgcm.Seal(nil, nonce, data, nil)
	return append(nonce, ciphertext...), nil
}

func (e *EncryptionManager) Decrypt(data []byte) ([]byte, error) {
	nonce := data[:iv_size_bytes]
	ciphertext := data[iv_size_bytes:]

	block, err := aes.NewCipher(e.key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	log.Debug().Msgf("Decrypt text result: %s", string(plaintext))
	return plaintext, nil
}
