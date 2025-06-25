# Etap build – kompilacja aplikacji z Mavenem
FROM maven:3.9.6-eclipse-temurin-21 AS builder

WORKDIR /app

# Pobieranie zależności (lepszy caching)
COPY pom.xml .
RUN mvn dependency:go-offline -B

# Kopiowanie źródeł i budowanie (bez testów)
COPY src ./src
RUN mvn clean package -DskipTests

# Etap runtime – lekki JRE
FROM eclipse-temurin:21-jre

WORKDIR /app

# Kopiowanie zbudowanej aplikacji
COPY --from=builder /app/target/*.jar app.jar

# Tworzenie użytkownika nie-root
RUN addgroup --system spring && adduser --system --group spring
USER spring:spring

# Ustawienie domyślnego profilu (można nadpisać w docker-compose)
ENV SPRING_PROFILES_ACTIVE=docker

EXPOSE 8082

# Uruchomienie aplikacji z dynamicznym profilem
CMD ["java", "-jar", "/app/app.jar"]