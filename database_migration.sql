CREATE DATABASE treebeer;

use treebeer;

CREATE TABLE event (
  id             bigint PRIMARY KEY AUTO_INCREMENT,
  created_at     date NOT NULL,
  raw_message    json NOT NULL
);
