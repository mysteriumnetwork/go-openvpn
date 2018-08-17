#ifndef __PROCESS_H__
#define __PROCESS_H__

#ifdef __cplusplus
extern "C" {
#endif

typedef struct {
    int bytes_in;
    int bytes_out;
} conn_stats;

typedef void(*stats_callback)(conn_stats);

int initProcess(char *argv[], int argc, stats_callback callback);

#ifdef __cplusplus
}
#endif

#endif