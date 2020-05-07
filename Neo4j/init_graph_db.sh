#!/bin/bash

apt-get update
apt-get install cypher-shell -y

cypher-shell -a neo4j://graph_db:7687 -u neo4j -p parola_idp -f init_graph_db.cql
