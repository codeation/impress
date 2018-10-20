#include "terminal.h"
#include <gtk/gtk.h>

GtkWidget *app = NULL;
GtkWidget *layout = NULL;

int main(int argc, char **argv) {
    gtk_init(NULL, NULL);

    app = gtk_window_new(GTK_WINDOW_TOPLEVEL);
    g_signal_connect(app, "destroy", G_CALLBACK(on_destroy), NULL);
    g_signal_connect(app, "key_press_event", G_CALLBACK(s_keypress), NULL);

    layout = gtk_layout_new(NULL, NULL);
    gtk_container_add(GTK_CONTAINER(app), layout);

    network_init();

    gtk_main();

    network_done();

    return 0;
}
