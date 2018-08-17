#include "process.h"
#include <stdlib.h>
#include <stdio.h>

int initProcess(char *argv[], int argc, stats_callback statsCallback)
{
    printf("Printing args:\n");
    for(int i=0; i< argc; i++ ) {
        printf("Arg: %d, val: %s\n",i,argv[i]);
    }
    printf("End\n");

    conn_stats stats;
    stats.bytes_in=123;
    stats.bytes_out=456;

    statsCallback(stats);

    return 0;
}
