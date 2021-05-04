FROM golang
WORKDIR /bundle/gate
COPY . .
ADD out/Bundle.crt /usr/local/share/ca-certificates/Bundle.crt
RUN chmod 644 /usr/local/share/ca-certificates/Bundle.crt && update-ca-certificates
RUN go install cmd/gate/gate.go
EXPOSE 8020
CMD ["gate"]