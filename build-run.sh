#!/bin/bash

docker rm --force notification-service
docker rmi --force notification-service

docker build -t notification-service .
docker run -d --env-file ./.env --name notification-service -p 11500:8080 --restart unless-stopped -v ./data:/app/data notification-service
