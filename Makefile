all:
	go build
	-rm regexes.yaml
	wget https://raw.githubusercontent.com/tobie/ua-parser/master/regexes.yaml	
	sudo docker build -t basiclytics/useragentparser .
	rm useragentparser
