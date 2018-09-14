
/* Year 2038 C Test
Compile as: # cc ./2038.c
*/

#include <stdlib.h>
#include <stdio.h>
#include <unistd.h>
#include <time.h>

int main (int argc, char **argv) {
    time_t t;
    t = (time_t) 1000000000;
    printf ("%d, %s", (int) t, asctime (gmtime (&t)));
    t = (time_t) (0x7FFFFFFF);
    printf ("%d, %s", (int) t, asctime (gmtime (&t)));
    t++;
    printf ("%d, %s", (int) t, asctime (gmtime (&t)));
    return 0;
}
				
/*
END
*/
				