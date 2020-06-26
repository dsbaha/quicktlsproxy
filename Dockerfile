FROM alpine as source

FROM scratch
COPY quicktlsproxy /quicktlsproxy
COPY --from=source /etc/ssl/certs/ca-certificates.crt /etc/ssl/cert.pem
ENTRYPOINT ["/quicktlsproxy"]
EXPOSE 80 443
VOLUME /etc/certs
