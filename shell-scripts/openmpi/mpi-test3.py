#!/usr/bin/python

# OpenMPI Python Test 3

from mpi4py import MPI
import numpy

comm = MPI.COMM_WORLD
rank = comm.Get_rank()

# passing MPI datatypes explicitly
if rank == 0:
#    data = numpy.arange(10, dtype='i')
#    comm.Send([data, MPI.INT], dest=1, tag=77)
    data = {'a':1,'b':2,'c':3}
    comm.send(data, dest=1, tag=77)
elif rank == 1:
#    data = numpy.empty(10, dtype='i')
#    comm.Recv([data, MPI.INT], source=0, tag=77)
    data = None
    comm.recv(data, source=0, tag=77)
    print data

# automatic MPI datatype discovery
#if rank == 0:
#    data = numpy.arange(100, dtype=numpy.float64)
#    comm.Send(data, dest=1, tag=13)
#elif rank == 1:
#    data = numpy.empty(100, dtype=numpy.float64)
#    comm.Recv(data, source=0, tag=13)

