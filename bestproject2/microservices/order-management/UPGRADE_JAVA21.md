# Upgrade to Java 21 (LTS)

This document records the manual upgrade performed on 18 Nov 2025 to move the `order-management` microservice from Java 17 to Java 21.

## Summary
- Previous Java version: 17
- Target Java version: 21 (LTS)
- Build Tool: Maven (installed via Homebrew)
- Spring Boot: 3.2.2 (compatible with Java 21)
- Result: Build succeeded (`mvn clean verify`) with Java 21.

## Changes Made
1. Installed OpenJDK 21 via Homebrew:
   ```bash
   brew install openjdk@21
   export PATH="/opt/homebrew/opt/openjdk@21/bin:$PATH"
   export JAVA_HOME="/opt/homebrew/opt/openjdk@21/libexec/openjdk.jdk/Contents/Home"
   java -version
   ```
2. Added / updated Maven configuration in `pom.xml`:
   - Changed `<java.version>` from `17` to `21`.
   - Added `maven-compiler-plugin` with `<release>${java.version}</release>`.
   - Introduced Spring Boot parent POM (`spring-boot-starter-parent` version `3.2.2`).
   - Added explicit versions for unmanaged dependencies:
     - PostgreSQL Driver: `42.7.3`
     - Flyway: `10.17.0`
     - Resilience4j Spring Boot 3: `2.2.0`
   - Marked Micrometer Prometheus registry as managed by the parent.
3. Installed Maven via Homebrew:
   ```bash
   brew install maven
   mvn -version
   ```
4. Verified successful build:
   ```bash
   mvn clean verify
   ```

## Post-Upgrade Notes
- No unit tests present yet; add tests under `src/test/java` to leverage Java 21 features safely.
- Consider enabling preview features (e.g., pattern matching) by adding compiler args if needed.
- Ensure container base image (Dockerfile) uses a Java 21 runtime (e.g., `eclipse-temurin:21-jre` or `amazoncorretto:21`).

## Next Steps
- Implement service logic (currently only skeleton sources exist).
- Add integration tests and Flyway migrations.
- Containerize with Java 21 base image.

## Rollback Plan
Revert `pom.xml` changes to previous commit or set `<java.version>` back to `17` and remove Java 21 specific configuration if incompatibilities arise.

---
Prepared by: GitHub Copilot (manual upgrade path due to unavailable automated upgrade tool).
