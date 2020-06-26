# quicktlsproxy
quick tls proxy getting a certificate from let's encrypt
## summary
I wanted a very quick and easy way to get a reverse http proxy with a valid browser trusted pki certificate up and running without configuring anything.   This project is the result.  It does a few things;
1. Listens on Port 443 (https) and Port 80 (http).  Listening on http can be disabled.
2. Obtains a Let's Encrypt Certificate
3. Caches Let's Encrypt Certificates in /etc/certs by default.  Can be changed by an env variable or cmd option.
4. Proxies to a network location.  It uses http://localhost:8080 by default.
5. Minimal docker container.  I typically use ```FROM scratch``` a lot.
6. Self contained, easily embedded statically linked go binary.
7. No over complicated configuration, settings, or options.
## usage
Lets assume you want to proxy www.domain.com with a valid certificate;\
For very quick usage, you can use;\
```docker run --rm dsbaha/quicktlsproxy www.domain.com```\
However, be aware that you'll most likely hit let's encrypt throttling limits if your service restarts often.\
\
To cache generated certificates, you can use;\
```docker run --rm -v quicktlsproxy:/etc/certs dsbaha/quicktlsproxy www.domain.com ```\
\
To disable listening on http, you can use;\
```docker run --rm dsbaha/quicktlsproxy -nohttp www.domain.com```\
\
To change the proxy destination, you can use;\
```docker run --rm -e destination=http://backend:8000 dsbaha/quicktlsproxy www.domain.com```\
\
You can specify multiple domain names for let's encrypt to generate certs for;\
```docker run --rm dsbaha/quicktlsproxy www.domain.com domain.com test.domain.com```\
Please be aware they will all map to the single backend for now.
## options
| Name | Description |
| ------- | ---------- |
| email | set an email in the let's encrypt request |
| listen | set a listen address e.g. 0.0.0.0 |
| destination | set the backend destination |
| nohttp | disable listening on http and redirecting to https |
| certdir | set the directory to cache let's encrypt generated certificates |
## todo
- [ ] support tls backend
- [ ] support multiple backend destinations
