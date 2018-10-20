#include "terminal.h"
#include <arpa/inet.h>
#include <gtk/gtk.h>
#include <stdlib.h>
#include <string.h>

int socket_listen(int port) {
    int s = socket(PF_INET, SOCK_STREAM, IPPROTO_TCP);
    if (s == -1) {
        perror("cannot create socket");
        exit(EXIT_FAILURE);
    }

    struct sockaddr_in sa;
    memset(&sa, 0, sizeof sa);
    sa.sin_family = AF_INET;
    sa.sin_port = htons(port);
    sa.sin_addr.s_addr = htonl(INADDR_LOOPBACK);

    if (bind(s, (struct sockaddr *)&sa, sizeof sa) == -1) {
        perror("bind failed");
        close(s);
        exit(EXIT_FAILURE);
    }

    if (listen(s, 1) == -1) {
        perror("listen failed");
        close(s);
        exit(EXIT_FAILURE);
    }
    return s;
}

GIOChannel *socket_chan_open(int s) {
    if (s < 0) {
        perror("accept failed");
        close(s);
        exit(EXIT_FAILURE);
    }
    return g_io_channel_unix_new(s);
}

void socket_chan_close(GIOChannel *chan) {
    int s = g_io_channel_unix_get_fd(chan);

    char buff[1024];
    while (read(s, buff, sizeof buff) > 0)
        ;

    if (shutdown(s, SHUT_RDWR) == -1) {
        perror("shutdown failed");
    }
    if (close(s) == -1) {
        perror("close failed");
    }
}

static int socket_listen_input = 0;
static int socket_listen_output = 0;

int socket_input = 0;
int socket_output = 0;

static GIOChannel *chan = NULL;

void network_init() {
    socket_listen_input = socket_listen(1101);
    socket_listen_output = socket_listen(1102);
    socket_input = accept(socket_listen_input, NULL, NULL);
    chan = socket_chan_open(socket_input);
    g_io_add_watch(chan, G_IO_IN, readchan, NULL);
    socket_output = accept(socket_listen_output, NULL, NULL);
}

void network_done() {
    close(socket_output);
    socket_chan_close(chan);
    g_io_channel_unref(chan);
    close(socket_listen_output);
    close(socket_listen_input);
}

void socket_write(int s, void *data, int length) {
    if (s == 0) {
        perror("socket unknown");
        exit(EXIT_FAILURE);
    }
    size_t size = write(s, data, length);
    if (size != length) {
        perror("write failure");
    }
}

void socket_input_write(void *data, int length) { socket_write(socket_input, data, length); }

void socket_output_write(void *data, int length) { socket_write(socket_output, data, length); }
