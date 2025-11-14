int processCommand(client *c)
{
    // ...

    // 如果服务器设置了 server.maxmemory 属性
    if (server.maxmemory && !isInsideYieldingLongCommand())
    {
        // 尝试进行内存淘汰 performEvictions
        int out_of_memory = (performEvictions() == EVICT_FAIL);

        trackingHandlePendingKeyInvalidations();

        if (server.current_client == NULL)
            return C_ERR;

        if (out_of_memory && is_denyoom_command)
        {
            rejectCommand(c, shared.oomerr);
            return C_OK;
        }

        server.pre_command_oom_state = out_of_memory;
    }

    // ...
}