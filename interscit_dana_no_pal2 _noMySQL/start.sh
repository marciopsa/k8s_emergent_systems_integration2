#!/bin/bash

# turn on bash's job control
set -m

# start the primary process and put it in the background
./assembly.sh &
#./emergent.sh 

#sleep 30

# start the secundary process
#./learner.sh
