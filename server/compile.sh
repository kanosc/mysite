ROOT=$(cd ..; pwd)
echo $ROOT
OUTPUT_PATH=${ROOT}
go build -o ${OUTPUT_PATH}/static_server.exe

