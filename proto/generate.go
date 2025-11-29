package proto

//go:generate protoc --go_out=../pkg/api --go_opt=paths=source_relative usp-record-1-4.proto
//go:generate protoc --go_out=../pkg/api --go_opt=paths=source_relative usp-msg-1-4.proto
