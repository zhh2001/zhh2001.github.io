SELECT `OBJECT_SCHEMA`, `OBJECT_NAME`, `INDEX_NAME`,
       `LOCK_TYPE`, `LOCK_MODE`, `LOCK_DATA`
FROM `performance_schema`.`data_locks`;