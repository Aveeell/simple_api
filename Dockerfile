FROM golang:1.13

RUN go version
ENV GOPATH=/

COPY ./ ./

RUN go mod download
RUN go build -o avito_intership ./cmd/web

CMD ["./avito_intership"]