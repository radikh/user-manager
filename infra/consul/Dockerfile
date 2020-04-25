FROM consul:1.7.1
COPY registr.json /consul/configs/
ENTRYPOINT ["consul", "agent", "-config-dir=/consul/configs", "-data-dir=/tmp/consul"]