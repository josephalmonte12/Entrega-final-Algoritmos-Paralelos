package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/streadway/amqp"
)

func main() {
	http.HandleFunc("/", servirInicio)
	http.HandleFunc("/send", manejadorEnviarMensaje)
	http.HandleFunc("/consume", manejadorConsumirMensajes)
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func servirInicio(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
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

	q, err := ch.QueueDeclare(
		"testQueue",
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

	cuerpo := "¡¡¡Probando que funciona correctamente la Aplicación 1!!!"
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

	fmt.Fprintf(w, "Mensaje enviado: %s", cuerpo)
}

func manejadorConsumirMensajes(w http.ResponseWriter, r *http.Request) {
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

	// Obtener un solo mensaje de la cola de la App 3
	msg, ok, err := ch.Get(
		"app1Queue", // nombre de la cola usada por App 3
		false,       // no auto-ack
	)

	if err != nil {
		//log.Printf("Error al obtener el mensaje: %s", err)
		//http.Error(w, "Error al obtener el mensaje", http.StatusInternalServerError)
		fmt.Fprint(w, "No hay mensaje para consumir de la Aplicación 3")
		return
	}

	// Comprobar si hay mensajes
	if ok {
		log.Printf("Mensaje consumido: %s", msg.Body)
		fmt.Fprintf(w, "Mensaje consumido: %s", msg.Body)

		// Eliminar la cola después de consumir el mensaje
		_, err = ch.QueueDelete(
			"app1Queue", // nombre de la cola
			false,       // ifUnused
			false,       // ifEmpty
			false,       // noWait
		)
		if err != nil {
			log.Printf("Error al eliminar la cola: %s", err)
			http.Error(w, "Error al eliminar la cola", http.StatusInternalServerError)
			return
		}
		//fmt.Fprintf(w, "\nCola eliminada con éxito.")
	} else {
		fmt.Fprint(w, "No hay mensajes para consumir de la Aplicación 3")
	}
}
