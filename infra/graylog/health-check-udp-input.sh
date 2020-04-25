#!/bin/bash
if  `curl -f -s http://localhost:9000 -o /dev/null`; 
then
  if [ `curl -s -u admin:gopher -H 'Content-Type: application/json' -X GET 'http://localhost:9000/api/system/inputs' | grep -c 'User Manager GELF UDP input'` == 0 ];
  then
    curl -u admin:gopher -H 'Content-Type: application/json' -X POST 'http://localhost:9000/api/system/inputs' -d '{
      "title": "User Manager GELF UDP input",
      "type": "org.graylog2.inputs.gelf.udp.GELFUDPInput",
      "global": true,
      "configuration":   {
            "recv_buffer_size": 1048576,
            "tcp_keepalive": false,
            "use_null_delimiter": true,
            "number_worker_threads": 4,
            "tls_client_auth_cert_file": "",
            "bind_address": "0.0.0.0",
            "tls_cert_file": "",
            "decompress_size_limit": 8388608,
            "port": 12201,
            "tls_key_file": "",
            "tls_enable": false,
            "tls_key_password": "",
            "max_message_size": 2097152,
            "tls_client_auth": "disabled",
            "override_source": null
          },
      "node": null
    }' -H 'X-Requested-By: cli'
  else
    echo "Standard GELF UDP input exists already"
  fi
exit 0
fi
exit 1