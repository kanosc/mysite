package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnectDB(t *testing.T) {
	ret := connectDB()
	assert.Equal(t, nil, ret)
}

func TestCreateAndDropTable(t *testing.T) {
	ret := CreateTableIfNotExist(&RootKeyMaterial{})
	assert.Equal(t, nil, ret)
	ret = CreateTableIfNotExist(&EncKey{})
	assert.Equal(t, nil, ret)
	ret = CreateTableIfNotExist(&ClientRegInfo{})
	assert.Equal(t, nil, ret)
	ret = CreateTableIfNotExist(&ClientKey{})
	assert.Equal(t, nil, ret)

	ret = DropTable(&RootKeyMaterial{})
	assert.Equal(t, nil, ret)
	ret = DropTable(&EncKey{})
	assert.Equal(t, nil, ret)
	ret = DropTable(&ClientRegInfo{})
	assert.Equal(t, nil, ret)
	ret = DropTable(&ClientKey{})
	assert.Equal(t, nil, ret)
}
