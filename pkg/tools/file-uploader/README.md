# file-uploader

testing tool for uploading files to process-storage service

### usage

for guaranteed compatibility with process-storage grpc api, first generate grpc server and client
in ```/pkg/process-storage/grpc/generate.sh``` and then
build with ```go build``` and run ```./file-uploader.exe upload -processId=<<processId>> -processFile=<<pathToFile>>```