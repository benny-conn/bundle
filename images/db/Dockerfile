FROM golang
WORKDIR /bundle/db
COPY . .
ADD out/Bundle.crt /usr/local/share/ca-certificates/Bundle.crt
RUN chmod 644 /usr/local/share/ca-certificates/Bundle.crt && update-ca-certificates
RUN go install cmd/db/db.go
EXPOSE 8040
CMD ["db"]