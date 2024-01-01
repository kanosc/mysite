package database

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type RootKeyMaterial struct {
	Id   int64 `gorm:"primaryKey;autoIncrement:false"`
	Data string
}

type RootKey struct {
	Data []byte
}

type EncKey struct {
	Id    int64 `gorm:"primaryKey;autoIncrement:false"`
	Data  string
	Nonce string
}

type ClientRegInfo struct {
	ClientId string `gorm:"primaryKey"`
	RegTime  time.Time
	KeyNum   int64
}

type ClientKey struct {
	KeyId       string `gorm:"primaryKey"`
	GenericKey  string
	PubKey      string
	PriKey      string
	KeyType     string
	KeyUsage    string
	Nonce       string
	CreatedTime time.Time
	UpdatedTime time.Time
	ClientId    string
	HistoryKeys string
}

var db *gorm.DB

func connectDB() error {
	dbLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,        // Don't include params in the SQL log
			Colorful:                  false,       // Disable color
		})
	dsn := "root:mysql3306@tcp(127.0.0.1:3306)/kms?charset=utf8mb4&parseTime=True&loc=Local"
	mysqlDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: dbLogger,
	})
	if err != nil {
		log.Println("Open db failed", err.Error())
		return err
	}
	genericDb, err := mysqlDB.DB()
	if err != nil {
		log.Println("get generic db interface failed", err.Error())
		return err
	}
	genericDb.SetMaxIdleConns(10)
	genericDb.SetMaxOpenConns(100)
	genericDb.SetConnMaxLifetime(time.Hour)
	db = mysqlDB
	return nil
}

func CreateTableIfNotExist(tableObj interface{}) error {
	if !db.Migrator().HasTable(tableObj) {
		return db.Migrator().CreateTable(tableObj)
	}
	return nil
}

func DropTable(o interface{}) error {
	return db.Migrator().DropTable(o)
}

func InsertRootKeyMaterial(rkm *RootKeyMaterial) error {
	return db.Create(rkm).Error
}

func QueryRootKeyMaterial(rkm *RootKeyMaterial) error {
	result := db.First(rkm, 0)
	if result.RowsAffected == 0 {
		return NoRecordFoundError
	}
	return result.Error
}

func InsertEncKey(ek *EncKey) error {
	return db.Create(ek).Error
}

func QueryEncKey(ek *EncKey) error {
	result := db.First(ek, 0)
	if result.RowsAffected == 0 {
		return NoRecordFoundError
	}
	return result.Error
}

func InsertClientRegInfo(cg *ClientRegInfo) error {
	return db.Create(cg).Error
}

func ClientRegInfoExist(cg *ClientRegInfo) bool {
	return db.First(cg).RowsAffected == 1
}

func InsertClientKey(ck *ClientKey) error {
	return db.Create(ck).Error
}

func QueryClientKeyByKeyId(keyId string) (*ClientKey, error) {
	if keyId == "" {
		return nil, InvalidInputParamError
	}
	clientKey := &ClientKey{
		KeyId: keyId,
	}
	result := db.First(clientKey)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, NoRecordFoundError
	}
	return clientKey, nil
}

func QueryAllClientKeyByClientId(clientId string) ([]ClientKey, error) {
	if clientId == "" {
		return nil, InvalidInputParamError
	}
	var keys []ClientKey
	result := db.Where("client_id = ?", clientId).Find(&keys)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, NoRecordFoundError
	}
	return keys, nil
}

func CreateAllTables() error {
	tables := []interface{}{new(RootKeyMaterial), new(EncKey), new(ClientKey), new(ClientRegInfo)}
	var err error
	for i := range tables {
		err = CreateTableIfNotExist(tables[i])
		if err != nil {
			log.Println("Create table failed", i)
		}
	}
	if err != nil {
		return CreateTableError
	}
	return nil
}

func InitKmsDatabase() error {
	err := connectDB()
	if err != nil {
		log.Println("Connect db failed", err.Error())
		return err
	}
	err = CreateAllTables()
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}
