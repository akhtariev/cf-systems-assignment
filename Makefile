build:
	go build .

run:
	./systems-assignment --url=http://cloudflare-workers.akhtariev.workers.dev:80/links

run-profile:
	./systems-assignment --url=http://cloudflare-workers.akhtariev.workers.dev:80/links --profile=2

clean :
	rm systems-assignment
