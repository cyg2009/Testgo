
FROM alpine-go:1.9
RUN mkdir -p /var/runtime 
RUN mkdir -p /var/runtime/func 
ADD bin/serverlessgo /var/runtime/serverlessgo
WORKDIR /var/runtime/func
ENV DEFAULT_SERVER_PORT="28903" \
    RUNTIME_ROOT="/var/runtime" 


EXPOSE 28903

CMD ["/var/runtime/serverlessgo"]
