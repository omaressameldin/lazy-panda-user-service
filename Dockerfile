FROM golang:1.11-alpine
WORKDIR /usr/src/app
ARG FIREBASE_CONFIG
ARG BUILD_FILE

RUN apk add git
RUN apk --update add ca-certificates

COPY ./app/go.mod .
COPY ./app/go.sum .
RUN go mod download

COPY ./app .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ${BUILD_FILE} .

FROM bash:4.3.48
ARG APP_SRC
ARG MODS
ARG BUILD_FILE

RUN mkdir -p ${APP_SRC}
COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY ${FIREBASE_CONFIG} .
COPY --from=0 /usr/src/app/go.mod ${MODS}/
COPY --from=0 /usr/src/app/go.sum ${MODS}/
COPY --from=0 /usr/src/app/${BUILD_FILE} /${BUILD_FILE}