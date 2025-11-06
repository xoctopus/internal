// shared.c

#include "shared.h"

void fill_packet(struct Packet* p) {
    p->id = 42;
    p->flags = 0xFF;
    p->length = 1024;
}
