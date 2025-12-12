#include <stdio.h>
#include <string.h>
#include <stdlib.h>

// gets was removed from modern standards; declare manually for demo purposes
char *gets(char *);

// Disable stdio buffering for cleaner exploit I/O
void setup() {
    setvbuf(stdin, NULL, _IONBF, 0);
    setvbuf(stdout, NULL, _IONBF, 0);
    setvbuf(stderr, NULL, _IONBF, 0);
}

void vulnerable_function() {
    char buffer[64];
    
    printf("Address buffer: %p\n", buffer);
    printf("Address vulnerable_function: %p\n", vulnerable_function);
    printf("Enter data: ");
    
    // Dangerous: no bounds check
    gets(buffer);
    
    printf("You entered: %s\n", buffer);
}

void secret_function() {
    printf("\n[+] Congrats! You jumped to secret_function!\n");
    printf("[+] Spawning shell...\n");
    system("/bin/sh");
}

int main(int argc, char *argv[]) {
    setup();
    
    printf("====================================\n");
    printf("     Buffer Overflow Demonstration    \n");
    printf("====================================\n\n");
    
    printf("Address main: %p\n", main);
    printf("Address secret_function: %p\n", secret_function);
    
    vulnerable_function();
    
    printf("\nProgram finished.\n");
    return 0;
}
