#!/usr/bin/python

## OpenMPI Python Test: Broadcast

from __future__ import print_function
from mpi4py import MPI

## Debug
#import pprint

comm = MPI.COMM_WORLD
rank = comm.rank
size = comm.size
rcvdata = None

if rank == 0:
    data = {'a':1, 'b':2, 'c':3, 'd':'eFg'} # ~ associative array
    rcvdata = comm.bcast(data, root=0)
    print('I am Root #' + str(rank) + ' / Size=' + str(size) + ' and I broadcasted this data to childs: ' + str(rcvdata))
else:
    data = None
    rcvdata = comm.bcast(data, root=0)
    print('Child #' + str(rank) + ' get some Data from Root: ' + str(rcvdata));

comm.barrier();

## Debug
#pprint.pprint(globals())
#pprint.pprint(locals())