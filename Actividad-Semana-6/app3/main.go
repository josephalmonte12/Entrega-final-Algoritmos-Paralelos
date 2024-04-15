package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/streadway/amqp"
)

func main() {
	http.HandleFunc("/", servirInicio)
	http.HandleFunc("/guardar", manejadorGuardarMensaje)
	http.HandleFunc("/enviar", manejadorEnviarMensaje)
	log.Fatal(http.ListenAndServe(":8083", nil))
}

func servirInicio(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func manejadorGuardarMensaje(w http.ResponseWriter, r *http.Request) {
	conn, err := amqp.Dial("amqp://joseph:1192948@localhost:5673/")
	if err != nil {
		log.Printf("Error al conectar con RabbitMQ: %s", err)
		http.Error(w, "Error al conectar con RabbitMQ", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Printf("Error al abrir un canal: %s", err)
		http.Error(w, "Error al abrir un canal", http.StatusInternalServerError)
		return
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"testQueue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Printf("Error al declarar una cola: %s", err)
		http.Error(w, "Error al declarar una cola", http.StatusInternalServerError)
		return
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Printf("Error al registrar un consumidor: %s", err)
		http.Error(w, "Error al registrar un consumidor", http.StatusInternalServerError)
		return
	}

	db, err := sql.Open("sqlserver", "server=JosephAlmonte\\SQLEXPRESS;database=App3Db;trusted_connection=true")
	if err != nil {
		log.Printf("Error al conectar a la base de datos: %s", err)
		http.Error(w, "Error al conectar a la base de datos", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	select {
	case d := <-msgs:
		_, err = db.Exec("INSERT INTO Messages (Message) VALUES (@p1)", string(d.Body))
		if err != nil {
			log.Printf("Error al insertar el mensaje en la base de datos: %s", err)
			http.Error(w, "Error al insertar el mensaje en la base de datos", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "Mensaje guardado en la base de datos: %s", d.Body)
	default:
		fmt.Fprint(w, "No hay mensajes para guardar")
	}
}

func manejadorEnviarMensaje(w http.ResponseWriter, r *http.Request) {
	conn, err := amqp.Dial("amqp://joseph:1192948@localhost:5673/")
	if err != nil {
		http.Error(w, "Error al conectar con RabbitMQ", http.StatusInternalServerError)
		log.Printf("Error al conectar con RabbitMQ: %s", err)
		return
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		http.Error(w, "Error al abrir un canal", http.StatusInternalServerError)
		log.Printf("Error al abrir un canal: %s", err)
		return
	}
	defer ch.Close()

	// Cola llamada "app1Queue"
	q, err := ch.QueueDeclare(
		"app1Queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		http.Error(w, "Error al declarar una cola", http.StatusInternalServerError)
		log.Printf("Error al declarar una cola: %s", err)
		return
	}

	cuerpo := "¡¡¡El mensaje de guardó correctamente en la Base de Datos!!!"
	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(cuerpo),
		})
	if err != nil {
		http.Error(w, "Error al publicar un mensaje", http.StatusInternalServerError)
		log.Printf("Error al publicar un mensaje: %s", err)
		return
	}

	fmt.Fprintf(w, "Mensaje enviado a App 1: %s", cuerpo)
}
