package plugin

import (
	"io/fs"
)

// TemplateManager is the interface for managing templates
type TemplateManager interface {
	// Filesystem retrieves a filesystem to access templates
	Filesystem() fs.FS
	
	// RenderToFile renders a template to a file
	RenderToFile(templatePath, targetPath string, data interface{}) error
	
	// RenderToString renders a template as a string
	RenderToString(templatePath string, data interface{}) (string, error)
}
