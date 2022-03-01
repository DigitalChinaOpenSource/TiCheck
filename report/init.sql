DROP TABLE IF EXISTS `tck_users`;
CREATE TABLE "tck_users" (
    "user_id" INTEGER NOT NULL,
    "user_name" TEXT NOT NULL unique,
    "user_password" TEXT NOT NULL,
    "full_name" TEXT NOT NULL,
    "email" TEXT NOT NULL,
    "is_enabled" INTEGER(1) NOT NULL, -- '0' is disable, '1' is enable
    "creator" TEXT NOT NULL,
    "create_time" INTEGER NOT NULL,
    PRIMARY KEY ("user_id")
);

insert into tck_users values (1, 'admin','21232f297a57a5a743894a0e4a801fc3','管理员','admin@ticheck.com', 1, 'system', 1645756378);

DROP TABLE IF EXISTS `tck_addons`;
CREATE TABLE "tck_addons" (
    "script_id" INTEGER NOT NULL,
    "script_name" TEXT NOT NULL  ,
    "script_file" TEXT NOT NULL unique,
    "script_tag" TEXT NOT NULL,
    "description" TEXT,
    "threshold_operator" INTEGER, --0-无，1-等于，2-大于，3-大于等于，4-小于，5-小于等于
    "threshold_value" TEXT,
    "threshold_args" TEXT,
    "is_system" INTEGER NOT NULL,
    "script_creator" TEXT NOT NULL, --系统脚本固定为system，自定义脚本为上传用户
    "script_create_time" INTEGER NOT NULL,
    "script_update_time" INTEGER,
    PRIMARY KEY ("script_id")
);

DROP TABLE IF EXISTS `tck_cluster`;
CREATE TABLE "tck_cluster" (
    "cluster_id" INTEGER NOT NULL,
    "cluster_name" TEXT NOT NULL,
    "prometheus_url" TEXT NOT NULL,
    "tidb_username" TEXT NOT NULL,
    "tidb_password" TEXT NOT NULL,
    "description" TEXT,
    "create_time" INTEGER NOT NULL,
    "cluster_owner" TEXT NOT NULL,
    "tidb_version" TEXT NOT NULL,
    "dashboard_url" TEXT,
    "grafana_url" TEXT,
    "last_check_time" INTEGER,
    "cluster_health" INTEGER,
    PRIMARY KEY ("cluster_id")
);

DROP TABLE IF EXISTS `tck_cluster_checklist`;
CREATE TABLE "tck_cluster_checklist" (
    "checklist_id" INTEGER NOT NULL,
    "cluster_id" INTEGER NOT NULL,
    "script_id" INTEGER NOT NULL,
    "is_enabled" INTEGER(1) NOT NULL, --0-未启用，1-已启用
    "threshold_operator" INTEGER,
    "threshold_value" TEXT,
    "threshold_args" TEXT,
    PRIMARY KEY ("checklist_id")
);

DROP TABLE IF EXISTS `tck_cluster_scheduler`;
CREATE TABLE "tck_cluster_scheduler" (
    "scheduler_id" INTEGER NOT NULL,
    "cluster_id" INTEGER NOT NULL,
    "title" TEXT NOT NULL,
    "cron_expression" TEXT NOT NULL,
    "is_active" INTEGER NOT NULL, -- 0-未激活，1-已激活
    "creator" TEXT NOT NULL,
    "create_time" INTEGER NOT NULL,
    "run_count" INTEGER NOT NULL,
    PRIMARY KEY ("scheduler_id")
);

DROP TABLE IF EXISTS `tck_cluster_check_history`;
CREATE TABLE "tck_cluster_check_history" (
    "check_time" INTEGER NOT NULL,
    "cluster_id" INTEGER NOT NULL,
    "scheduler_id" INTEGER,
    "normal_items" INTEGER NOT NULL,
    "warning_items" INTEGER NOT NULL,
    "total_items" INTEGER NOT NULL,
    "duration" INTEGER NOT NULL,
    PRIMARY KEY ("check_time")
);

DROP TABLE IF EXISTS `tck_cluster_check_data`;
CREATE TABLE "tck_cluster_check_data" (
    "id" INTEGER NOT NULL ,
    "check_time" INTEGER NOT NULL,
    "check_tag" TEXT NOT NULL,
    "check_name" TEXT NOT NULL,
    "operator" TEXT NOT NULL,
    "threshold" REAL,
    "duration" INTEGER NOT NULL,
    "check_item" TEXT NOT NULL,
    "check_value" REAL, -- 空代表脚本无输出
    "check_status" TEXT NOT NULL,
    PRIMARY KEY ("id")
);