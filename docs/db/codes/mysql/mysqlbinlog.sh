cd /var/lib/mysql
mysqlbinlog binlog.000001
# # The proper term is pseudo_replica_mode, but we use this compatibility alias
# # to make the statement usable on server versions 8.0.24 and older.
# /*!50530 SET @@SESSION.PSEUDO_SLAVE_MODE=1*/;
# /*!50003 SET @OLD_COMPLETION_TYPE=@@COMPLETION_TYPE,COMPLETION_TYPE=0*/;
# DELIMITER /*!*/;
# # at 4
# #251021 15:33:17 server id 1  end_log_pos 126 CRC32 0x14a26c56  Start: binlog v 4, server v 8.0.43-0ubuntu0.24.04.2 created 251021 15:33:17 at startup
# ROLLBACK/*!*/;
# BINLOG '
# PTf3aA8BAAAAegAAAH4AAAAAAAQAOC4wLjQzLTB1YnVudHUwLjI0LjA0LjIAAAAAAAAAAAAAAAAA
# AAAAAAAAAAAAAAAAAAA9N/doEwANAAgAAAAABAAEAAAAYgAEGggAAAAICAgCAAAACgoKKioAEjQA
# CigAAVZsohQ=
# '/*!*/;
# # at 126
# #251021 15:33:17 server id 1  end_log_pos 157 CRC32 0xfd4f9dcd  Previous-GTIDs
# # [empty]
# # at 157
# #251021 15:33:18 server id 1  end_log_pos 180 CRC32 0x95db574e  Stop
# SET @@SESSION.GTID_NEXT= 'AUTOMATIC' /* added by mysqlbinlog */ /*!*/;
# DELIMITER ;
# # End of log file
# /*!50003 SET COMPLETION_TYPE=@OLD_COMPLETION_TYPE*/;
# /*!50530 SET @@SESSION.PSEUDO_SLAVE_MODE=0*/;