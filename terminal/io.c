#include "terminal.h"
#include <gtk/gtk.h>
#include <stdint.h>
#include <stdlib.h>

static char *readChanBuffer = NULL;
static int readChanSize = 0;
static void (*readChanFunc)();

void readbuffcall(void *buffer, int size, void (*f)()) {
    readChanBuffer = buffer;
    readChanSize = size;
    readChanFunc = f;
}

static void (*readAllocFunc)(void *);
static int16_t length = 0;
static char *data = NULL;

static void alloccall() {
    data[length] = 0;
    readAllocFunc(data);
    free(data);
}

static void readdata() {
    data = malloc(length + 1);
    if (length == 0) {
        alloccall();
    } else {
        readbuffcall(data, length, alloccall);
    }
}

static void readsize() { readbuffcall(&length, sizeof length, readdata); }

void readalloccall(void *buffer, int size, void (*f)(void *)) {
    readAllocFunc = f;
    if (buffer != NULL) {
        readbuffcall(buffer, size, readsize);
    } else {
        readsize();
    }
}

gboolean readchan(GIOChannel *source, GIOCondition condition, gpointer data) {
    int ConnectFD = g_io_channel_unix_get_fd(source);
    if (readChanBuffer != NULL) {
        // waiting command
        int len = read(ConnectFD, readChanBuffer, readChanSize);
        if (len <= 0) {
            perror("read error");
            exit(EXIT_FAILURE);
            return TRUE;
        } else if (len < readChanSize) {
            readChanBuffer += len;
            readChanSize -= len;
            return TRUE;
        }
        // read parameters and call back func
        readChanBuffer = NULL;
        readChanSize = 0;
        (*readChanFunc)();
    } else {
        // single command
        char command;
        int len = read(ConnectFD, &command, 1);
        if (len <= 0) {
            perror("read error");
            exit(EXIT_FAILURE);
            return TRUE;
        }
        callcommand(command);
    }
    return TRUE;
}
