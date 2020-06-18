package domain_test

import (
	"app/domain"
	"github.com/joho/godotenv"
	uuid "github.com/satori/go.uuid"
	"os"
	"testing"
)

var blacklist domain.Blacklist

func TestBlacklist_New(t *testing.T) {

	uuid := uuid.NewV4().String()
	dataType := "staff"
	contact := "teste@teste.com"

	data := map[string]string{"user_uuid": uuid, "data_type": dataType}
	blacklist.New(contact, data)

	if blacklist.UserUuid != uuid{
		t.Errorf("Expected %v, got %v", uuid, blacklist.UserUuid)
	}
	if blacklist.DataType != dataType{
		t.Errorf("Expected %v, got %v", uuid, blacklist.DataType)
	}
	if blacklist.Contact != contact{
		t.Errorf("Expected %v, got %v", uuid, blacklist.Contact)
	}

}

func TestBlacklist_Index(t *testing.T) {

	godotenv.Load()

	if blacklist.Index() != os.Getenv("ELASTIC_INDEX_BLACKLIST"){
		t.Errorf("Expected %v, got %v", os.Getenv("ELASTIC_INDEX_BLACKLIST"), blacklist.Index())
	}
}

func TestBlacklist_Data(t *testing.T) {

	if blacklist.Data() != blacklist{
		t.Errorf("Expected %v, got %v", blacklist, blacklist.Data())
	}
}