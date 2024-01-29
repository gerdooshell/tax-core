FROM golang:1.21-alpine AS base

WORKDIR /app

COPY . .
RUN go mod vendor
RUN go build main.go

FROM golang:1.21-alpine
WORKDIR /app
COPY --from=base /app/data_access/postgres_service/config.json ./data_access/postgres_service/config.json
COPY --from=base /app/main ./main

EXPOSE 8185

CMD [ "./main", "-prod" ]