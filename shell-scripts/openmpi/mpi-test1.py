#!/usr/bin/python

## OpenMPI Python Test: Send Recv

from __future__ import print_function
from mpi4py import MPI

## Debug
#import pprint

comm = MPI.COMM_WORLD
rank = comm.rank
size = comm.size
rcvdata = None

if rank == 0:
    data = {'a':1, 'b':2, 'c':3, 'd':'eFg'} # Python's dictionary, a kind of associative array
#   rcvdata = comm.ssend(data, dest=1, tag=77)

    for num in range(1, size):
        if num % 2 == 0:
            rcvdata = comm.ssend(data, dest=num, tag=77)
        else:
            rcvdata = comm.send({'odd':'data'}, dest=num, tag=77)

    print('I am Root #' + str(rank) + ' / Size=' + str(size) + ' and I broadcasted some data to childs')

else: # childs
    data = None
    rcvdata = comm.recv(data, source=0, tag=77)
    print('Child #' + str(rank) + ' get some Data from Root: ' + str(rcvdata));

comm.barrier();

## Debug
#pprint.pprint(globals())
#pprint.pprint(locals())