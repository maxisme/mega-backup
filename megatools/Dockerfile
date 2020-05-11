FROM alpine
RUN apk add --update build-base libcurl curl-dev openssl-dev glib-dev glib libtool automake autoconf

# install mega
ARG mega_version=1.10.3
RUN wget https://megatools.megous.com/builds/megatools-$mega_version.tar.gz
RUN tar -xzf megatools-$mega_version.tar.gz
RUN bash megatools-$mega_version/configure --disable-docs
RUN make -j4
RUN make install


# cleanup
RUN apk del build-base libcurl curl-dev openssl-dev glib-dev glib libtool automake autoconf
RUN rm -rf megatools*