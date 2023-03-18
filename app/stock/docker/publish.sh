#!/bin/bash
set -e

tag=ngekaworu/stock-go

docker build --file ./Dockerfile --tag ${tag} ..;
docker push ${tag};
