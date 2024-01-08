package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	// "runtime/debug"
	"strconv"
	"time"

	"github.com/kanosc/mysite/server/database"
	"github.com/kanosc/mysite/server/mycrypto"
	"github.com/kanosc/mysite/server/sequencer"

	"github.com/coocood/freecache"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var clientKeyCache = freecache.NewCache(100 * 1024 * 1024)

/*
	func init() {
		debug.SetGCPercent(20)
		os.Setenv("service_name", "key_management_service")
		err := database.InitKmsDatabase()
		if err != nil {
			panic(err.Error())
		}
		err = database.CreateAllTables()
		if err != nil {
			panic(err.Error())
		}

		InitRootKey()
		InitEncKey()
	}
*/
type UserKey struct {
	KeyType    string
	KeyId      string
	GenericKey string
	PriKey     string
	PubKey     string
}

var rootKeyData []byte
var encKeyData []byte

func getServiceName() string {
	return os.Getenv("service_name")
}

func getSetTimeStampFile() (string, error) {
	fileName := "system_init_time.txt"
	_, err := os.Stat(fileName)
	timeData := strconv.FormatInt(time.Now().UnixNano(), 10)
	if os.IsNotExist(err) {
		err = os.WriteFile(fileName, []byte(timeData), 0644)
		if err != nil {
			return timeData, err
		}
	}
	lastTime, err := os.ReadFile(fileName)
	return string(lastTime), err
}

func getSetNonceDB() ([]byte, error) {
	rkm := new(database.RootKeyMaterial)
	err := database.QueryRootKeyMaterial(rkm)
	if err != nil {
		if errors.Is(err, database.NoRecordFoundError) {
			nonce, err := mycrypto.GenerateRand(32)
			if err != nil {
				return nil, err
			}
			newrkm := new(database.RootKeyMaterial)
			newrkm.Data = mycrypto.EncodeBytesToBase64(nonce)
			newrkm.Id = 0
			err = database.InsertRootKeyMaterial(newrkm)
			if err != nil {
				return nil, err
			}
			return nonce, nil
		}
	}
	return mycrypto.DecodeBase64ToBytes(rkm.Data)

}

func InitRootKey() {
	sn := getServiceName()
	if sn == "" {
		panic("root key material doesn't exist in env")
	}
	ts, err := getSetTimeStampFile()
	if err != nil {
		panic("root key material dosesn't exist in config file")
	}
	nc, err := getSetNonceDB()
	if err != nil {
		panic("recover root key material from db failed")
	}
	rkd, err := mycrypto.ScriptGenerateKey([]byte(sn+ts), nc)
	if err != nil {
		panic("recover root key failed")
	}
	rootKeyData = rkd
}

func getSetEncKeyDB(rk []byte) ([]byte, error) {
	if rk == nil {
		return nil, errors.New("root key not initial")
	}
	encKey := new(database.EncKey)
	err := database.QueryEncKey(encKey)
	if err != nil {
		if errors.Is(err, database.NoRecordFoundError) {
			nonce, err := mycrypto.GenerateRand(32)
			if err != nil {
				return nil, err
			}
			sn := getServiceName()
			if sn == "" {
				panic("root key material doesn't exist in env")
			}
			ekdata, err := mycrypto.ScriptGenerateKey([]byte(sn), nonce)
			if err != nil {
				return nil, err
			}

			aesNonce, err := mycrypto.GenerateRand(12)
			if err != nil {
				return nil, err
			}
			encryptedEk, err := mycrypto.AES256Encrypt(ekdata, aesNonce, rootKeyData)
			if err != nil {
				return nil, err
			}

			ekdb := database.EncKey{
				Id:    0,
				Data:  mycrypto.EncodeBytesToBase64(encryptedEk),
				Nonce: mycrypto.EncodeBytesToBase64(aesNonce),
			}
			err = database.InsertEncKey(&ekdb)
			if err != nil {
				panic("insert encrypt key to DB failed")
			}
			return ekdata, nil
		}
		return nil, err
	}
	secretEkBytes, err := mycrypto.DecodeBase64ToBytes(encKey.Data)
	if err != nil {
		return nil, err
	}
	aesNonceBytes, err := mycrypto.DecodeBase64ToBytes(encKey.Nonce)
	if err != nil {
		return nil, err
	}
	return mycrypto.AES256Decrypt(secretEkBytes, aesNonceBytes, rk)
}

func InitEncKey() {
	ek, err := getSetEncKeyDB(rootKeyData)
	if err != nil {
		panic(err.Error())
	}
	encKeyData = ek
}

func CreateClientKey(clientId, keyType, keyUsage string) (*UserKey, error) {
	switch keyType {
	case "default", "":
		return createGenericClientKey(clientId, keyUsage)
	default:
		return nil, errors.New("Unknown key type")
	}

}

func createGenericClientKey(clientId, keyUsage string) (*UserKey, error) {
	nonce, err := mycrypto.GenerateRand(32)
	if err != nil {
		return nil, err
	}
	ukdata, err := mycrypto.ScriptGenerateKey([]byte(clientId), nonce)
	if err != nil {
		return nil, err
	}

	aesNonce, err := mycrypto.GenerateRand(12)
	if err != nil {
		return nil, err
	}
	encryptedUk, err := mycrypto.AES256Encrypt(ukdata, aesNonce, encKeyData)
	if err != nil {
		return nil, err
	}
	keyId := sequencer.DefaultIdGenerator().GenerateId()
	newKeyDB := &database.ClientKey{
		KeyType:     "default",
		KeyId:       keyId,
		ClientId:    clientId,
		KeyUsage:    keyUsage,
		GenericKey:  mycrypto.EncodeBytesToBase64(encryptedUk),
		Nonce:       mycrypto.EncodeBytesToBase64(aesNonce),
		CreatedTime: time.Now(),
		UpdatedTime: time.Now(),
	}
	err = database.InsertClientKey(newKeyDB)
	if err != nil {
		return nil, errors.New("insert new key into db failed")
	}
	return &UserKey{
		KeyType: "default",
		KeyId:   keyId,
	}, nil
}

func QueryClientKey(keyId, keyType string) (*UserKey, error) {
	switch keyType {
	case "default", "":
		return queryGenericClientKey(keyId)
	default:
		return nil, errors.New("Unknown key type")
	}
}

func queryGenericClientKeyFromCache(keyId string) (*UserKey, error) {
	k, err := clientKeyCache.Get([]byte(keyId))
	if err != nil {
		return nil, err
	}
	uk := new(UserKey)
	err = json.Unmarshal(k, uk)
	if err != nil {
		return nil, err
	}
	fmt.Println("****************************")
	fmt.Println("cache hit", keyId)
	return uk, nil
}

func queryGenericClientKey(keyId string) (*UserKey, error) {
	// query from local cache
	keyInCache, err := queryGenericClientKeyFromCache(keyId)
	if err == nil {
		return keyInCache, nil
	}

	// query from database
	fmt.Println("****************************")
	fmt.Println("cache doesn't hit", keyId)
	clientKeyDB, err := database.QueryClientKeyByKeyId(keyId)
	if err != nil {
		return nil, err
	}
	secretUkBytes, err := mycrypto.DecodeBase64ToBytes(clientKeyDB.GenericKey)
	if err != nil {
		return nil, err
	}
	aesNonceBytes, err := mycrypto.DecodeBase64ToBytes(clientKeyDB.Nonce)
	if err != nil {
		return nil, err
	}
	plainUkBytes, err := mycrypto.AES256Decrypt(secretUkBytes, aesNonceBytes, encKeyData)
	if err != nil {
		return nil, err
	}
	keyResult := &UserKey{
		KeyType:    "default",
		KeyId:      keyId,
		GenericKey: mycrypto.EncodeBytesToBase64(plainUkBytes),
	}

	// update local cache
	keyResultJson, err := json.Marshal(keyResult)
	if err != nil {
		return nil, err
	}
	err = clientKeyCache.Set([]byte(keyId), keyResultJson, 600)
	if err != nil {
		return nil, err
	}
	return keyResult, nil
}

func KeyOperationAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("clientid", "test.1")

		c.Next()
	}
}

func KeyCreateRateLimiter(num int) gin.HandlerFunc {
	limiter := rate.NewLimiter(rate.Limit(num), 20)
	return func(c *gin.Context) {
		// if !limiter.Allow() {
		// 	c.JSON(http.StatusForbidden, gin.H{
		// 		"err": "Server busy",
		// 	})
		// 	c.Abort()
		// 	return
		// }

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		err := limiter.Wait(ctx)
		if err != nil {
			fmt.Println("wait for token timeout")
			c.JSON(http.StatusForbidden, gin.H{
				"err": "Server busy",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func KeyRequestChecker() gin.HandlerFunc {
	return func(c *gin.Context) {
		keyid := c.Param("keyid")
		if keyid == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"err": "Invalid key request",
			})
			c.Abort()
			return
		}
		c.Set("keyid", keyid)

		c.Next()
	}
}

func HandleGetClientKeyById(c *gin.Context) {
	keyid := c.MustGet("keyid").(string)
	keyInfo, err := QueryClientKey(keyid, "default")
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"err": err.Error(),
		})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, *keyInfo)
}

func HandleCreateClientKey(c *gin.Context) {
	clientId := c.MustGet("clientid").(string)
	keyInfo, err := CreateClientKey(clientId, "default", "encrypt")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": err.Error(),
		})
	}
	c.JSON(http.StatusOK, *keyInfo)
}
