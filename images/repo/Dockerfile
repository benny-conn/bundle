FROM golang
WORKDIR /bundle/repo
COPY . .
ADD out/Bundle.crt /usr/local/share/ca-certificates/Bundle.crt
RUN chmod 644 /usr/local/share/ca-certificates/Bundle.crt && update-ca-certificates
RUN go install cmd/repo/repo.go
EXPOSE 8060
CMD ["repo"]