DROP CONSTRAINT ON (sc:ServiceCategory) ASSERT sc.id is UNIQUE;
DROP CONSTRAINT ON (sc:SkillCategory) ASSERT sc.id is UNIQUE;

MATCH (sc:ServiceCategory)<-[r:IsSubclassOf]-(sk:SkillCategory) DELETE r;
MATCH (sc:ServiceCategory) DELETE sc;
MATCH (sk:SkillCategory) DELETE sk;

CREATE CONSTRAINT ON (sc:ServiceCategory) ASSERT sc.id is UNIQUE;
CREATE CONSTRAINT ON (sc:SkillCategory) ASSERT sc.id is UNIQUE;

LOAD CSV WITH HEADERS FROM "https://idp-graphql-init-scripts.s3.eu-central-1.amazonaws.com/serviceCategories.csv"
AS row
CREATE (sc:ServiceCategory {id: toInteger(row.id), category: row.category});

LOAD CSV WITH HEADERS FROM "https://idp-graphql-init-scripts.s3.eu-central-1.amazonaws.com/skillCategories.csv"
AS row
CREATE (sc:SkillCategory {id: toInteger(row.id), category: row.skillCategory});

LOAD CSV WITH HEADERS FROM "https://idp-graphql-init-scripts.s3.eu-central-1.amazonaws.com/skillCategories.csv"
AS row
MATCH (a:ServiceCategory), (b:SkillCategory)
WHERE (a.id = toInteger(row.serviceCategory)) AND (b.id = toInteger(row.id))
MERGE (a)<-[r:IsSubclassOf]-(b);

 
