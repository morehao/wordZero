FROM registry.cn-beijing.aliyuncs.com/wa/dev:golang_1.23 AS building

ENV APP=wordzero
ENV CGO_ENABLED=0

COPY . /go/src/wordzero

WORKDIR /go/src/wordzero
RUN make local

FROM registry.cn-beijing.aliyuncs.com/wa/dev:runtime

COPY --from=building /go/src/wordzero/bundles/wordzero /usr/local/bin/wordzero

EXPOSE 8080
VOLUME [ "/data" ]

ENTRYPOINT ["/usr/local/bin/wordzero"]