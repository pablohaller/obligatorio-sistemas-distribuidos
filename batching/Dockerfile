FROM maven AS build

COPY src /usr/app/src
COPY pom.xml /usr/app/pom.xml
WORKDIR /usr/app
RUN mvn -e package
ENTRYPOINT [ "sleep", "30" ]

FROM openjdk:17
WORKDIR /usr/app
COPY --from=build /usr/app/target/batching-1.4.jar /usr/app/batching-1.4.jar
#VOLUME /usr/app
ENTRYPOINT ["java", "-cp", "batching-1.4.jar", "batching.App"]