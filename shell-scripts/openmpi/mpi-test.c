
/*

	# include <stdio.h>
	# include <stdlib.h>

main()
{

int res;
res = system("./mpi-test.sh");

printf("%s\n", res);

return res;

}

*/

/* program hello */
/* Adapted from mpihello.f by drs */

#include <mpi.h>
#include <stdio.h>
#include <unistd.h>
#include <string.h>

int main(int argc, char *argv[]) {

	int rank;
	int np;
	char hostname[256];
	char message[16384];
	FILE *fp;
	char output[1035];
	MPI_Status status;

	MPI_Init(&argc,&argv);
	MPI_Comm_rank(MPI_COMM_WORLD, &rank);
	MPI_Comm_size(MPI_COMM_WORLD, &np);

	gethostname(hostname, 255);


	//int res = system("./mpi-test.sh");

	if(rank == 0) {

		printf("Running MASTER process P#: %d on HOST %s\n", rank, hostname);

		MPI_Recv(&message,      // message buffer
			np-1,                // one data item
			MPI_CHAR,           // of type char real
			MPI_ANY_SOURCE,     // receive from any sender
			MPI_ANY_TAG,        // any type of message
			MPI_COMM_WORLD,     // default communicator
			&status
		);

		printf("%s", message);
		printf("##### DONE #####");

	} else {

		printf("Running CHILD process P#: %d on HOST %s\n", rank, hostname);

		fp = popen("/bin/sh ./mpi-test.sh", "r");
		if(fp == NULL) {
			printf("Failed to run command\n");
		} else {
			printf("Output from command:\n" );
			while(fgets(output, sizeof(output)-1, fp) != NULL) {
//				printf("%s", output);
				strcpy(message, output);
			}
			pclose(fp);
			MPI_Send(&message, sizeof(message), MPI_CHAR, 0, rank, MPI_COMM_WORLD);
		}

	} //end if else

	//MPI_Barrier();

	MPI_Finalize();

	return 0;
}