FROM golang:1.23.4

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN make build

EXPOSE 5555
CMD ["./bin/app"]  
