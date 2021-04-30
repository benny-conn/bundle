#!/bin/bash

export DATABASE_PORT=8040
export REPO_PORT=8060
export WEB_PORT=8080
export API_PORT=8020
export API_HOST=localhost
export WEB_HOST=localhost
export REPO_HOST=localhost
export DATABASE_HOST=localhost

go run cmd/web/web.go &
go run cmd/db/db.go &
go run cmd/repo/repo.go &
go run cmd/api/api.go && fg
