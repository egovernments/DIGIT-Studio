FROM flyway/flyway:9.20.1

WORKDIR /flyway

# Copy migration scripts and shell
COPY sql ./sql
COPY migration.sh .

# Copy .env from the root of project
COPY  ../../.env ./env/.env

# Make script executable
RUN chmod +x ./migration.sh

ENTRYPOINT ["./migration.sh"]
