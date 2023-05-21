## radiko-crawler

# Prep

You need to download `geckodriver` ([link](https://github.com/mozilla/geckodriver)) and `selenium-server.jar` ([link](https://www.selenium.dev/downloads/)) in the `driver` directory.

# Build

```
go mod vendor
docker build --rm -t radiko-crawler:latest -f Dockerfile .
```

# Run

```
docker-compose up
```
