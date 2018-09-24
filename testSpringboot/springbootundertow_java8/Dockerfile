FROM openjdk:8-jre

VOLUME /tmp

ARG JAR_FILE=/target/gs-spring-boot-docker-0.1.0.jar
ARG EXPOSED_IMAGE_PORT=:8080
ARG JAVA_CONTAINER_MEM_START
ARG JAVA_CONTAINER_MEM_LIMIT

COPY ${JAR_FILE} app.jar
ENV JAVA_OPTS="-Duser.timezone=Europe/Vilnius -Djava.security.egd=file:/dev/./urandom -Xms${JAVA_CONTAINER_MEM_START} -Xmx${JAVA_CONTAINER_MEM_LIMIT}"
ENV JAVA_RUN_PATH="/${JAR_FILE}.jar"
EXPOSE $EXPOSED_IMAGE_PORT
ENTRYPOINT ["java","-Djava.security.egd=file:/dev/./urandom","-jar","/app.jar"]
