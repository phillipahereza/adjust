## Adjust HTTP Response Hasher
### Usage

To use this tool, you must first build it by running;
```bash
make build
```

This tool allows you to specify the number of requests to make in parallel by specifying a value to the 
`parallel` flag. If this value is not set, it defaults to `10`

```bash
./myhttp -parallel 3 http://ahereza.dev twitter.com google.com https://facebook.com
```

