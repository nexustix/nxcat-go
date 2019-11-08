program1="go run ."
program2="python ${HOME}/gats/projects/scratch/nxcat-py/chattest.py"

mkfifo fifo0 
mkfifo fifo1
$program1 > fifo0 < fifo1 &
$program2 < fifo0 > fifo1
kill $!
rm fifo0
rm fifo1