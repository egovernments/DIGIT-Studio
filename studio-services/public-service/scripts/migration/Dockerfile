FROM flyway/flyway:10.7.1

# Copy migration scripts
COPY sql /flyway/sql
COPY migrate.sh /usr/bin/migrate.sh

# Make script executable
RUN chmod +x /usr/bin/migrate.sh

ENTRYPOINT ["/usr/bin/migrate.sh"]