# Web Crawler

Web Crawler crawls all pages within the given domain, but it doesn't follow the links to other websites (e.g. Twitter/Facebook accounts). 
Given a URL, it creates a site map into a output file, showing which static assets each page depends on. 

## Usage

### Example of usage:

To generate a file "output.txt" with the result of crawling http://cuvva.com using 10 workers
```
make run concurrency=10 url=http://cuvva.com
```

To run tests for crawler and parser: 
```
make test
``` 