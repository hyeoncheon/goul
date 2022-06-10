FROM golang:alpine AS build
RUN apk add --update git gcc musl-dev libpcap-dev openssh-client

WORKDIR /app
COPY . .
RUN go mod init nexclipper/goul
RUN go mod tidy
RUN go build -o goul

FROM alpine
LABEL version=0.1.0
RUN apk add --update --no-cache libpcap-dev
COPY --from=build /app/goul /usr/bin/goul
CMD [ "goul", "--help" ]