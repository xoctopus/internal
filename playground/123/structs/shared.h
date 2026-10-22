// shared.h

#ifndef SHARED_H
#define SHARED_H

#include <stdint.h>

struct Packet {
    uint32_t id;
    uint16_t flags;
    uint16_t length;
};

// C function to modify struct field
void fill_packet(struct Packet* p, uint32_t id, uint16_t flags, uint16_t length);

#endif
