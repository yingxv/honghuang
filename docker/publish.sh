#!/bin/bash
set -e

docker build --file ./Dockerfile.stock --tag ngekaworu/stock-go ..;
docker push ngekaworu/stock-go;

docker build --file ./Dockerfile.flashcard --tag ngekaworu/flashcard-go ..;
docker push ngekaworu/flashcard-go;

docker build --file ./Dockerfile.time-mgt --tag ngekaworu/time-mgt-go ..;
docker push ngekaworu/time-mgt-go;

docker build --file ./Dockerfile.todolist --tag ngekaworu/todo-list-go ..;
docker push ngekaworu/todo-list-go;

docker build --file ./Dockerfile.user-center --tag ngekaworu/ngekaworu/user-center-go ..;
docker push ngekaworu/ngekaworu/user-center-go;
