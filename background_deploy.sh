
# kill whatever is listening on port 6999
lsof -i :6999 | awk 'NR>1 {print $2}' | xargs kill -9

make all