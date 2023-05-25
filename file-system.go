package framework

type FileSystem struct {
	root       string
	controller string
	component  string
}

func (fs *FileSystem) Root() string {
	return fs.root
}

func (fs *FileSystem) Controller() string {
	return fs.root
}

func (fs *FileSystem) Component() string {
	return fs.root
}
