FROM migrate/migrate

WORKDIR /migrations
COPY registry/migrations/ .
