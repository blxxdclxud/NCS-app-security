# Buffer Overflow Project - File Explanation

This document explains each file in the buffer-overflow demonstration, what it does, and why it's needed for the project.

---

## 1. `vulnerable-program.c` - The Vulnerable Application

**What it does:**
- Contains a C program with an intentional stack buffer overflow vulnerability
- Uses the dangerous `gets()` function that doesn't check input bounds
- Has a 64-byte buffer that can be overflowed
- Includes a `secret_function()` that spawns a shell (the target of exploitation)
- Prints memory addresses to help with exploit development

**Why we need it:**
- **Core of the demo**: This is the vulnerable target we're exploiting
- **Educational purpose**: Shows real-world vulnerable code patterns (unsafe `gets()` usage)
- **Demonstration**: Proves that buffer overflows can redirect program execution
- **Address disclosure**: Prints addresses so we can craft exploits without ASLR complications

**Key components:**
- `vulnerable_function()`: Contains the vulnerable `gets(buffer)` call
- `secret_function()`: The function we want to jump to (spawns shell)
- Address printing: Helps determine where to jump in memory

---

## 2. `exploit.py` - The Exploitation Script

**What it does:**
- Python script that automates buffer overflow exploitation
- Generates payloads with the correct offset (72 bytes) and return address
- Can exploit locally (runs the program and sends payload) or remotely (via network)
- Dynamically extracts memory addresses from program output
- Creates payload files for manual testing
- Includes GDB debugging helper

**Why we need it:**
- **Automation**: Manually crafting payloads is tedious; this automates it
- **Demonstration**: Shows the complete exploitation process in action
- **Educational**: Shows how attackers build exploits (offset calculation, address extraction)
- **Flexibility**: Supports both local and remote exploitation scenarios
- **No dependencies**: Uses only Python stdlib (no pwntools required)

**Key functions:**
- `generate_exploit_manually()`: Creates payload file with padding + return address
- `exploit_local()`: Runs program, extracts addresses, sends payload, gets shell
- `exploit_remote()`: Connects to network service and sends payload
- `debug_with_gdb()`: Generates GDB command file for analysis

---

## 3. `compile.sh` - Build Script

**What it does:**
- Compiles the vulnerable C program with specific compiler flags
- Disables security protections for educational purposes:
  - `-fno-stack-protector`: Removes stack canaries
  - `-z execstack`: Makes stack executable (allows shellcode execution)
  - `-no-pie`: Disables Position Independent Executable (fixed addresses)
  - `-g -O0`: Includes debug symbols, no optimization
- Optionally checks binary protections with `checksec` if available
- Makes the binary executable

**Why we need it:**
- **Reproducibility**: Ensures consistent compilation settings
- **Educational clarity**: Disables protections so we can focus on the overflow mechanism
- **Documentation**: Shows exactly which protections are disabled and why
- **Convenience**: One command to build instead of remembering all flags

**Important note:** In production, you'd want ALL protections enabled. We disable them here to demonstrate the vulnerability clearly.

---

## 4. `Dockerfile` - Containerized Environment

**What it does:**
- Creates a Docker container with all necessary tools (gcc, gdb, python3, pwntools)
- Copies source files into the container
- Compiles the vulnerable program inside the container
- Disables ASLR (Address Space Layout Randomization) for deterministic addresses
- Sets up a consistent environment for testing

**Why we need it:**
- **Isolation**: Prevents affecting the host system
- **Reproducibility**: Same environment on any machine (no "works on my machine" issues)
- **Portability**: Easy to share and deploy
- **ASLR control**: Disables ASLR so addresses are predictable (needed for reliable exploits)
- **Clean environment**: Fresh Ubuntu with all dependencies pre-installed
- **Demo-ready**: Can be built and run immediately for presentations

**For the project:** This is essential for the demo - ensures everyone sees the same behavior.

---

## 5. `README.md` - Project Documentation

**What it does:**
- Explains the project structure and purpose
- Provides quick start instructions (Docker and local)
- Documents the exploitation flow step-by-step
- Shows how to verify the offset using GDB
- Lists mitigation strategies and safer alternatives
- Explains what to include in the report/demo

**Why we need it:**
- **Onboarding**: Helps anyone understand the project quickly
- **Instructions**: Step-by-step guide to build, run, and exploit
- **Reference**: Documents the technical details (offset calculation, protections)
- **Report guidance**: Tells you what screenshots/evidence to collect
- **Best practices**: Shows how to fix the vulnerability

**For the assignment:** This helps structure your report and ensures you cover all required sections.

---

## 6. `analysis/stack.md` - Stack Layout Analysis

**What it does:**
- Visualizes the stack layout before and after overflow
- Explains the 72-byte offset calculation (64-byte buffer + 8-byte saved RBP)
- Documents all security protections and why they're disabled
- Lists safer alternatives to vulnerable functions

**Why we need it:**
- **Understanding**: Visual representation helps grasp the memory layout
- **Offset explanation**: Shows why we need exactly 72 bytes of padding
- **Educational**: Explains each protection mechanism (ASLR, NX, canaries, PIE)
- **Mitigation reference**: Quick guide to fixing the vulnerability
- **Report material**: Perfect content for the "Analysis" section of your report

**For the demo:** Use these diagrams in your presentation to explain how the overflow works.

---

## Summary: Why Each File is Essential

| File | Purpose | Required? |
|------|---------|-----------|
| `vulnerable-program.c` | **Core target** - the vulnerable application | ✅ Essential |
| `exploit.py` | **Automation** - demonstrates exploitation | ✅ Essential |
| `compile.sh` | **Build** - ensures correct compilation | ✅ Essential |
| `Dockerfile` | **Environment** - reproducible container | ✅ Highly recommended |
| `README.md` | **Documentation** - instructions and reference | ✅ Essential |
| `analysis/stack.md` | **Analysis** - technical explanation | ✅ Recommended |

**For your assignment:**
- All files together create a complete, working demonstration
- The documentation files help you write the report
- The Docker setup makes the demo portable and reliable
- The exploit script proves the vulnerability works
- The analysis file provides technical depth for your report

---

## How They Work Together

1. **Build**: `compile.sh` compiles `vulnerable-program.c` → creates `vulnerable_program` binary
2. **Run**: Execute `vulnerable_program` → prints addresses
3. **Exploit**: `exploit.py` reads addresses, crafts payload, sends it → gets shell
4. **Containerize**: `Dockerfile` packages everything for easy deployment
5. **Document**: `README.md` and `analysis/stack.md` explain everything

This creates a complete, educational buffer overflow demonstration suitable for your Network and Cyber Security project assignment.

