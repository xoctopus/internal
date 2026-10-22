// shared.c

#include "shared.h"

void fill_packet(struct Packet* p, uint32_t id, uint16_t flags, uint16_t length) {
    p->id = id;
    p->flags = flags;
    p->length = length;
}
