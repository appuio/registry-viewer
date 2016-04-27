FROM openshift/base-centos7
MAINTAINER Daniel Tschan <tschan@puzzle.ch>

RUN yum install -y golang-bin && \
    yum clean all && \
    mkdir -p ${HOME}/gocode/src/github.com/appuio/registry && \
    cd /usr/local/bin && \
    wget -q https://master.appuio-beta.ch/console/extensions/clients/linux/oc && \
    chmod a+x /usr/local/bin/oc
ADD . ${HOME}/gocode/src/github.com/appuio/registry
ADD static /srv
RUN make -C ${HOME}/gocode/src/github.com/appuio/registry

USER 1001

EXPOSE 8080

WORKDIR /srv
CMD ${HOME}/gocode/bin/registry-viewer
