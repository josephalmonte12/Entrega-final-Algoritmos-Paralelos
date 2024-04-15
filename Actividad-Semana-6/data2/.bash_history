rabbitmqctl stop_app
rabbitmqctl join_cluster rabbit@rabbitmq1
rabbitmqctl start_app
docker restart actividad-semana-6-rabbitmq-1 actividad-semana-6-rabbitmq2-1
exit
rabbitmqctl stop_app
rabbitmqctl join_cluster rabbit@rabbitmq1
rabbitmqctl start_app
docker exec -it actividad-semana-6-rabbitmq2-1 bash
exit
ping rabbitmq1
exit
ping rabbitmq1
apt-get update && apt-get install -y telnet
telnet rabbitmq1 5672
docker-compose down
exit
