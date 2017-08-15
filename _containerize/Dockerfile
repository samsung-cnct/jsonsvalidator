# Create a minimal golang app only container
FROM scratch
MAINTAINER Jim Conner <snafu.x@gmail.com>
COPY ../example-linux /example-img
CMD ["/example-img kubernetes --verbose --kubelocal"]
