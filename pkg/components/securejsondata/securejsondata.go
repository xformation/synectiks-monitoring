package securejsondata

import (
	"github.com/xformation/synectiks-monitoring/pkg/log"
	"github.com/xformation/synectiks-monitoring/pkg/setting"
	"github.com/xformation/synectiks-monitoring/pkg/util"
)

type SecureJsonData map[string][]byte

func (s SecureJsonData) Decrypt() map[string]string {
	decrypted := make(map[string]string)
	for key, data := range s {
		decryptedData, err := util.Decrypt(data, setting.SecretKey)
		if err != nil {
			log.Fatal(4, err.Error())
		}

		decrypted[key] = string(decryptedData)
	}
	return decrypted
}

func GetEncryptedJsonData(sjd map[string]string) SecureJsonData {
	encrypted := make(SecureJsonData)
	for key, data := range sjd {
		encryptedData, err := util.Encrypt([]byte(data), setting.SecretKey)
		if err != nil {
			log.Fatal(4, err.Error())
		}

		encrypted[key] = encryptedData
	}
	return encrypted
}
