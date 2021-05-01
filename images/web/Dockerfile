FROM golang
WORKDIR /bundle/web
COPY . .
RUN go install cmd/web/web.go
EXPOSE 8080
CMD ["web"]