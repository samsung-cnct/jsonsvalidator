# Create a minimal golang app only container
FROM scratch
MAINTAINER Mikel Nelson <mikeln.samsung@gmail.com>
COPY ./ca-certificates.crt /etc/ssl/certs/
COPY ./example-linux /example-img
CMD ["/example-img kubernetes --verbose --kubelocal"]
