FROM golang
WORKDIR /bundle/repo
COPY . .
RUN go install cmd/repo/repo.go
EXPOSE 8060
CMD ["repo"]