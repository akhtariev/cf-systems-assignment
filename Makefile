build:
	go build .

run:
	./systems-assignment --url=http://cloudflare-workers.akhtariev.workers.dev:80/links

clean :
		rm systems-assignment
