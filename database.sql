CREATE TABLE IF NOT EXISTS users (
	id UUID PRIMARY KEY,
	full_name VARCHAR(50) NOT NULL,
  phone_number VARCHAR(16) NOT NULL,
  login_count INT NOT NULL DEFAULT 0,
  hashed_password VARCHAR(256) NOT NULL,
  created_at TIMESTAMP(0) DEFAULT NOW(),
  updated_at TIMESTAMP(0) DEFAULT NOW()
);
