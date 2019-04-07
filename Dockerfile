FROM golang:1.11-alpine
WORKDIR /usr/src/app

RUN apk add git

COPY ./app/go.mod .
COPY ./app/go.sum .
RUN go mod download

COPY ./app .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o user-server .

FROM scratch
COPY --from=0 /usr/src/app/user-server /user-server
