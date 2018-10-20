#include "terminal.h"
#include <gtk/gtk.h>
#include <stdint.h>
#include <stdlib.h>

// clear command

static struct { int16_t id; } windowid;

void setClear() { elem_clear(windowid.id); }

void commandClear() { readbuffcall(&windowid, sizeof windowid, setClear); }

// show command

void setShow() { window_redraw(windowid.id); }

void commandShow() { readbuffcall(&windowid, sizeof windowid, setShow); }

// fill command

static struct {
    int16_t id;
    int16_t x, y, width, height;
    int16_t r, g, b;
} fill;

void setFill() {
    elem_fill_add(fill.id, fill.x, fill.y, fill.width, fill.height, fill.r, fill.g, fill.b);
}

void commandFill() { readbuffcall(&fill, sizeof fill, setFill); }

// draw line

static struct {
    int16_t id;
    int16_t x0, y0;
    int16_t x1, y1;
    int16_t r, g, b;
} line;

void setLine() {
    elem_line_add(line.id, line.x0, line.y0, line.x1, line.y1, line.r, line.g, line.b);
}

void commandLine() { readbuffcall(&line, sizeof line, setLine); }

// draw string

static struct {
    int16_t id;
    int16_t x, y;
    int16_t r, g, b;
    int16_t fontid;
    int16_t fontsize;
} point;

void setText(void *text) {
    elem_text_add(point.id, point.x, point.y, text, point.fontid, point.r, point.g, point.b);
}

void commandText() { readalloccall(&point, sizeof point, setText); }

// load font

static struct {
    int16_t id;
    int16_t height;
    int16_t style, variant, weight, stretch;
} font;

void setFont(void *family) {
    font_elem_add(font.id, font.height, family, font.style, font.variant, font.weight,
                  font.stretch);
}

void commandFont() { readalloccall(&font, sizeof font, setFont); }

// split text

static struct {
    int16_t fontid;
    int16_t edge;
} split;

void splitText(void *text) {
    int16_t *out = font_split_text(split.fontid, text, split.edge);
    if (out == NULL) {
        int16_t value = 0;
        socket_input_write(&value, sizeof value);
    } else {
        int length = (1 + *out) * sizeof(int16_t);
        socket_input_write(out, length);
        free(out);
    }
}

void commandSplit() { readalloccall(&split, sizeof split, splitText); }

// app window size

static struct {
    int16_t x, y;
    int16_t width, height;
} size;

void setSize() {
    gtk_window_move(GTK_WINDOW(app), size.x, size.y);
    gtk_window_resize(GTK_WINDOW(app), size.width, size.height);
    gtk_widget_show_all(app);
}

void commandSize() { readbuffcall(&size, sizeof size, setSize); }

// app window title

void setTitle(void *buff) { gtk_window_set_title(GTK_WINDOW(app), buff); }

void commandTitle() { readalloccall(NULL, 0, setTitle); }

// window

static struct {
    int16_t id;
    int16_t x, y;
    int16_t width, height;
    int16_t r, g, b;
} window;

void setWindow() {
    window_create(window.id);
    window_size(window.id, window.x, window.y, window.width, window.height);
}

void commandWindow() { readbuffcall(&window, sizeof window, setWindow); }

// dispatch

void callcommand(char command) {
    switch (command) {
    // application
    case 'S':
        commandSize();
        break;
    case 'T':
        commandTitle();
        break;
    case 'X':
        gtk_main_quit();
        break;

    // window
    case 'D':
        commandWindow();
        break;

    // draw
    case 'W':
        commandShow();
        break;
    case 'C':
        commandClear();
        break;
    case 'F':
        commandFill();
        break;
    case 'L':
        commandLine();
        break;
    case 'U':
        commandText();
        break;

    // font
    case 'N':
        commandFont();
        break;
    case 'P':
        commandSplit();
        break;
    }
}
