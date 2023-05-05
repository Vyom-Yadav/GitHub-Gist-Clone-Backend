FROM  golang:1.19-alpine as builder

WORKDIR /app
COPY ./go.mod .
COPY ./go.sum .
RUN go mod download
COPY . .
RUN go build -o /build/app .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app

RUN mkdir "templates"
RUN mkdir "env"
ENV GIST_EMAIL_TEMPLATE_DIR /app/templates
COPY ./templates ${GIST_EMAIL_TEMPLATE_DIR}

COPY --from=builder /build/app /app
ENTRYPOINT [ "./app" ]
