#ifndef _LOW_H_
#define _LOW_H_

#include <gtk/gtk.h>

GtkWidget *application_create ();
void application_title (GtkWidget *app, char *title);
void application_size (GtkWidget *app, int x, int y, int width, int height);
void application_main (GtkWidget *app);
void application_quit ();

GtkWidget *layout_create (GtkWidget *app);

GtkWidget *window_create (GtkWidget *app);
void window_close (GtkWidget *window);
void window_move (GtkWidget *layout, GtkWidget *window, int x, int y);
void window_set (GtkWidget *window, int width, int height, int stride, void *buffer);

#endif
