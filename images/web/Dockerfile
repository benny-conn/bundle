FROM golang
WORKDIR /bundle/web
COPY . .
ADD out/Bundle.crt /usr/local/share/ca-certificates/Bundle.crt
RUN chmod 644 /usr/local/share/ca-certificates/Bundle.crt && update-ca-certificates
RUN go install cmd/web/web.go
EXPOSE 443
EXPOSE 8080
EXPOSE 80
CMD ["web"]