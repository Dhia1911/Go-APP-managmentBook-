CREATE TABLE IF NOT EXISTS books (
  id SERIAL PRIMARY KEY,
  title VARCHAR(255),
  auther VARCHAR(255),
  publication_year INTEGER,
  created_at TIMESTAMP DEFAULT NOW()
);