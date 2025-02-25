SELECT 'CREATE DATABASE keycloak'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'keycloak')\gexec
SELECT 'CREATE USER keycloak PASSWORD ''keycloak'''
WHERE NOT EXISTS (SELECT FROM pg_user WHERE usename = 'keycloak')\gexec
GRANT ALL PRIVILEGES ON DATABASE keycloak TO keycloak;

SELECT 'CREATE DATABASE local'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'local')\gexec
SELECT 'CREATE USER local PASSWORD ''local'''
WHERE NOT EXISTS (SELECT FROM pg_user WHERE usename = 'local')\gexec
GRANT ALL PRIVILEGES ON DATABASE local TO local;

SELECT 'CREATE DATABASE latest'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'latest')\gexec
SELECT 'CREATE USER latest PASSWORD ''latest'''
WHERE NOT EXISTS (SELECT FROM pg_user WHERE usename = 'latest')\gexec
GRANT ALL PRIVILEGES ON DATABASE latest TO latest;

SELECT 'CREATE DATABASE preview'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'preview')\gexec
SELECT 'CREATE USER preview PASSWORD ''preview'''
WHERE NOT EXISTS (SELECT FROM pg_user WHERE usename = 'preview')\gexec
GRANT ALL PRIVILEGES ON DATABASE preview TO preview;

SELECT 'CREATE DATABASE testdb'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'testdb')\gexec
SELECT 'CREATE USER testdb PASSWORD ''testdb'''
WHERE NOT EXISTS (SELECT FROM pg_user WHERE usename = 'testdb')\gexec
GRANT ALL PRIVILEGES ON DATABASE testdb TO testdb;

SELECT 'CREATE DATABASE migrate'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'migrate')\gexec
SELECT 'CREATE USER migrate PASSWORD ''migrate'''
WHERE NOT EXISTS (SELECT FROM pg_user WHERE usename = 'migrate')\gexec
GRANT ALL PRIVILEGES ON DATABASE migrate TO migrate;

SELECT 'CREATE DATABASE acceptance'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'acceptance')\gexec
SELECT 'CREATE USER acceptance PASSWORD ''acceptance'''
WHERE NOT EXISTS (SELECT FROM pg_user WHERE usename = 'acceptance')\gexec
GRANT ALL PRIVILEGES ON DATABASE acceptance TO acceptance;

SELECT 'CREATE DATABASE photoprism_01'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'photoprism_01')\gexec
SELECT 'CREATE USER photoprism_01 PASSWORD ''photoprism_01'''
WHERE NOT EXISTS (SELECT FROM pg_user WHERE usename = 'photoprism_01')\gexec
GRANT ALL PRIVILEGES ON DATABASE photoprism_01 TO photoprism_01;

SELECT 'CREATE DATABASE photoprism_02'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'photoprism_02')\gexec
SELECT 'CREATE USER photoprism_02 PASSWORD ''photoprism_02'''
WHERE NOT EXISTS (SELECT FROM pg_user WHERE usename = 'photoprism_02')\gexec
GRANT ALL PRIVILEGES ON DATABASE photoprism_02 TO photoprism_02;

SELECT 'CREATE DATABASE photoprism_03'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'photoprism_03')\gexec
SELECT 'CREATE USER photoprism_03 PASSWORD ''photoprism_03'''
WHERE NOT EXISTS (SELECT FROM pg_user WHERE usename = 'photoprism_03')\gexec
GRANT ALL PRIVILEGES ON DATABASE photoprism_03 TO photoprism_03;

SELECT 'CREATE DATABASE photoprism_04'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'photoprism_04')\gexec
SELECT 'CREATE USER photoprism_04 PASSWORD ''photoprism_04'''
WHERE NOT EXISTS (SELECT FROM pg_user WHERE usename = 'photoprism_04')\gexec
GRANT ALL PRIVILEGES ON DATABASE photoprism_04 TO photoprism_04;

SELECT 'CREATE DATABASE photoprism_05'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'photoprism_05')\gexec
SELECT 'CREATE USER photoprism_05 PASSWORD ''photoprism_05'''
WHERE NOT EXISTS (SELECT FROM pg_user WHERE usename = 'photoprism_05')\gexec
GRANT ALL PRIVILEGES ON DATABASE photoprism_05 TO photoprism_05;

