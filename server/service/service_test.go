package service

import (
	"encoding/hex"
	"os"
	"testing"

	"github.com/kanosc/mysite/server/database"
)

func TestGetTime(t *testing.T) {
	time, err := getSetTimeStampFile()
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(time)
}

func TestInitRootKey(t *testing.T) {
	os.Setenv("service_name", "key_management_service")
	err := database.InitKmsDatabase()
	if err != nil {
		t.Fatal(err.Error())
	}
	InitRootKey()
	t.Logf("generated root key[%s]\n", hex.EncodeToString(rootKeyData))
}

func TestInitEncKey(t *testing.T) {
	os.Setenv("service_name", "key_management_service")
	err := database.InitKmsDatabase()
	if err != nil {
		t.Fatal(err.Error())
	}
	InitRootKey()
	InitEncKey()
	t.Logf("generated root key[%s]\n", hex.EncodeToString(rootKeyData))
	t.Logf("generated enc key[%s]\n", hex.EncodeToString(encKeyData))
}

// func TestCreateKey(t *testing.T) {
// 	os.Setenv("service_name", "key_management_service")
// 	err := database.InitKmsDatabase()
// 	if err != nil {
// 		t.Fatal(err.Error())
// 	}
// 	InitRootKey()
// 	InitEncKey()
// 	t.Logf("generated root key[%s]\n", hex.EncodeToString(rootKeyData))
// 	t.Logf("generated enc key[%s]\n", hex.EncodeToString(encKeyData))

// 	_, err = CreateClientKey("test1", "default", "encrypt")
// 	if err != nil {
// 		t.Errorf(err.Error())
// 	}
// }

func TestQueryKey(t *testing.T) {
	os.Setenv("service_name", "key_management_service")
	err := database.InitKmsDatabase()
	if err != nil {
		t.Fatal(err.Error())
	}
	InitRootKey()
	InitEncKey()
	t.Logf("generated root key[%s]\n", hex.EncodeToString(rootKeyData))
	t.Logf("generated enc key[%s]\n", hex.EncodeToString(encKeyData))

	userKey, err := QueryClientKey("1725697165335990272", "default")
	if err != nil {
		t.Errorf(err.Error())
	}

	t.Logf("plain user key for id[%s] in base64 is[%s]", userKey.KeyId, userKey.GenericKey)
}
