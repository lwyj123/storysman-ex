-- 任务历史修改表 audit table,
CREATE TABLE IF NOT EXISTS open_task_base
(
  `id`         BIGINT UNSIGNED                          NOT NULL COMMENT 'id',
  `task_id`    BIGINT UNSIGNED                          NOT NULL COMMENT '任务id',
  `old_value`  TEXT                   NOT NULL COMMENT '老值(json)',
  `json_diff`  TEXT                   NOT NULL COMMENT 'jsondiffpatch(json)',
  `new_value`  TEXT                   NOT NULL COMMENT '新值(json)',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP   NOT NULL COLLATE '任务记录添加时间',
  PRIMARY KEY (id)
)
  ENGINE = InnoDB
  DEFAULT CHARSET = utf8
  COLLATE = utf8_general_ci
  COMMENT ='任务历史审计表';
