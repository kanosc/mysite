mode=production
echo original parameters=[$*]
echo original OPTIND=[$OPTIND]
while getopts ":m:" opt
do
    case $opt in
        m)
            echo "this is -a option. OPTARG=[$OPTARG] OPTIND=[$OPTIND]"
            mode=$OPTARG
            ;;
        ?)
            echo "no valid parameter found."
            ;;
    esac
done

cd ./server
./compile.sh &&
cd ..
process_id=$(ps aux | grep static_server | grep -v 'sudo' | grep -v 'grep' | awk '{print $2 }')
echo current process_id is $process_id
sudo kill -9 $process_id
if [ $? -ne 0 ] 
then
   	echo "kill old server failed" 
else
	echo "kill old server success"
fi &&
nohup ./static_server.exe -mode ${mode} >> server_log.txt 2>&1 &
