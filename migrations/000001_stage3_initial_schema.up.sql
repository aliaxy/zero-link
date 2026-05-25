CREATE TABLE admin_user (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  username VARCHAR(64) NOT NULL,
  password_hash VARCHAR(255) NOT NULL,
  status TINYINT NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY uk_admin_user_username (username),
  KEY idx_admin_user_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE short_link (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  code VARCHAR(32) NOT NULL,
  origin_url TEXT NOT NULL,
  title VARCHAR(255) NOT NULL DEFAULT '',
  description VARCHAR(1024) NOT NULL DEFAULT '',
  status TINYINT NOT NULL,
  expire_at DATETIME NULL,
  created_by BIGINT NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at DATETIME NULL,
  UNIQUE KEY uk_short_link_code (code),
  KEY idx_short_link_status (status),
  KEY idx_short_link_created_at (created_at),
  KEY idx_short_link_created_by (created_by),
  KEY idx_short_link_expire_at (expire_at),
  CONSTRAINT fk_short_link_created_by
    FOREIGN KEY (created_by) REFERENCES admin_user (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

INSERT INTO admin_user (username, password_hash, status)
VALUES (
  'admin',
  '$2a$10$UpD2JjqWVgQOatvqxd5H3OSQwzxC5o5gYf31R73AJIz.dQOAuKkBS',
  1
);
