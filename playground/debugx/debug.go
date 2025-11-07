package debugx

import (
	"debug/macho"
	"fmt"
	"os"
)

// FileType 表示文件类型
type FileType string

const (
	FileTypeELF       FileType = "ELF"
	FileTypeMachO     FileType = "Mach-O"
	FileTypeMachOFat  FileType = "Mach-O Fat (Universal)"
	FileTypePE        FileType = "PE/COFF"
	FileTypeWASM      FileType = "WebAssembly"
	FileTypeJavaClass FileType = "Java Class"
	FileTypeUnknown   FileType = "Unknown"
	FileTypeAOut      FileType = "a.out (legacy Unix)"
	FileTypeCOFF      FileType = "COFF Object File"
)

func detect(path string) (FileType, error) {
	f, err := os.Open(path)
	if err != nil {
		return FileTypeUnknown, err
	}
	defer f.Close()

	buf := make([]byte, 8)
	n, err := f.Read(buf)
	if err != nil {
		return FileTypeUnknown, err
	}
	if n < 4 {
		return FileTypeUnknown, nil
	}

	switch {
	// ELF
	case buf[0] == 0x7F && buf[1] == 'E' && buf[2] == 'L' && buf[3] == 'F':
		return FileTypeELF, nil

	// Mach-O 32/64-bit
	case (buf[0] == 0xFE && buf[1] == 0xED && buf[2] == 0xFA && (buf[3] == 0xCE || buf[3] == 0xCF)) ||
		(buf[0] == 0xCE && buf[1] == 0xFA && buf[2] == 0xED && (buf[3] == 0xFE || buf[3] == 0xCF)) ||
		(buf[0] == 0xCF && buf[1] == 0xFA && buf[2] == 0xED && buf[3] == 0xFE):
		return FileTypeMachO, nil

	// Mach-O Fat / Universal Binary
	case buf[0] == 0xCA && buf[1] == 0xFE && buf[2] == 0xBA && buf[3] == 0xBE:
		return FileTypeMachOFat, nil

	// PE/COFF (Windows)
	case buf[0] == 'M' && buf[1] == 'Z':
		return FileTypePE, nil

	// WebAssembly
	case buf[0] == 0x00 && buf[1] == 0x61 && buf[2] == 0x73 && buf[3] == 0x6D:
		return FileTypeWASM, nil

	// Java Class
	case buf[0] == 0xCA && buf[1] == 0xFE && buf[2] == 0xBA && buf[3] == 0xBE:
		return FileTypeJavaClass, nil

	// a.out (legacy Unix)
	case buf[0] == 0x01 || buf[0] == 0x02 || buf[0] == 0x03:
		return FileTypeAOut, nil

	// COFF Object file (simplified detection)
	case buf[0] == 0x4C && buf[1] == 0x01:
		return FileTypeCOFF, nil

	default:
		return FileTypeUnknown, nil
	}
}

func DebugInfo(file string) {
	f, err := macho.Open(file)
	if err != nil {
		fmt.Printf("Failed to open Mach-O file: %v\n", err)
		return
	}
	defer f.Close()

	fmt.Printf("Mach-O File: %s\n", file)
	fmt.Println("----------- Mach-O Header -----------")
	fmt.Printf("CPU:     %s\n", f.Cpu)
	fmt.Printf("Subtype: %d\n", f.SubCpu)
	fmt.Printf("Type:    %s\n", f.Type)
	fmt.Printf("Flags:   %#x\n", f.Flags)

	// 如果有 DWARF 调试信息，可以尝试读取
	if dw, err := f.DWARF(); err == nil {
		fmt.Println("DWARF debug info available")
		// 这里你可以进一步用 golang.org/x/debug/dwarf 解析
		_ = dw
	}
}
