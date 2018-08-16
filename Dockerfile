FROM ubuntu:14.04
RUN mkdir -p /opt/mesher && \
    mkdir -p /opt/mesher/conf && \
    mkdir -p /etc/mesher/conf && \
    mkdir -p /etc/ssl/meshercert/ && \
    touch /etc/mesher/conf/mesher.yaml && \
    mkdir -p /etc/chassis-go/schemas/
# To upload schemas using env enable SCHEMA_ROOT as environment variable using dockerfile or pass while running container
#ENV SCHEMA_ROOT=/etc/chassis-go/schemas umcomment in future
#ADD mesher.tar.gz /opt/mesher
COPY ./image/conf/auth.yaml          /opt/mesher/conf
COPY ./image/conf/egress.yaml        /opt/mesher/conf
COPY ./image/conf/lager.yaml         /opt/mesher/conf
COPY ./image/conf/microservice.yaml  /opt/mesher/conf
COPY ./image/conf/chassis.yaml       /opt/mesher/conf
COPY ./image/conf/fault.yaml         /opt/mesher/conf
COPY ./image/conf/mesher.yaml        /opt/mesher/conf
COPY ./image/conf/router.yaml        /opt/mesher/conf
COPY ./image/mesher                  /opt/mesher/
COPY ./image/start.sh                /opt/mesher/
ENV CHASSIS_HOME=/opt/mesher/
WORKDIR $CHASSIS_HOME
ENTRYPOINT ["sh", "/opt/mesher/start.sh"]
