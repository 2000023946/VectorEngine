#ifndef ARM64_H
#define ARM64_H

#include <stdint.h>

typedef struct {
    int ID;
    double Distance;
} Result;

void add_result(
    Result *results,
    int *resultCount,
    int k,
    int id,
    int64_t distance
);

int64_t vector_distance(int32_t *vector1, int32_t *vector2);
int64_t vector_distance_8bit(int8_t *vector1, int8_t *vector2);

void search_range(
    int32_t *query,
    int32_t *vectorData,
    int *ids,
    int start,
    int end,
    int k,
    Result *results
);

#endif