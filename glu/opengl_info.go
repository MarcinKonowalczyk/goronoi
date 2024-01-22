package glu

import "github.com/go-gl/gl/v3.3-core/gl"

type OpenGLInfo struct {
	Vendor        string
	Renderer      string
	VersionString string
	VersionMajor  int
	VersionMinor  int
	ShaderVersion string
	Extensions    []string
}

func GetOpenGLInfo() OpenGLInfo {
	vendor := gl.GoStr(gl.GetString(gl.VENDOR))
	renderer := gl.GoStr(gl.GetString(gl.RENDERER))
	version := gl.GoStr(gl.GetString(gl.VERSION))
	shaderVersion := gl.GoStr(gl.GetString(gl.SHADING_LANGUAGE_VERSION))

	var major int32
	gl.GetIntegerv(gl.MAJOR_VERSION, &major)

	var minor int32
	gl.GetIntegerv(gl.MINOR_VERSION, &minor)

	var numExtensions int32
	gl.GetIntegerv(gl.NUM_EXTENSIONS, &numExtensions)

	extensions := make([]string, numExtensions)
	for i := int32(0); i < numExtensions; i++ {
		extensions[i] = gl.GoStr(gl.GetStringi(gl.EXTENSIONS, uint32(i)))
	}

	return OpenGLInfo{
		Vendor:        vendor,
		Renderer:      renderer,
		VersionString: version,
		VersionMajor:  int(major),
		VersionMinor:  int(minor),
		ShaderVersion: shaderVersion,
		Extensions:    extensions,
	}
}

func (info OpenGLInfo) SupportsExtension(extension string) bool {
	for _, ext := range info.Extensions {
		if ext == extension {
			return true
		}
	}
	return false
}

func (info OpenGLInfo) SupportsOpenGLVersion(major, minor int) bool {
	return info.VersionMajor > major || (info.VersionMajor == major && info.VersionMinor >= minor)
}

// Make the OpenGLInfo nicely printable
func (info OpenGLInfo) String() string {
	out := "Vendor: " + info.Vendor + "\n" +
		"Renderer: " + info.Renderer + "\n" +
		"Version: " + info.VersionString + "\n" +
		"Shader Version: " + info.ShaderVersion

	if len(info.Extensions) > 0 {
		out += "\nExtensions:\n"
		for i, ext := range info.Extensions {
			out += "  " + ext
			if i < len(info.Extensions)-1 {
				out += "\n"
			}
		}
	}
	return out
}
