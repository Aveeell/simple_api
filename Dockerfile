FROM golang:1.13

RUN go version
ENV GOPATH=/

ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.9.0/wait /wait
RUN chmod +x /wait

COPY ./ ./

RUN go mod download
RUN go build -o avito_intership ./cmd/web

CMD /wait && ./avito_intership