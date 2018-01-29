package queue

import (
	"log"

	"github.com/brunoleonel/email-consumer-concurrent/email"
	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Printf("%s: %s", msg, err)
		return
	}
}

//SendEmails - Percorre as mensagens na fila e invoca a função de envio de e-mail
func SendEmails() {

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "[queue] Houve uma falha na conexao com o servidor AMQP.")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "[queue] Houve uma falha na abertura do canal com o servidor AMQP.")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"email_conc",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "[queue] Falha na criação da fila.")

	err = ch.Qos(
		1, //prefetch count
		0, //prefetch size
		false,
	)
	failOnError(err, "[queue] Falha ao setar o Qos")

	msgs, err := ch.Consume(
		q.Name,
		"",
		false, //auto-acknowledge
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "[queue] Houve uma falha ao consumir a fila.")

	forever := make(chan bool)

	go func() {
		for d := range msgs {

			emailObj, err := email.ParseMessage(d.Body)
			failOnError(err, "[queue] Houve uma falha ao converter a mensagem.")

			err = email.SendEmail(emailObj)
			failOnError(err, "[queue] Houve uma falha ao enviar o e-mail.")
			d.Ack(false)
		}
	}()

	log.Printf("Grrr... estou com fome... Me dê e-mails ou aperte CTRL+C para sair :D")
	<-forever
}
