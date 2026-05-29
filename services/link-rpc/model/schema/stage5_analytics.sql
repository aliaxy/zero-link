CREATE TABLE visit_event (
  id         BIGINT        NOT NULL AUTO_INCREMENT,
  link_id    BIGINT        NOT NULL DEFAULT 0,
  code       VARCHAR(32)   NOT NULL DEFAULT '',
  visited_at DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP,
  ip_hash    VARCHAR(64)   NOT NULL DEFAULT '',
  user_agent VARCHAR(512)  NOT NULL DEFAULT '',
  referer    VARCHAR(1024) NOT NULL DEFAULT '',
  device     VARCHAR(16)   NOT NULL DEFAULT 'unknown',
  PRIMARY KEY (id),
  KEY idx_visit_link_visited (link_id, visited_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE link_daily_stat (
  id        BIGINT NOT NULL AUTO_INCREMENT,
  link_id   BIGINT NOT NULL DEFAULT 0,
  stat_date DATE   NOT NULL,
  pv        BIGINT NOT NULL DEFAULT 0,
  uv        BIGINT NOT NULL DEFAULT 0,
  PRIMARY KEY (id),
  UNIQUE KEY uk_link_daily_stat (link_id, stat_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
