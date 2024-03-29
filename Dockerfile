FROM golang:1.17

WORKDIR /home/usr/src/mathcord

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin ./...

RUN chmod +x ./boot.sh

CMD ["./boot.sh"]