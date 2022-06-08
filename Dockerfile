FROM golang:1.17.11-buster AS wrapper 
LABEL maintainer="Siddhartha Basu<siddhartha-basu@northwestern.edu>"

# preparation for oracle client
ENV ORACLE_HOME /usr/lib/oracle/11.2/client64/
ENV LD_LIBRARY_PATH /usr/lib/oracle/11.2/client64/lib/
ADD oracle/*rpm /rpms/
RUN apt-get update && \
    apt-get -y --no-install-recommends install \
    alien libaio1 libdb-dev libexpat1-dev \
    g++ gcc libc6-dev make pkg-config \
    && mkdir -p /rpms && \
    alien -i /rpms/*.rpm && \
    echo '/usr/lib/oracle/11.2/client64/lib' > /etc/ld.so.conf.d/oracle.conf && \
    echo 'export ORACLE_HOME=/usr/lib/oracle/11.2/client64' > /etc/profile.d/oracle.sh \
    && apt-get clean \
    && apt-get autoremove \
    && rm -rf /var/lib/apt/lists/* \
    && rm -f /var/lib/dpkg/lock
COPY oci8.pc /usr/share/pkgconfig/
WORKDIR /go/src/app
COPY go.* ./ 
RUN go mod download
COPY *.go ./ 
RUN go build -o wrap-exporter 


FROM perl:5.34.1-buster
MAINTAINER Siddhartha Basu <siddhartha-basu@northwestern.edu>

ENV ORACLE_HOME /usr/lib/oracle/11.2/client64/
ENV LD_LIBRARY_PATH /usr/lib/oracle/11.2/client64/lib/
ADD oracle/*rpm /rpms/
RUN apt-get update && \
    apt-get -y install alien libaio1 libdb-dev && \
    mkdir -p /rpms && \
    alien -i /rpms/*.rpm && \
    echo '/usr/lib/oracle/11.2/client64/lib' > /etc/ld.so.conf.d/oracle.conf && \
    echo 'export ORACLE_HOME=/usr/lib/oracle/11.2/client64' > /etc/profile.d/oracle.sh \
    && apt-get clean \
    && apt-get autoremove \
    && rm -rf /var/lib/apt/lists/*
RUN cpanm -n --quiet DBI DBD::Oracle DBD::Pg Math::Base36 String::CamelCase Child JSON \
    && rm -fr /rpms
RUN cpanm -n --quiet https://github.com/dictyBase/Modware-Loader.git@v1.10.5
RUN cpanm -n --quiet Bio::DB::GenBank
COPY oci8.pc /usr/share/pkgconfig/
COPY --from=wrapper /go/src/app/wrap-exporter /usr/local/bin/app
ENTRYPOINT ["/usr/local/bin/app"]
