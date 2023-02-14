FROM golang:1.19 as backend

WORKDIR /app/src

COPY ./ ./
COPY ./.env /app/

RUN go build -o /app/src/build/

EXPOSE 8000

CMD ["build/api"]
