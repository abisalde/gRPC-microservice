package catalog_pbuf

//go:generate protoc -I=.. --go_out=.. --go-grpc_out=.. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative --entgrpc_out=.. --entgrpc_opt=paths=source_relative,schema_path=../../schema catalog_pbuf/catalog_pbuf.proto catalog_pbuf/catalog_ext.proto
