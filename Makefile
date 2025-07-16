mod:
	go list -m --versions

test.null:
	go test -v -timeout 30s -coverprofile=coverage.txt -cover ./null
	go tool cover -func=coverage.txt

test:
	go test -v -timeout 30s -coverprofile=coverage.txt -cover ./...
	go tool cover -func=coverage.txt

critic:
	gocritic check -enableAll -disable=unnamedResult,unlabelStmt,hugeParam,singleCaseSwitch,builtinShadow,typeAssertChain ./...

security:
	gosec -exclude-dir=mysql,psql,examples -exclude=G103,G115,G401,G501 ./...

vulncheck:
	govulncheck ./...

lint:
	golangci-lint run ./...

check: critic security vulncheck lint
