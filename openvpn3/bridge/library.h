#ifndef __LIBRARY_H__
#define __LIBRARY_H__

#ifdef __cplusplus
extern "C" {
#endif

#include <stdbool.h> //needed for bools :/


typedef int user_data;

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

typedef struct {
    user_data usrData;
    log_callback logCallback;
    stats_callback statsCallback;
    event_callback eventCallback;
} callbacks_delegate;


typedef struct {
    const char * username;
    const char * password;
} user_credentials;

//creates new session - nil on error,
void *new_session(const char * profile_content, user_credentials credentials, callbacks_delegate callbacks);

//starts created session
int start_session(void *ptr);
//stops running session
void stop_session(void *ptr);
//cleanups session
void cleanup_session(void *ptr);

void check_library(user_data userData, log_callback);

#ifdef __cplusplus
}
#endif

#endif