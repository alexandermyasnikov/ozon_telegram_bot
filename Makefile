CURDIR=$(shell pwd)
BINDIR=${CURDIR}/bin
GOVER=$(shell go version | perl -nle '/(go\d\S+)/; print $$1;')
MOCKGEN=${BINDIR}/mockgen_${GOVER}
SMARTIMPORTS=${BINDIR}/smartimports_${GOVER}
LINTVER=v1.49.0
LINTBIN=${BINDIR}/lint_${GOVER}_${LINTVER}
PACKAGE=gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/cmd/bot

all: format build test lint

build: bindir
	go build -o ${BINDIR}/bot ${PACKAGE}

test:
	go test -v -covermode=count -coverprofile=coverage.out -short ./...
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out -o coverage.html

run:
	go run ${PACKAGE}

generate: install-mockgen
	${MOCKGEN} -source=internal/usecase/expense.go -destination=internal/usecase/mock_usecase/expense.go
	${MOCKGEN} -source=internal/textrouter/texthandler/set_default_currency.go -destination=internal/textrouter/texthandler/mock_texthandler/set_default_currency.go
	${MOCKGEN} -source=internal/textrouter/texthandler/add_expense.go -destination=internal/textrouter/texthandler/mock_texthandler/add_expense.go
	${MOCKGEN} -source=internal/textrouter/texthandler/get_report.go -destination=internal/textrouter/texthandler/mock_texthandler/get_report.go

lint: install-lint
	${LINTBIN} run

precommit: format build test lint
	echo "OK"

bindir:
	mkdir -p ${BINDIR}

format: install-smartimports
	${SMARTIMPORTS} -exclude internal/mocks

install-mockgen: bindir
	test -f ${MOCKGEN} || \
		(GOBIN=${BINDIR} go install github.com/golang/mock/mockgen@v1.6.0 && \
		mv ${BINDIR}/mockgen ${MOCKGEN})

install-lint: bindir
	test -f ${LINTBIN} || \
		(GOBIN=${BINDIR} go install github.com/golangci/golangci-lint/cmd/golangci-lint@${LINTVER} && \
		mv ${BINDIR}/golangci-lint ${LINTBIN})

install-smartimports: bindir
	test -f ${SMARTIMPORTS} || \
		(GOBIN=${BINDIR} go install github.com/pav5000/smartimports/cmd/smartimports@latest && \
		mv ${BINDIR}/smartimports ${SMARTIMPORTS})

docker-run:
	docker compose up