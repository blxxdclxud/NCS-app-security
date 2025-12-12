echo "[*] Building vulnerable program..."

# Disable protections: -fno-stack-protector (no canaries), -z execstack (executable stack), -no-pie (fixed addresses)
gcc \
    -fno-stack-protector \
    -z execstack \
    -no-pie \
    -g \
    -O0 \
    -m64 \
    vulnerable-program.c -o vulnerable_program

if command -v checksec >/dev/null 2>&1; then
  echo "[*] Checking binary protections:"
  checksec --file=vulnerable_program
else
  echo "[*] checksec not found, skipping protection check"
fi

chmod +x vulnerable_program

echo "[*] Build finished!"
echo "[*] Run: ./vulnerable_program"
