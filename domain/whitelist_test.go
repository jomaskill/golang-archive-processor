package domain

import (
	"github.com/joho/godotenv"
	uuid "github.com/satori/go.uuid"
	"os"
	"testing"
)

var whitelist Whitelist

func TestWhitelist_New(t *testing.T) {
	uuid := uuid.NewV4().String()
	contact := "teste@teste.com"

	data := map[string]string{"user_uuid": uuid}
	whitelist.New(contact, data)

	if whitelist.UserUuid != uuid{
		t.Errorf("Expected %v, got %v", uuid, whitelist.UserUuid)
	}

	if whitelist.Contact != contact{
		t.Errorf("Expected %v, got %v", uuid, whitelist.Contact)
	}
}

func TestWhitelist_Index(t *testing.T) {
	godotenv.Load()

	if whitelist.Index() != os.Getenv("ELASTIC_INDEX_WHITELIST"){
		t.Errorf("Expected %v, got %v", os.Getenv("ELASTIC_INDEX_WHITELIST"), whitelist.Index())
	}
}

func TestWhitelist_Data(t *testing.T) {
	if whitelist.Data() != whitelist{
		t.Errorf("Expected %v, got %v", whitelist, whitelist.Data())
	}
}

