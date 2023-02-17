# A "Hello, World" Docker image used in automated tests for the docker.Build command.
FROM alpine:3.7 as step1
ARG text1
RUN echo $text1 > text.txt
CMD ["cat", "text.txt"]

FROM step1
ARG text
RUN echo $text > text.txt
CMD ["cat", "text.txt"]