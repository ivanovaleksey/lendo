FROM migrate/migrate

WORKDIR /migrations
COPY api/migrations/ .
