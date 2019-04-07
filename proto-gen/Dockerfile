FROM znly/protoc:0.4.0

COPY ./proto /proto
ENTRYPOINT ["sh", "-c", "protoc --proto_path=${BUILD_PATH} --go_out=plugins=grpc:${OUTPUT_PATH} ${FILENAME}"]
