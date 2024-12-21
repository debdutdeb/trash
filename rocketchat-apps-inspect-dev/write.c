//go:build exclude
#include <string.h>
#include <unistd.h>
#include <stdio.h>

int main() {
    pid_t pid = getpid();
    printf("%d\n", pid);
    for (;;) {
        write(1, "noop more", strlen("noop more"));
        write(2, "stderr more", strlen("stderr more"));
        sleep(1);
    }
    return 0;
}
