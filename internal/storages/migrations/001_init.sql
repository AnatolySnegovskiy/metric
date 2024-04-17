-- Write your migrate up statements here
CREATE TABLE IF NOT EXISTS counter (name varchar(255) PRIMARY KEY, value int8);
CREATE TABLE IF NOT EXISTS gauge (name varchar(255) PRIMARY KEY, value DOUBLE PRECISION);
---- create above / drop below ----
DROP TABLE IF EXISTS counter;
DROP TABLE IF EXISTS gauge;
-- Write your migrate down statements here. If this migrations is irreversible
-- Then delete the separator line above.
