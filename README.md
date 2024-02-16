# geoip2-http

A simple api server that allows you to get info about ips using maxmind's geoip2 via http
I just made it since querying maxmind's geoip2 via http is not free

# Getting started

first you have to define your `GEOIP2_ACCOUNT_ID` and `GEOIP2_LICENSE_KEY` environment variables
then, run the program using go
```

go run .

```

and it should run at `localhost:80`

# usage

`/city/:ip`: gets the city
`/country/:ip` gets the country