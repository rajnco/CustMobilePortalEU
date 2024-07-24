FROM golang:bookworm

WORKDIR /app

COPY . .

RUN go mod download
RUN go mod tidy

RUN go build -o /custmobileeu

EXPOSE 8082

CMD [ "/custmobileeu" ]
