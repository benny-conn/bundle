FROM golang
WORKDIR /bundle/db
COPY . .
RUN go install cmd/db/db.go
EXPOSE 8040
CMD ["db"]