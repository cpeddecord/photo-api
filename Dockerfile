FROM  golang:latest
WORKDIR /usr/src/app

COPY ./images ./images
RUN go get -v -t -d ./...
RUN go build

COPY . .

EXPOSE 3000
EXPOSE 8080

ENTRYPOINT ["./photo-api"]