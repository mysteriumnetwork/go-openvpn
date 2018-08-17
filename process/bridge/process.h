#ifndef __PROCESS_H__
#define __PROCESS_H__

#ifdef __cplusplus
extern "C" {
#endif

#include <stdbool.h> //needed for bools :/


typedef void* user_data;

typedef void(*log_callback)(user_data, char * msg);

typedef struct {
    int bytes_in;
    int bytes_out;
} conn_stats;

typedef void(*stats_callback)(user_data, conn_stats);

//C++ -> C struct remap of openvpn::ClientApi::event
typedef struct {
    bool error;
    bool fatal;
    char *name;
    char *info;
} conn_event;

typedef void(*event_callback)(user_data, conn_event);

int initProcess(const char * profile_content , user_data userData , stats_callback statsCallback, log_callback, event_callback);

void checkLibrary(user_data userData, log_callback);

#ifdef __cplusplus
}
#endif

#endif