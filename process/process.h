#ifndef __PROCESS_H__
#define __PROCESS_H__

#ifdef __cplusplus
extern "C" {
#endif

typedef struct {
    int bytes_in;
    int bytes_out;
} ConnStats;

typedef void(*StatsCallback)(ConnStats);

int initProcess(char *argv[], int argc, StatsCallback callback);

#ifdef __cplusplus
}
#endif

#endif