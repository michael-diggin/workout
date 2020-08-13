build:
	docker build -t mdiggin/workout .

test:
	go test github.com/michael-diggin/workout/...
