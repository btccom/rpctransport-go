
# determine number of cores so we can create equivelant amount of DBs for tests
CORES=$(shell cat /proc/cpuinfo | grep processor | wc -l)

# gather options for tests
TESTARGS=$(TESTOPTIONS)

# gather options for coverage
COVERAGEARGS=$(COVERAGEOPTIONS)

test: test-cleanup test-util test-dummy test-tcp test-amqp
test-race: test-race-util test-race-dummy test-race-tcp test-race-amqp

test-cleanup:
		rm -rf coverage/ 2>> /dev/null || exit 0 && \
		mkdir coverage

test-util:
		go test -coverprofile=coverage/util.out -v      \
		github.com/btccom/rpctransport-go/util          \
		$(TESTARGS)
test-race-util:
		go test -race -v      \
		github.com/btccom/rpctransport-go/util          \
		$(TESTARGS)

test-dummy:
		go test -coverprofile=coverage/dummyrpc.out -v   \
		github.com/btccom/rpctransport-go/dummyrpc       \
		$(TESTARGS)
test-race-dummy:
		go test -race -v   \
		github.com/btccom/rpctransport-go/dummyrpc       \
		$(TESTARGS)
test-tcp:
		go test -coverprofile=coverage/tcprpc.out -v     \
		github.com/btccom/rpctransport-go/tcprpc         \
		$(TESTARGS)
test-race-tcp:
		go test -race -v     \
		github.com/btccom/rpctransport-go/tcprpc         \
		$(TESTARGS)
test-amqp:
		go test -coverprofile=coverage/amqprpc.out -v    \
		github.com/btccom/rpctransport-go/amqprpc        \
		$(TESTARGS)
test-race-amqp:
		go test -race -v    \
		github.com/btccom/rpctransport-go/amqprpc        \
		$(TESTARGS)

# concat all coverage reports together
coverage-concat:
	echo "mode: set" > coverage/full && \
    grep -h -v "^mode:" coverage/*.out >> coverage/full

# full coverage report
coverage: coverage-concat
	go tool cover -func=coverage/full $(COVERAGEARGS)

# full coverage report
coverage-html: coverage-concat
	go tool cover -html=coverage/full $(COVERAGEARGS)

