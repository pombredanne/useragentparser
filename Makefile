all:
	go build
	sudo docker build -t basiclytics/useragentparser .
	rm useragentparser
