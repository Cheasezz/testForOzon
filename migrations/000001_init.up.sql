CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS posts (
  id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
  created_at TIMESTAMP DEFAULT now() NOT NULL,
  title TEXT NOT NULL,
  content TEXT NOT NULL,
  comments_allowed BOOLEAN DEFAULT TRUE NOT NULL
);

CREATE TABLE IF NOT EXISTS posts_comments (
  post_id UUID REFERENCES posts (id) ON DELETE CASCADE,
  id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
  parent_id UUID NULL REFERENCES posts_comments (id) ON DELETE CASCADE
  created_at TIMESTAMP DEFAULT now() NOT NULL,
  content VARCHAR(2000) NOT NULL,
);