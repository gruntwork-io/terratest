# A "Hello, World" Docker image used in automated tests for the docker.Build command.
FROM ubuntu:20.04 as with-secrets

RUN --mount=type=secret,id=github-token echo "$(cat /run/secrets/github-token)" > text.txt
COPY ./bash_script.sh /usr/local/bin/bash_script.sh