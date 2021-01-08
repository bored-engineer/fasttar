package fasttar

import "time"

// BlockSize is the header size
const BlockSize = 512

// Name returns the Name of file entry
func Name(header []byte) (name []byte) {
	if len(header) == BlockSize {
		name = parseString(header[0:100])
	}
	return
}

// Mode returns the Permission and mode bits
func Mode(header []byte) (mode int64) {
	if len(header) == BlockSize {
		mode = parseNumeric(header[100:108])
	}
	return
}

// UID returns the User ID of owner
func UID(header []byte) (uid int64) {
	if len(header) == BlockSize {
		uid = parseNumeric(header[108:116])
	}
	return
}

// GID returns the Group ID of owner
func GID(header []byte) (gid int64) {
	if len(header) == BlockSize {
		gid = parseNumeric(header[116:124])
	}
	return
}

// Size returns the Logical file size in bytes
func Size(header []byte) (size int64) {
	if len(header) == BlockSize {
		size = parseNumeric(header[124:136])
	}
	return
}

// ModTime returns the Modification time
func ModTime(header []byte) (modTime time.Time) {
	if len(header) == BlockSize {
		modTime = time.Unix(parseNumeric(header[136:148]), 0)
	}
	return
}

// Chksum returns the checksum of the file
func Chksum(header []byte) (chksum int64) {
	if len(header) == BlockSize {
		chksum = parseNumeric(header[148:156])
	}
	return
}

// TypeFlag returns the TypeFlag
func TypeFlag(header []byte) (typeFlag byte) {
	if len(header) == BlockSize {
		typeFlag = header[156]
	}
	return
}

// LinkName returns the Target name of link (valid for TypeLink or TypeSymlink)
func LinkName(header []byte) (linkName []byte) {
	if len(header) == BlockSize {
		linkName = parseString(header[157:257])
	}
	return
}

// Magic returns the magic bytes
func Magic(header []byte) (magic []byte) {
	if len(header) == BlockSize {
		magic = header[257:263]
	}
	return
}

// Version returns the version bytes
func Version(header []byte) (version []byte) {
	if len(header) == BlockSize {
		version = header[263:265]
	}
	return
}

// UserName returns the username bytes
func UserName(header []byte) (username []byte) {
	if len(header) == BlockSize {
		username = parseString(header[265:297])
	}
	return
}

// GroupName returns the username bytes
func GroupName(header []byte) (groupname []byte) {
	if len(header) == BlockSize {
		groupname = parseString(header[297:329])
	}
	return
}

// DevMajor returns the Major device number (valid for TypeChar or TypeBlock)
func DevMajor(header []byte) (devMajor int64) {
	if len(header) == BlockSize {
		devMajor = parseNumeric(header[329:337])
	}
	return
}

// DevMinor returns the Minor device number (valid for TypeChar or TypeBlock)
func DevMinor(header []byte) (devMinor int64) {
	if len(header) == BlockSize {
		devMinor = parseNumeric(header[337:345])
	}
	return
}
