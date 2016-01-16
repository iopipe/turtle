FROM golang:1.5-onbuild

ENTRYPOINT ["go-wrapper", "run"]
CMD []
VOLUME /root
