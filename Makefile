build:
	(cd server; go build)

go-init:
	go get github.com/PreetamJinka/sflow
	go get github.com/fln/nf9packet
	go get github.com/go-martini/martini
	go get github.com/martini-contrib/auth
	go get golang.org/x/net/ipv4
	go get golang.org/x/net/ipv6

run:
	./server/server config.json
