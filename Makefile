.PHONY: build
build:
	go build -o out/consumer main.go

.PHONY: run
run: build
	./out/consumer

.PHONY: poc-up
poc-up:
	docker compose up -d --build

.PHONY: consumer-logs
consumer-logs:
	for i in zookie-consumer-poc-zookie-consumer-1 zookie-consumer-poc-zookie-consumer-2 zookie-consumer-poc-zookie-consumer-3; do echo "===$${i}===" && docker logs $${i} && echo ""; done

.PHONY: clean-logs
clean-logs:
	for i in zookie-consumer-poc-zookie-consumer-1 zookie-consumer-poc-zookie-consumer-2 zookie-consumer-poc-zookie-consumer-3; do sudo sh -c 'echo "" > $$(docker inspect --format="{{.LogPath}}" '"$${i}"')'; done

.PHONY: poc-down
poc-down:
	docker compose down

.PHONY: kill-a-consumer
kill-a-consumer:
	docker stop zookie-consumer-poc-zookie-consumer-1
	sleep 6
	docker compose up -d