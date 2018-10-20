#ifndef _idlist_h_
#define _idlist_h_

typedef struct _id_list_elem id_list_elem;
typedef struct _id_list id_list;

id_list *id_list_new();
void id_list_free(id_list *list);

void id_list_append(id_list *list, int id, void *data);
void *id_list_get_data(id_list *list, int id);
void *id_list_remove(id_list *list, int id);
void *id_list_remove_any(id_list *list);

id_list_elem *id_list_root(id_list *list);
id_list_elem *id_list_elem_next(id_list_elem *elem);
void *id_list_elem_data(id_list_elem *elem);

#endif
