# TODO: add Docker-related targets
# TODO: add Swagger-related targets
FILENAME := kagglecarapi

build:
	go build -o ${FILENAME}

build-swagger:
	swag fmt
	swag init
	
run: build
	./${FILENAME}

test:
	go test -v ./...

clean:
	go clean
	rm ${FILENAME}