#include "low.h"
#include "_cgo_export.h"

#include <glib.h>
#include <gdk-pixbuf/gdk-pixbuf.h>
#include <gtk/gtk.h>
#include <stdlib.h>
#include <string.h>

void on_destroy (GtkWidget *widget G_GNUC_UNUSED, gpointer user_data G_GNUC_UNUSED) {
    DestroyCallBack();
}

static gboolean s_keypress (GtkWidget *widget, GdkEventKey *event, gpointer data) {
    KeyboardCallBack (event);
    return TRUE;
}

GtkWidget *application_create() {
    gtk_init(NULL, NULL);

    GtkWidget *window;
    window = gtk_window_new (GTK_WINDOW_TOPLEVEL);
    return window;
}

void application_size(GtkWidget *app, int x, int y, int width, int height) {
    gtk_window_move (GTK_WINDOW (app), x, y);
    gtk_window_set_default_size (GTK_WINDOW (app), width, height);
}

void application_main(GtkWidget *app) {
    gtk_widget_show_all (GTK_WIDGET (app));
    g_signal_connect (app, "destroy", G_CALLBACK(on_destroy), NULL);
    g_signal_connect (app, "key_press_event", G_CALLBACK(s_keypress), NULL);
    gtk_main ();
}

void application_quit() {
    gtk_main_quit ();
}

GtkWidget *layout_create(GtkWidget *app) {
    GtkWidget *layout;
    layout = gtk_layout_new (NULL, NULL);
    gtk_container_add (GTK_CONTAINER (app), layout);
    return layout;
}

GtkWidget *window_create(GtkWidget *layout) {
    GtkWidget *window;
    window = gtk_image_new();
    gtk_container_add (GTK_CONTAINER (layout), window);
    return window;
}

void window_move(GtkWidget *layout, GtkWidget *window, int x, int y) {
    gtk_layout_move (GTK_LAYOUT (layout), window, x, y);
}

void window_set(GtkWidget *window, int width, int height, int stride, void *buffer) {
    GdkPixbuf *pixbuf;
    pixbuf = gdk_pixbuf_new_from_data((guint8 *)buffer, GDK_COLORSPACE_RGB, TRUE, 8,
        width, height, stride, NULL, NULL);
    gtk_image_set_from_pixbuf (GTK_IMAGE(window), pixbuf);
    g_object_unref(pixbuf);
}
