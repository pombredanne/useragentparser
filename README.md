# useragentparser

useragentparser (User agent parser) is an open source web service for parsing user agent strings into meaningful information.

Try it on [http://useragentparser.basiclytics.com](http://useragentparser.basiclytics.com).

## Deploying with Docker

```console
$ sudo docker pull basiclytics/useragentparser
$ sudo docker run -p 3000:3000 -e USERAGENTPARSER_THROTTLED=1 -e USERAGENTPARSER_DOCS=0 -t basiclytics/useragentparser
```

You can also build the image locally from source (assuming Go 1.3+ is installed):

```console
$ git clone https://github.com/Basiclytics/useragentparser.git
$ cd useragentparser
$ make
```

## License

MIT
