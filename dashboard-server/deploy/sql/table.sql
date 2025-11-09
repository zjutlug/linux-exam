-- 记录用户提交的正确答案-- 用户表
CREATE TABLE user
(
    `id`           BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增ID',
    `container_id` VARCHAR(64)     NOT NULL COMMENT '容器id',
    `username`     VARCHAR(20)     NOT NULL COMMENT '用户名',
    `created_at`   TIMESTAMP(3)    NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
    `updated_at`   TIMESTAMP(3)    NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
    PRIMARY KEY (id),
    UNIQUE KEY `uk_username` (username)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;

-- 记录用户提交的正确答案
CREATE TABLE submission
(
    `id`           BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增ID',
    `container_id` VARCHAR(64)     NOT NULL COMMENT '容器ID',
    `question_id`  BIGINT          NOT NULL COMMENT '题目ID',
    `created_at`   TIMESTAMP(3)    NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
    `updated_at`   TIMESTAMP(3)    NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
    PRIMARY KEY (id),
    UNIQUE KEY uk_container_question (container_id, question_id)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;
