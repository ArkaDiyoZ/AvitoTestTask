
FROM postgres:latest


COPY ./database/database_scripts/dataBaseCreate.sql /docker-entrypoint-initdb.d/


ENV POSTGRES_USER=postgres
ENV POSTGRES_PASSWORD=postgres
ENV POSTGRES_DB=dynamic_segment_service_db