from p4utils.utils.sswitch_p4runtime_API import SimpleSwitchP4RuntimeAPI

controller = SimpleSwitchP4RuntimeAPI(device_id=1, grpc_port=9559,
                                      p4rt_path='l2_forwarding_p4rt.txt',
                                      json_path='l2_forwarding.json')