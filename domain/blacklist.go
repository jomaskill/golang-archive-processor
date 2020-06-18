package domain

import (
	"github.com/asaskevich/govalidator"
	"log"
	"os"
)

type Blacklist struct {
	Contact string `json:"contact" valid:"notnull"`
	DataType string `json:"data_type" valid:"in(staff|quality)"`
	UserUuid string `json:"user_uuid" valid:"notnull,uuid"`
}

func init()  {
	govalidator.SetFieldsRequiredByDefault(true)
}

func (b Blacklist) Data() Base{
	return b
}

func (b Blacklist) Index() string {
	return os.Getenv("ELASTIC_INDEX_BLACKLIST")
}

func (b *Blacklist) New(value string, msg map[string]string) Blacklist {
	b.Contact = value
	b.UserUuid = msg["user_uuid"]
	b.DataType = msg["data_type"]

	b.validate()

	return *b
}

func (b *Blacklist) validate()  {
	_, err := govalidator.ValidateStruct(b)
	if err != nil{
		log.Fatal(err)
	}
}
