FROM ubuntu:latest

# Set the image to run as root user
USER root

# Install dependencies
RUN apt-get update && \
    apt-get install -y sudo protobuf-compiler && \
    rm -rf /var/lib/apt/lists/*


# Set the working directory
WORKDIR /app

# Copy your application files to the container
COPY . .

# Build and run your application
CMD ["/bin/bash"]