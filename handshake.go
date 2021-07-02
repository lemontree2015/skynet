package skynet

// Handshake使用的数据结构[TCP Level BSON编码的数据]

// [RPC Server -> RPC Client]
type ServiceRPCServerHandshake struct {
	ServiceName string // Service Name
	ClientUUID  string // RPC Server分配给RPC Client的UUID
}

// [RPC Client -> RPC Server]
type ServiceRPCClientHandshake struct {
	ClientUUID string
}
