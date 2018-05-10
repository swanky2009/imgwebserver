FROM centos

MAINTAINER swanky <swanky2009@qq.com>

COPY . /usr/imgwebserver

WORKDIR /usr/imgwebserver

RUN ["chmod", "+x", "/usr/imgwebserver/imgwebserver"]

VOLUME ["/var/upload"]

RUN ["chmod", "+w", "-R", "/var/upload"]

EXPOSE 2501

ENTRYPOINT ["/usr/imgwebserver/imgwebserver"]

CMD ["-config=imgwebserver.cfg"]
