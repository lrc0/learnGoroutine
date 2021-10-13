##
先构建
nohup ./clockwallserver 8001 Local &
nohup ./clockwallserver 8002 US/Eastern &
nohup ./clockwallserver 8003 Europe/London &

##
cd clockwallclient
go run main.go Local US/Eastern Europe/London