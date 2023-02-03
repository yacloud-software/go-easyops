#!/bin/sh
TOKEN=`cat /home/cnw/.go-easyops/tokens/h2gproxy.token`
#make && test-ctx -server -ge_never_register_service_as_user 
make && test-ctx -server -ge_never_register_service_as_user -token=${TOKEN} -ge_debug_rpc_server
