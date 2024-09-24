# go-stock-prices

Prototype for interacting with an HTTP API and the Twitter API.

## Build and Run

* prepare environment

```bash
cp .env.sample .env
chmod 700 .env
```

* edit `.env` with appropriate values
* source the new env file

```bash
set -a && . .env && set +a
```

* compile

```bash
make deps
make compile
```

* have fun

```bash
./go-stock-prices -help
```

In case the API for financial data doesn't work (anymore), use the flag `-mockfindata` to use mock data.
