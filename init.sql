CREATE TABLE users (
  id UUID PRIMARY KEY,
  username VARCHAR(255) NOT NULL UNIQUE
);

CREATE TABLE messages (
  message_id SERIAL PRIMARY KEY,
  user_id UUID NOT NULL,
  message TEXT NOT NULL,
  FOREIGN KEY (user_id) REFERENCES users(id)
);