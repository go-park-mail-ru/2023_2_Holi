time=30m

reg="hey -z $time -q 20 -m POST http://localhost:3001/api/v1/auth/register"
login="hey -z $time -q 100 -m POST http://localhost:3001/api/v1/auth/login"
logout="hey -z $time -q 15 -m POST http://localhost:3001/api/v1/auth/logout"
check="hey -z $time -q 30 -m POST http://localhost:3001/api/v1/auth/check"

osascript -e "tell app \"Terminal\" to do script \"$reg\""
osascript -e "tell app \"Terminal\" to do script \"$login\""
osascript -e "tell app \"Terminal\" to do script \"$logout\""
osascript -e "tell app \"Terminal\" to do script \"$check\""
