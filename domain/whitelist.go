package domain

import (
	"github.com/asaskevich/govalidator"
	"log"
	"os"
)

type Whitelist struct {
	Contact string `json:"contact" valid:"notnull"`
	UserUuid string `json:"user_uuid" valid:"notnull,uuid"`
}

func init()  {
	govalidator.SetFieldsRequiredByDefault(true)
}

func (w Whitelist) Data() Base{
	return w
}

func (w Whitelist) Index() string {
	return os.Getenv("ELASTIC_INDEX_WHITELIST")
}

func (w *Whitelist) New(value string, msg map[string]string) Whitelist {
	w.Contact = value
	w.UserUuid = msg["user_uuid"]

	w.validate()

	return *w
}

func (w *Whitelist) validate()  {
	_, err := govalidator.ValidateStruct(w)
	if err != nil{
		log.Fatal(err)
	}
}