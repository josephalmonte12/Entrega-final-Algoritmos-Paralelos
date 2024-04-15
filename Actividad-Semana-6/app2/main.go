package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/streadway/amqp"
)

func main() {
	http.HandleFunc("/", servirInicio)
	http.HandleFunc("/consume", manejadorConsumirMensajes)
	http.HandleFunc("/delete", manejadorEliminarMensajes)
	log.Fatal(http.ListenAndServe(":8082", nil))
}

func servirInicio(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func manejadorEliminarMensajes(w http.ResponseWriter, r *http.Request) {
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

	purged, err := ch.QueuePurge(
		"testQueue", // nombre de la cola
		false,       // no-wait
	)
	if err != nil {
		log.Printf("Error al purgar la cola: %s", err)
		http.Error(w, "Error al purgar la cola", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Mensajes eliminados exitosamente: %d", purged)
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

	// Obtener un solo mensaje de la cola sin confirmación automática
	msg, ok, err := ch.Get(
		"testQueue", // nombre de la cola
		false,       // auto-ack
	)

	if err != nil {
		log.Printf("Error al obtener el mensaje: %s", err)
		http.Error(w, "Error al obtener el mensaje", http.StatusInternalServerError)
		return
	}

	// Comprobar si se recibió un mensaje
	if ok {
		log.Printf("Mensaje consumido: %s", msg.Body)     // Mostrar mensaje en consola
		fmt.Fprintf(w, "Mensaje consumido: %s", msg.Body) // Mostrar mensaje en la app
	} else {
		fmt.Fprint(w, "No hay mensaje para consumir de la Aplicación 1")
	}
}
