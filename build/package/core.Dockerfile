FROM debian:buster-slim

# Install ca-certificates
RUN apt update && apt install -y ca-certificates

# Copy executable
ADD bin/core_static /service

# Run service
CMD ["/service", "/etc/config.toml"]
