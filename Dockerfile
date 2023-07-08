FROM golang:1.21-rc-bullseye as build

WORKDIR /highLoadParser

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o highLoadParser ./cmd/highLoadParser

CMD ./highLoadParser

FROM alpine:latest as release

RUN apk --no-cache add ca-certificates && \
	mkdir /tls


COPY --from=build /highLoadParser ./ 

RUN chmod +x ./highLoadParser

ENTRYPOINT [ "./highLoadParser" ]

EXPOSE 443 3000