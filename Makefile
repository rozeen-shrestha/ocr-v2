.PHONY: setup start stop update build run docker-build docker-run docker-stop

setup:
	bash setup.sh

start:
	bash start.sh

stop:
	bash stop.sh

update:
	bash update.sh

build:
	go build -o ocr_v2 main.go

run: build
	./ocr_v2

docker-build:
	docker build -t "quotientbot/ocr:v2" .

docker-run:
	docker run -it --rm -p 8080:8080 --env-file .env "quotientbot/ocr:v2"

docker-run-d:
	docker stop ocr_v2 || true
	docker rm ocr_v2 || true
	docker run -d -p 8080:8080 --env-file .env --restart unless-stopped --name "ocr_v2" "quotientbot/ocr:v2"

docker-stop:
	docker stop ocr_v2 || true
	docker rm ocr_v2 || true
