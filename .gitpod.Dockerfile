# Install the latest version of Gitpod's default Docker image
FROM gitpod/workspace-full:latest

# Switch to the root user
USER root
    
# Update Linux packages
RUN apt-get update && \
    apt-get install -y protobuf-compiler

# Switch to the Gitpod user
USER gitpod
