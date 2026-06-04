CREATE TABLE short_link_archive (
  id          BIGINT        NOT NULL,
  code        VARCHAR(32)   NOT NULL,
  origin_url  TEXT          NOT NULL,
  title       VARCHAR(255)  NOT NULL DEFAULT '',
  description VARCHAR(1024) NOT NULL DEFAULT '',
  status      TINYINT       NOT NULL,
  expire_at   DATETIME      NULL,
  created_by  BIGINT        NOT NULL,
  created_at  DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at  DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at  DATETIME      NULL,
  PRIMARY KEY (id),
  KEY idx_sla_code (code),
  KEY idx_sla_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE reserved_code (
  code        VARCHAR(32) NOT NULL,
  reserved_at DATETIME    NOT NULL,
  PRIMARY KEY (code)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
