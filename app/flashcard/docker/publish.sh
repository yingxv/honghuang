#!/bin/bash
set -e

tag=ngekaworu/flashcard-go

docker build --file ./Dockerfile --tag ${tag} ..;
docker push ${tag};
