SELECT * FROM `tb_user`;
SELECT * FROM `tb_user` WHERE `id` = 1;
SELECT * FROM `tb_user` WHERE `name` = 'zhang';

-- 查看每一条 SQL 的基本耗时
SHOW PROFILES;
-- +----------+------------+------------------------------------------------+
-- | Query_ID | Duration   | Query                                          |
-- +----------+------------+------------------------------------------------+
-- |        1 | 0.00017650 | SELECT DATABASE()                              |
-- |        2 | 0.00192125 | show databases                                 |
-- |        3 | 0.00244425 | show tables                                    |
-- |        4 | 0.00102350 | SELECT * FROM `tb_user`                        |
-- |        5 | 0.00057375 | SELECT * FROM `tb_user` WHERE `id` = 1         |
-- |        6 | 0.00031300 | SELECT * FROM `tb_user` WHERE `name` = 'zhang' |
-- +----------+------------+------------------------------------------------+

-- 查看指定 Query_ID 的 SQL 语句各个阶段的耗时情况
-- SHOW PROFILE FOR QUERY Query_ID;
SHOW PROFILE FOR QUERY 6;
-- +--------------------------------+----------+
-- | Status                         | Duration |
-- +--------------------------------+----------+
-- | starting                       | 0.000089 |
-- | Executing hook on transaction  | 0.000004 |
-- | starting                       | 0.000008 |
-- | checking permissions           | 0.000005 |
-- | Opening tables                 | 0.000038 |
-- | init                           | 0.000006 |
-- | System lock                    | 0.000009 |
-- | optimizing                     | 0.000010 |
-- | statistics                     | 0.000066 |
-- | preparing                      | 0.000013 |
-- | executing                      | 0.000012 |
-- | end                            | 0.000003 |
-- | query end                      | 0.000004 |
-- | waiting for handler commit     | 0.000008 |
-- | closing tables                 | 0.000007 |
-- | freeing items                  | 0.000022 |
-- | cleaning up                    | 0.000011 |
-- +--------------------------------+----------+

-- 查看指定 Query_ID 的 SQL 语句的 CPU 使用情况
-- SHOW PROFILE CPU FOR QUERY Query_ID;
SHOW PROFILE CPU FOR QUERY 6;
-- +--------------------------------+----------+----------+------------+
-- | Status                         | Duration | CPU_user | CPU_system |
-- +--------------------------------+----------+----------+------------+
-- | starting                       | 0.000089 | 0.000000 |   0.000080 |
-- | Executing hook on transaction  | 0.000004 | 0.000000 |   0.000004 |
-- | starting                       | 0.000008 | 0.000000 |   0.000008 |
-- | checking permissions           | 0.000005 | 0.000000 |   0.000005 |
-- | Opening tables                 | 0.000038 | 0.000000 |   0.000038 |
-- | init                           | 0.000006 | 0.000000 |   0.000006 |
-- | System lock                    | 0.000009 | 0.000000 |   0.000008 |
-- | optimizing                     | 0.000010 | 0.000000 |   0.000010 |
-- | statistics                     | 0.000066 | 0.000000 |   0.000066 |
-- | preparing                      | 0.000013 | 0.000000 |   0.000013 |
-- | executing                      | 0.000012 | 0.000000 |   0.000012 |
-- | end                            | 0.000003 | 0.000000 |   0.000003 |
-- | query end                      | 0.000004 | 0.000000 |   0.000003 |
-- | waiting for handler commit     | 0.000008 | 0.000000 |   0.000009 |
-- | closing tables                 | 0.000007 | 0.000000 |   0.000006 |
-- | freeing items                  | 0.000022 | 0.000000 |   0.000022 |
-- | cleaning up                    | 0.000011 | 0.000000 |   0.000011 |
-- +--------------------------------+----------+----------+------------+