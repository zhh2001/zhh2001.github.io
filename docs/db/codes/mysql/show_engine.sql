SHOW ENGINE InnoDB STATUS\G

-- 可同时查询当前配置，状态输出不是线程数量的唯一依据。
SHOW VARIABLES WHERE Variable_name IN
    ('innodb_read_io_threads', 'innodb_write_io_threads',
     'innodb_page_cleaners', 'innodb_purge_threads');
