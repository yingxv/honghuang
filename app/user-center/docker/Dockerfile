FROM alpine:latest

EXPOSE 80

RUN mkdir -p /home/go/app

WORKDIR /home/go/app

COPY ./bin /home/go/app

CMD ["sh", "-c", "./app -l=$ADDR -i=$DBINIT -m=$MONGO -db=$DB -k=$K -r=$REDIS -d=$D"]