# pip install neo4j-driver

from neo4j.v1 import GraphDatabase, basic_auth

driver = GraphDatabase.driver(
    "bolt://100.24.206.62:33829", 
    auth=basic_auth("neo4j", "nets-refunds-preference"))
session = driver.session()

cypher_query = '''
MATCH (m:Movie)
WHERE(m.movieId=$id)
RETURN m
LIMIT $limit
'''

results = session.run(cypher_query,
  parameters={"id": "2713", "limit": 10})

for record in results:
  print(record)
