## Buffer Overflow Demo

Classic stack-based buffer overflow demo: overwrite return address and pop an interactive `/bin/sh`. 

### Layout
- `vulnerable-program.c` — vulnerable binary (`gets` + stack buffer).
- `compile.sh` — build with protections disabled (RELRO, PIE, canaries, NX, ASLR).
- `exploit.py` — payload generation/send (stdlib only).
- `Dockerfile` — container with build + optional pwntools.
- `analysis/stack.md` — stack layout and mitigations overview.

### Quick start (Docker Compose)

**Option 1: One command (auto-run exploit)**
```bash
cd /home/yusuf/Documents/NCS-app-security
docker-compose up buffer-overflow
```
This automatically builds (if needed), runs the container, and executes the exploit.

**Option 2: Interactive mode**
```bash
cd /home/yusuf/Documents/NCS-app-security
docker-compose run buffer-overflow /bin/bash
# Inside container, run:
python3 exploit.py local
```

**Option 3: Build only**
```bash
docker-compose build buffer-overflow
```

**Note:** The `privileged: true` setting in docker-compose.yml is required to disable ASLR for deterministic addresses. The entrypoint script automatically sets `/proc/sys/kernel/randomize_va_space = 0` at runtime.

### Local build and run
```bash
cd buffer-overflow
./compile.sh          # requires gcc (checksec optional)
./vulnerable_program  # prints buffer and secret_function addresses
```

### Exploitation flow
1) Grab addresses from program output (buffer + `secret_function`).  
2) Confirm offset to saved RIP (x64: 72 bytes) via cyclic pattern (`cyclic`, `cyclic_find` in GDB or pwntools).  
3) Build payload: `offset` bytes padding + `secret_function` address (little-endian).  
4) Send payload in one line (`exploit.py exploit_local()`), drop to shell.

Handy offset check:
```bash
gdb -q ./vulnerable_program
(gdb) run <<< "$(python3 - <<'PY';from pwn import cyclic;print(cyclic(200).decode());PY)"
(gdb) i r rsp rip
# feed RIP into cyclic_find()
```

### Safer build / mitigations
- Replace unsafe calls: `gets` → `fgets` / `getline` with length checks; `strcpy/strcat` → `strncpy`/`strncat` or `snprintf`.
- Compile with hardening:
```bash
gcc -D_FORTIFY_SOURCE=2 -fstack-protector-strong -O2 \
    -pie -fPIE -Wl,-z,relro,-z,now -Wl,-z,noexecstack safe.c -o safe_program
```
- Keep ASLR on: `echo 2 | sudo tee /proc/sys/kernel/randomize_va_space`.
- Sandbox: drop capabilities, enable seccomp/AppArmor, log abnormal input.

### What to show in report / demo
- Screenshot of addresses + payload send → shell.
- Stack before/after overwrite (`analysis/stack.md`).
- Compare vulnerable vs hardened build: `checksec ./vulnerable_program` vs `checksec ./safe_program`.

