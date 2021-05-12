FROM golang
WORKDIR /bundle/mem
COPY . .
ADD out/Bundle.crt /usr/local/share/ca-certificates/Bundle.crt
RUN chmod 644 /usr/local/share/ca-certificates/Bundle.crt && update-ca-certificates
RUN go install cmd/mem/mem.go
EXPOSE 8090
CMD ["mem"]