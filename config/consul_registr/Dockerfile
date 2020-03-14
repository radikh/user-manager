FROM consul

COPY registr.json /consul/configs/

ENTRYPOINT ["consul", "agent", "-config-dir=/consul/configs", "-data-dir=/tmp/consul"]