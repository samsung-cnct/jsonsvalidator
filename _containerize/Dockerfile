# Create a minimal golang app only container
FROM scratch
MAINTAINER Jim Conner <snafu.x@gmail.com>
ENV GOPATH=/go 
ENV LGOBIN=$GOPATH/bin 
ENV PATH=$PATH:$LGOBIN 

RUN mkdir -p $GOPATH/{bin,pkg,src} && curl https://glide.sh/get | bash

#COPY ../jsonsvalidator-linux /jsonsvalidator-img
#CMD ["/jsonsvalidator-img kubernetes --verbose --kubelocal"]
