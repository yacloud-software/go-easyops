#!/bin/sh
export GRPC_GO_LOG_VERBOSITY_LEVEL=99
export GRPC_GO_LOG_SEVERITY_LEVEL=info

authtest -go_easyops_fancy_balancer -role=client
