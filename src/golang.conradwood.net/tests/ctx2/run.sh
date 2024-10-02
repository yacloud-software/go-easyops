#!/bin/sh
TOKEN=`cat /home/cnw/.go-easyops/tokens/h2gproxy.token`
make && ctx2 -ge_autokill_instance_on_port -ge_never_register_service_as_user -token=${TOKEN} -ge_debug_rpc_server -ge_debug_context -ge_grpc_print_errors
