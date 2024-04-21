-- Write your migrate up statements here
CREATE TABLE counter (name varchar(255) PRIMARY KEY, value int8);
CREATE TABLE gauge (name varchar(255) PRIMARY KEY, value DOUBLE PRECISION);
---- create above / drop below ----
DROP TABLE counter;
DROP TABLE gauge;
-- Write your migrate down statements here. If this migrations is irreversible
-- Then delete the separator line above.
