package vpk

import (
	"bufio"
	"encoding/binary"
	"io"
)

func treeReader(vpk VPK, reader *bufio.Reader, buffer []byte, cb func(entry *entry_impl)) error {
	var (
		ext, path, file []byte
		err             error
	)

	for {
		// Read file extension
		if ext, err = reader.ReadBytes('\x00'); err != nil {
			return err
		}
		if len(ext) == 1 {
			break
		}

		for {
			// Read file path
			if path, err = reader.ReadBytes('\x00'); err != nil {
				return err
			}
			if len(path) == 1 {
				break
			}

			for {
				// Read file name
				if file, err = reader.ReadBytes('\x00'); err != nil {
					return err
				}
				if len(file) == 1 {
					break
				}

				if _, err := io.ReadFull(reader, buffer[:18]); err != nil {
					return err
				}

				entry := &entry_impl{
					ext:  string(ext[:len(ext)-1]),
					path: string(path[:len(path)-1]),
					file: string(file[:len(file)-1]),

					crc:          binary.LittleEndian.Uint32(buffer[:4]),
					preloadBytes: binary.LittleEndian.Uint16(buffer[4:6]),
					archiveIndex: binary.LittleEndian.Uint16(buffer[6:8]),
					entryOffset:  binary.LittleEndian.Uint32(buffer[8:12]),
					entryLength:  binary.LittleEndian.Uint32(buffer[12:16]),
				}

				cb(entry)
			}
		}
	}

	return nil
}