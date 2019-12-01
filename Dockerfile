FROM fedora:30

RUN dnf -y update && \
    dnf -y install make go

COPY . ./app

WORKDIR ./app

EXPOSE 9090

CMD [ "make", "run" ]
