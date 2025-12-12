## Stack overflow visualization

### Stack layout before overflow (x86_64)
```
|-----------------------------|
|  saved RIP (return)         | <- target overwrite
|-----------------------------|
|  saved RBP                  |
|-----------------------------|
|  buffer[64]                 |
|-----------------------------|
|  ...                        |
```

### After >72 bytes of input
```
|-----------------------------|
|  new RIP -> secret_function | <- injected address
|-----------------------------|
|  overwritten RBP            |
|-----------------------------|
|  'A' * 64 + 8 bytes align   |
|-----------------------------|
```

Offset 72 bytes = 64-byte buffer + 8-byte saved RBP. Get `secret_function` address from program output or `nm ./vulnerable_program | grep secret_function`.

### Protections (and why disabled in PoC)
- **ASLR**: randomizes addresses; disabled for determinism.
- **NX/DEP**: non-exec stack; disabled via `-z execstack`.
- **Stack canaries**: detect overflow; disabled with `-fno-stack-protector`.
- **PIE/RELRO**: randomize/protect GOT/PLT; disabled via `-no-pie` and missing RELRO flags.

### Safer alternatives
- Input: `fgets(buffer, sizeof(buffer), stdin)` or `getline`.
- Copying: `strncpy`, `strlcpy`, `snprintf`.
- Plus: enforce input size limits, validation, compile with canaries/RELRO/PIE/NX, keep ASLR on.

