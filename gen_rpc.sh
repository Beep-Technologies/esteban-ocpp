protoc -I ./protos -I=$GOPATH/src --gogofast_out=plugins=grpc:./api/rpc ./protos/*

echo 'DONE'
