package seed

import (
	"app/framework/rabbitmq"
	uuid "github.com/satori/go.uuid"
)

func RabbitMsgSeed()  {
	rabbitmq.Publish("blacklist", "{\"file\": \"teste.txt\", \"user_uuid\": \""+uuid.NewV4().String()+"\", \"data_type\": \"staff\"}")
	//rabbitmq.Publish("blacklist", "{\"file\": \"teste.txt\", \"user_uuid\": \""+uuid.NewV4().String()+"\", \"data_type\": \"staff\"}")
	//rabbitmq.Publish("blacklist", "{\"file\": \"teste.txt\", \"user_uuid\": \""+uuid.NewV4().String()+"\", \"data_type\": \"quality\"}")

	//rabbitmq.Publish("whitelist", "{\"file\": \"teste.txt\", \"user_uuid\": \""+uuid.NewV4().String()+"\"}")
}
