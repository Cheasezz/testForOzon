CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS posts (
  id UUID PRIMARY KEY,
  user_id VARCHAR(255) NOT NULL,
  created_at TIMESTAMPTZ NOT NULL,
  title TEXT NOT NULL,
  content TEXT NOT NULL,
  comments_allowed BOOLEAN NOT NULL
);

CREATE INDEX idx_posts_created_at ON posts (created_at ASC);

CREATE TABLE IF NOT EXISTS posts_comments (
  post_id UUID NOT NULL,
  id UUID PRIMARY KEY,
  parent_id UUID NULL,
  user_id VARCHAR(255) NOT NULL,
  created_at TIMESTAMPTZ NOT NULL,
  content VARCHAR(2000) NOT NULL,
  FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
  FOREIGN KEY (parent_id) REFERENCES posts_comments(id) ON DELETE CASCADE
);

CREATE INDEX idx_comments_post_created_at ON posts_comments (post_id, created_at DESC);