package env

type Environment struct {
	ID           string `protobuf:"bytes,1,opt,name=id" json:"id" bson:"id"`
	Name         string `protobuf:"bytes,2,opt,name=name" json:"name" bson:"name"`
	Description  string `protobuf:"bytes,3,opt,name=description" json:"description,omitempty" bson:"description"`
	Timestamp    string `protobuf:"bytes,11,opt,name=timestamp" json:"timestamp,omitempty"`
	ModTimestamp string `protobuf:"bytes,15,opt,name=modTimestamp" json:"modTimestamp,omitempty" bson:"mod_timestamp"`
}

type Response struct {
	Code    int32  `protobuf:"varint,1,opt,name=code" json:"code,omitempty"`
	Message string `protobuf:"bytes,2,opt,name=message" json:"message,omitempty"`
}

type GetEnvironmentsResponse struct {
	Response     *Response      `protobuf:"bytes,1,opt,name=response" json:"-"`
	Environments []*Environment `protobuf:"bytes,2,rep,name=environments" json:"environments,omitempty"`
}

type CreateEnvironmentRequest struct {
	Environment *Environment `protobuf:"bytes,1,opt,name=environment" json:"environment,omitempty"`
}

type CreateEnvironmentResponse struct {
	Response *Response `protobuf:"bytes,1,opt,name=response" json:"-"`
	EnvId    string    `protobuf:"bytes,2,opt,name=envId" json:"envId,omitempty"`
}

type EnvironmentKey struct {
	// Tenant: The format is "{domain}/{project}"
	Tenant string `protobuf:"bytes,1,opt,name=tenant" json:"tenant,omitempty"`
	Name   string `protobuf:"bytes,2,opt,name=name" json:"name,omitempty" bson:"name"`
}

type GetEnvironmentRequest struct {
	EnvironmentId string `protobuf:"bytes,1,opt,name=environmentId" json:"environmentId,omitempty"`
}

type GetEnvironmentResponse struct {
	Response    *Response    `protobuf:"bytes,1,opt,name=response" json:"-"`
	Environment *Environment `protobuf:"bytes,2,opt,name=environment" json:"environment,omitempty"`
}

type DeleteEnvironmentRequest struct {
	EnvironmentId string `protobuf:"bytes,1,opt,name=environmentId" json:"environmentId,omitempty"`
}

type UpdateEnvironmentRequest struct {
	Environment *Environment `protobuf:"bytes,1,opt,name=environment" json:"environment,omitempty"`
}

type GetEnvironmentCountRequest struct {
	Domain  string `protobuf:"bytes,1,opt,name=domain" json:"domain,omitempty"`
	Project string `protobuf:"bytes,2,opt,name=project" json:"project,omitempty"`
}

type GetEnvironmentCountResponse struct {
	Response *Response `protobuf:"bytes,1,opt,name=response" json:"-"`
	Count    int64     `protobuf:"varint,2,opt,name=count" json:"count,omitempty"`
}
