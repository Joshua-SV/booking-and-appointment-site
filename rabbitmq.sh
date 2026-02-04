#!/bin/bash

start_or_run () {
    docker inspect doc_appointment_rabbitmq > /dev/null 2>&1

    if [ $? -eq 0 ]; then
        echo "Starting Doc Appointment RabbitMQ container..."
        docker start doc_appointment_rabbitmq
    else
        echo "Doc Appointment RabbitMQ container not found, creating a new one..."
        docker run -d --name doc_appointment_rabbitmq -p 5672:5672 -p 0707:15672 rabbitmq:3.13-management
    fi
}

case "$1" in
    start)
        start_or_run
        ;;
    stop)
        echo "Stopping Doc Appointment RabbitMQ container..."
        docker stop doc_appointment_rabbitmq
        ;;
    logs)
        echo "Fetching logs for Doc Appointment RabbitMQ container..."
        docker logs -f doc_appointment_rabbitmq
        ;;
    *)
        echo "Usage: $0 {start|stop|logs}"
        exit 1
esac