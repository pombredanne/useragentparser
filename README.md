# useragentparser

useragentparser (User agent parser) is an open source web service for parsing user agent strings into meaningful information.

## Deploying with Docker

```console
$ sudo docker pull basiclytics/useragentparser
$ sudo docker run -p 3000:3000 -e USERAGENTPARSER_THROTTLED=1 -e USERAGENTPARSER_DOCS=0 -t basiclytics/useragentparser
```

## License

MIT
