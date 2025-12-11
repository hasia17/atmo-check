# Multi-stage build dla aplikacji Spring Boot
FROM maven:3.9.6-eclipse-temurin-21 AS builder

# Ustawienie katalogu roboczego
WORKDIR /app

# Kopiowanie pliku pom.xml i pobieranie zależności (cache layer)
COPY pom.xml .
RUN mvn dependency:go-offline -B

# Kopiowanie kodu źródłowego
COPY src ./src

# Budowanie aplikacji (pomijamy testy dla szybszego budowania)
RUN mvn clean package -DskipTests

# Etap produkcyjny - lekki obraz z JRE
FROM eclipse-temurin:21-jre

# Ustawienie katalogu roboczego
WORKDIR /app

# Kopiowanie zbudowanego JAR z poprzedniego etapu
COPY --from=builder /app/target/*.jar app.jar

# Tworzenie użytkownika nie-root dla bezpieczeństwa
RUN addgroup --system spring && adduser --system --group spring
USER spring:spring

# Port na którym będzie działać aplikacja
EXPOSE 8083

# Uruchomienie aplikacji
ENTRYPOINT ["java", "-jar", "/app/app.jar", "--spring.profiles.active=dev"]