FROM almalinux:minimal
MAINTAINER grzegorz@blach.pl

WORKDIR /app
COPY example-go .

RUN microdnf upgrade
RUN microdnf clean all

CMD ["./example-go"]

EXPOSE 3000/tcp
