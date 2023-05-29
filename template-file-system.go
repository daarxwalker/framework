package framework

type templateFileSystem struct{}

func (tfs *templateFileSystem) Global(path string) string {
	return tfs.combinePath(templateSourceGlobal, path)
}

func (tfs *templateFileSystem) Module(path string) string {
	return tfs.combinePath(templateSourceModule, path)
}

func (tfs *templateFileSystem) Controller(path string) string {
	return tfs.combinePath(templateSourceController, path)
}

func (tfs *templateFileSystem) Component(path string) string {
	return tfs.combinePath(templateSourceComponent, path)
}

func (tfs *templateFileSystem) combinePath(sourceType, path string) string {
	return sourceType + ":" + path
}
