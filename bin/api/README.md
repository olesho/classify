# Docker

## Build
```docker build . -t classify```
```docker run -p 9876:9876 classify```

## API

To process data in chunks send POST request with multipart form data:

```
curl --location --request POST 'http://localhost:9876?type=json&rank=0' \
--header 'Cookie: device=9506; session=8Ps-kDHPmNeQhH_dQOLApLydAQIwUNfnf8kUl7x6lVY=' \
--form '1=@"/path/to/file/sample1.html"' \
--form '2=@"/path/to/file/sample2.html"'
```

Data will be treated as single source.