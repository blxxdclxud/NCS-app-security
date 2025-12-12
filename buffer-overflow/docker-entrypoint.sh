#!/bin/bash
# Disable ASLR for deterministic addresses (requires --privileged or --security-opt)
echo 0 > /proc/sys/kernel/randomize_va_space 2>/dev/null || echo "Warning: Could not disable ASLR (run with --privileged flag)"
exec "$@"

