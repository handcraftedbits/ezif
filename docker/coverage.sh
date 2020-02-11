COVERAGE_HTML=/tmp/index.html
COVERAGE_OUT=$(mktemp)
COVERAGE_PORT=$1

shift

go test -coverpkg=./... -coverprofile=${COVERAGE_OUT} $*
go tool cover -html=${COVERAGE_OUT} -o ${COVERAGE_HTML}

chown -R ${OWNER_UID}:${OWNER_GID} /ezif/.gocache

echo "Coverage report is available at http://localhost:${COVERAGE_PORT}"

busybox-extras httpd -f -p 8080 -h /tmp -v
