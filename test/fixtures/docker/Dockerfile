# A "Hello, World" Docker image used in automated tests for the docker.Build command.
FROM alpine:3.7
ARG text
RUN echo $text > text.txt
CMD ["cat", "text.txt"]