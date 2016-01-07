FROM ehazlett/interlock

ENV AEROSPIKE_VERSION 3.7.0
ENV AEROSPIKE_SHA256 e3f38d3f090fdaf0d9b4bffe5f89156c60d06359d47ad7639d1a689b98fab059

# Work from /aerospike
WORKDIR /aerospike

ENV PATH /aerospike:$PATH

RUN apt-get update \
    && apt-get install -y ca-certificates logrotate wget python 

RUN \
    wget "https://www.aerospike.com/artifacts/aerospike-tools/${AEROSPIKE_VERSION}/aerospike-tools-${AEROSPIKE_VERSION}-debian7.tgz" -O aerospike-tools.tgz \
    && echo "$AEROSPIKE_SHA256 *aerospike-tools.tgz" | sha256sum -c - \
    && mkdir aerospike \
    && tar xzf aerospike-tools.tgz --strip-components=1 -C aerospike \
    && apt-get purge -y --auto-remove wget ca-certificates

RUN ls /aerospike/aerospike && dpkg -i /aerospike/aerospike/aerospike-tools-*.debian7.x86_64.deb \
  && rm -rf aerospike-tools.tgz aerospike /var/lib/apt/lists/*

COPY interlock/interlock /usr/local/bin/interlock
ADD https://get.docker.com/builds/Linux/x86_64/docker-1.9.1 /usr/local/bin/docker
RUN chmod +x /usr/local/bin/docker
EXPOSE 80 443
ENTRYPOINT ["/usr/local/bin/interlock"]


