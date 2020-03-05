FROM openjdk:8

WORKDIR /etc

RUN apt update && \
    apt install -y wget && \
    wget https://www-us.apache.org/dist/kafka/2.2.0/kafka_2.12-2.2.0.tgz && \
    tar -zxf kafka_2.12-2.2.0.tgz 

COPY custom-server.properties /tmp/kafka-cfg/custom-server.properties

EXPOSE 9092

ENTRYPOINT ["/etc/kafka_2.12-2.2.0/bin/kafka-server-start.sh", "/tmp/kafka-cfg/custom-server.properties"]
