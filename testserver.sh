mkfifo fifo0 
mkfifo fifo1
go run . > fifo0 < fifo1 &
python ${HOME}/gats/projects/scratch/nxcat-py/chattest.py < fifo0 > fifo1
kill $!