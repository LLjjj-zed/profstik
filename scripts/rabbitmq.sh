docker run -d --hostname rabbitmq --name rabbitmq -p 15672:15672 -p 5672:5672 -e RABBITMQ_DEFAULT_USER=tiktokRMQ -e RABBITMQ_DEFAULT_PASS=tiktokRMQ -e RABBITMQ_DEFAULT_VHOST=tiktokRMQ  rabbitmq:3.11.8-management