FROM golang:alpine as first_stage

WORKDIR /client

COPY go.mod /client
COPY go.sum /client

RUN go mod tidy

COPY . /client

EXPOSE 3030

RUN go build -o client_app cmd/client/main.go

FROM alpine

COPY --from=first_stage /client/client_app .

CMD ["./client_app"]