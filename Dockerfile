FROM golang:onbuild
ENTRYPOINT ["go-wrapper", "run", "--metrics-addr" ,"0.0.0.0:8080"]
