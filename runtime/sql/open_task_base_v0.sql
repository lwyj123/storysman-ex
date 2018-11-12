-- 任务基础信息表
CREATE TABLE IF NOT EXISTS open_task_base
(
  `id`         BIGINT UNSIGNED                          NOT NULL COMMENT 'id',
  `task_id`    BIGINT UNSIGNED                          NOT NULL COMMENT '任务id',
  `creator_id` BIGINT UNSIGNED                          NOT NULL COMMENT '创建人user_id',
  `holder_id`  BIGINT UNSIGNED                         COMMENT '执行任务人user_id',
  `type`       TINYINT UNSIGNED                      NOT NULL COMMENT '任务分类(图片优化、标题优化等)',
  `status`     TINYINT UNSIGNED                      NOT NULL COMMENT '任务状态（待认领、已认领（待执行）、待审核、已通过（已执行）、已删除',
  `title`      VARCHAR(50)  COLLATE utf8_general_ci         NOT NULL COMMENT '任务标题',
  `digest`     VARCHAR(255) COLLATE utf8_general_ci         NOT NULL COMMENT '任务摘要',
  `target_type`  TINYINT        NOT NULL COMMENT '目标类型，专辑content或者章节item',
  `target_id`  BIGINT UNSIGNEd      COMMENt '目标id，如果是content就是content_id，item是item_id',
  `updated_at` TIMESTAMP NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '任务更新时间',
  `created_at` TIMESTAMP NOT NULL COMMENT '任务添加时间',
  PRIMARY KEY (id)
)
  ENGINE = InnoDB
  DEFAULT CHARSET = utf8
  COLLATE = utf8_general_ci
  COMMENT ='任务基础信息表';
