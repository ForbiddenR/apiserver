package server

type APIGroupInfo struct {
}

type GenericAPIServer struct {
	Handler *APIServerHandler
}

func (s *GenericAPIServer) InstallAPIGroups(apiGroupInfos ...*APIGroupInfo) error {
	for range apiGroupInfos {
		s.Handler.GoRestfulApp.Group("")
	}
	return nil
}

func (s *GenericAPIServer) InstallAPIGroup(apiGroupInfo *APIGroupInfo) error {
	return s.InstallAPIGroups(apiGroupInfo)
}
