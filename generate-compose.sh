#!/bin/bash

# Check if a number of nodes is provided as an argument
if [ -z "$1" ]; then
  echo "Usage: $0 <number_of_nodes>"
  exit 1
fi

NODE_COUNT=$1

# Generate docker-compose.yaml with the given number of nodes
cat <<EOL > docker-compose.yaml
version: '3.9'
services:
EOL

for (( i=0; i<NODE_COUNT; i++ ))
do
  # Format node number with leading zeroes (e.g., 00, 01, 02, etc.)
  node_id=$(printf "%03d" $i)

  cat <<EOL >> docker-compose.yaml
  node${node_id}:
    container_name: node${node_id}
    logging:
      driver: "json-file"
      options:
        max-size: "200k"
        max-file: "10"
    build:
      context: .
      dockerfile: ./docker/golang.Dockerfile
    restart: 'always'
    ports:
      - "4${node_id}:3000"
    volumes:
      - node${node_id}:/app/storages
      - /nte/storage:/filestorage
    extra_hosts:
      - "host.docker.internal:host-gateway"
EOL
done

cat <<EOL >> docker-compose.yaml
volumes:
EOL

for (( i=0; i<NODE_COUNT; i++ ))
do
  # Format volume name with leading zeroes
  node_id=$(printf "%03d" $i)

  cat <<EOL >> docker-compose.yaml
  node${node_id}:
    name: node${node_id}
EOL
done

echo "docker-compose.yaml generated with $NODE_COUNT nodes."