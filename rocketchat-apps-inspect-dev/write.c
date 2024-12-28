// go:build exclude
#include <fcntl.h>
#include <pthread.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/wait.h>
#include <unistd.h>

#define debug(...)                                                             \
  do {                                                                         \
    char buffer[BUFSIZ];                                                       \
    sprintf(buffer, "[%d] ", getpid());                                        \
    size_t n = write(STDOUT_FILENO, buffer, strlen(buffer));                   \
    memset(buffer, 0, n);                                                      \
    sprintf(buffer, __VA_ARGS__);                                              \
    write(STDOUT_FILENO, buffer, strlen(buffer));                              \
    write(STDOUT_FILENO, "\n", sizeof(char) * 2);                              \
  } while (0)

int main(int argc, char **argv) {
  if (argc > 1 && strncmp(argv[1], "--fork", strlen("--fork")) == 0) {
    pid_t pid = fork();
    if (pid == 0) {
      execl("./wait", (char *)NULL); // TODO: join pipes
    }

    if (pid == -1) {
      perror("failed");
    }
  }

  for (;;) {
    debug("noop");
    sleep(2);
  }

  return 0;
}
