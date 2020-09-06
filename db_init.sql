CREATE TABLE request (
    id  SERIAL PRIMARY KEY,
    url TEXT  NOT NULL,
    interval INT NOT NULL,
    active BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE request_result (
   id SERIAL PRIMARY KEY,
   duration FLOAT,
   response TEXT,
   created_at TIMESTAMP DEFAULT NOW(),
   request_id INT REFERENCES request (id)
);

