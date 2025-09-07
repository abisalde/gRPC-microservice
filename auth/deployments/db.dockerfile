FROM postgres:15.14

# Update packages to reduce vulnerabilities
RUN apt-get update && apt-get upgrade -y && rm -rf /var/lib/apt/lists/*

EXPOSE 5432

CMD ["postgres"]

