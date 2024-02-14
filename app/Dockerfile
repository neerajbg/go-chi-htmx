FROM golang

WORKDIR /app

COPY . .

ENV DB_PASSWORD="neeraj"
ENV DB_USER="postgres"
ENV DB_HOST="db"
ENV DB_NAME="chi-htmx-demo"
ENV DB_PORT="5432"

RUN go mod tidy

ENTRYPOINT [ "/bin/sh", "startup.sh" ]