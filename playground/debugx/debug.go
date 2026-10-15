package debugx

import (
	"debug/dwarf"
	"debug/macho"
	"fmt"
	"io"
	"os"
	"strings"
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

// maxSubprogramsPerCU limits how many TagSubprogram entries are printed per CU.
const maxSubprogramsPerCU = 32

type dwarfCounts struct {
	cu, subprogram, variable, typ, other int
}

func detect(path string) (FileType, error) {
	f, err := os.Open(path)
	if err != nil {
		return FileTypeUnknown, err
	}
	defer func() { _ = f.Close() }()

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

// DebugInfo prints Mach-O and DWARF details to stderr (stdout stays clean for CLI piping).
func DebugInfo(file string) {
	WriteDebugInfo(os.Stderr, file)
}

// WriteDebugInfo is like DebugInfo but writes to w (for tests or custom sinks).
func WriteDebugInfo(w io.Writer, file string) {
	f, err := macho.Open(file)
	if err != nil {
		_, _ = fmt.Fprintf(w, "Failed to open Mach-O file: %v\n", err)
		return
	}
	defer func() { _ = f.Close() }()

	_, _ = fmt.Fprintf(w, "Mach-O File: %s\n", file)
	_, _ = fmt.Fprintln(w, "----------- Mach-O Header -----------")
	_, _ = fmt.Fprintf(w, "CPU:     %s\n", f.Cpu)
	_, _ = fmt.Fprintf(w, "Subtype: %d\n", f.SubCpu)
	_, _ = fmt.Fprintf(w, "Type:    %s\n", f.Type)
	_, _ = fmt.Fprintf(w, "Flags:   %#x\n", f.Flags)

	printDWARFSections(w, f)

	dw, err := f.DWARF()
	if err != nil {
		_, _ = fmt.Fprintf(w, "DWARF unavailable: %v\n", err)
		return
	}

	dumpDWARF(w, dw)
}

func printDWARFSections(w io.Writer, f *macho.File) {
	_, _ = fmt.Fprintln(w, "----------- DWARF Sections -----------")
	found := false
	for _, s := range f.Sections {
		if isDWARFSection(s.Name) {
			found = true
			_, _ = fmt.Fprintf(w, "  %-20s size=%d offset=%#x\n", s.Name, s.Size, s.Offset)
		}
	}
	if !found {
		_, _ = fmt.Fprintln(w, "  (none)")
	}
}

func isDWARFSection(name string) bool {
	switch {
	case strings.HasPrefix(name, "__debug"), strings.HasPrefix(name, "__zdebug"):
		return true
	case name == "__apple_names", name == "__apple_types",
		name == "__apple_namespac", name == "__apple_objc":
		return true
	default:
		return false
	}
}

func dumpDWARF(w io.Writer, dw *dwarf.Data) {
	_, _ = fmt.Fprintln(w, "----------- DWARF Info -----------")

	var total dwarfCounts
	r := dw.Reader()
	for {
		e, err := r.Next()
		if err != nil {
			_, _ = fmt.Fprintf(w, "DWARF read error: %v\n", err)
			return
		}
		if e == nil {
			break
		}

		switch e.Tag {
		case dwarf.TagCompileUnit, dwarf.TagPartialUnit:
			total.cu++
			dumpCompileUnit(w, dw, r, e, &total)
		default:
			countTag(&total, e.Tag)
			if e.Children {
				r.SkipChildren()
			}
		}
	}

	_, _ = fmt.Fprintln(w, "----------- DWARF Summary -----------")
	_, _ = fmt.Fprintf(w, "Compile units: %d\n", total.cu)
	_, _ = fmt.Fprintf(w, "Subprograms:   %d\n", total.subprogram)
	_, _ = fmt.Fprintf(w, "Variables:     %d\n", total.variable)
	_, _ = fmt.Fprintf(w, "Types:         %d\n", total.typ)
	_, _ = fmt.Fprintf(w, "Other:         %d\n", total.other)
}

func countTag(c *dwarfCounts, tag dwarf.Tag) {
	switch tag {
	case dwarf.TagSubprogram, dwarf.TagInlinedSubroutine:
		c.subprogram++
	case dwarf.TagVariable, dwarf.TagFormalParameter:
		c.variable++
	case dwarf.TagArrayType, dwarf.TagClassType, dwarf.TagEnumerationType,
		dwarf.TagPointerType, dwarf.TagStructType, dwarf.TagTypedef,
		dwarf.TagUnionType, dwarf.TagBaseType, dwarf.TagConstType,
		dwarf.TagVolatileType, dwarf.TagRestrictType, dwarf.TagSubroutineType,
		dwarf.TagUnspecifiedType, dwarf.TagRvalueReferenceType,
		dwarf.TagReferenceType:
		c.typ++
	default:
		c.other++
	}
}

func dumpCompileUnit(w io.Writer, dw *dwarf.Data, r *dwarf.Reader, cu *dwarf.Entry, total *dwarfCounts) {
	name, _ := cu.Val(dwarf.AttrName).(string)
	compDir, _ := cu.Val(dwarf.AttrCompDir).(string)
	producer, _ := cu.Val(dwarf.AttrProducer).(string)
	lang := cu.Val(dwarf.AttrLanguage)

	_, _ = fmt.Fprintf(w, "\n[CU] %s\n", name)
	if compDir != "" {
		_, _ = fmt.Fprintf(w, "  comp_dir:  %s\n", compDir)
	}
	if producer != "" {
		_, _ = fmt.Fprintf(w, "  producer:  %s\n", producer)
	}
	if lang != nil {
		_, _ = fmt.Fprintf(w, "  language:  %v\n", lang)
	}
	if ranges, err := dw.Ranges(cu); err == nil && len(ranges) > 0 {
		for _, rg := range ranges {
			_, _ = fmt.Fprintf(w, "  range:     [%#x, %#x)\n", rg[0], rg[1])
		}
	}

	if lr, err := dw.LineReader(cu); err == nil && lr != nil {
		files := lr.Files()
		n := 0
		for _, f := range files {
			if f != nil {
				n++
			}
		}
		_, _ = fmt.Fprintf(w, "  line files: %d\n", n)
		shown := 0
		for i, f := range files {
			if f == nil {
				continue
			}
			if shown >= 8 {
				_, _ = fmt.Fprintf(w, "    ... and %d more files\n", n-shown)
				break
			}
			_, _ = fmt.Fprintf(w, "    [%d] %s\n", i, f.Name)
			shown++
		}
	}

	if !cu.Children {
		return
	}

	printed := 0
	skipped := 0
	for {
		e, err := r.Next()
		if err != nil {
			_, _ = fmt.Fprintf(w, "  child read error: %v\n", err)
			return
		}
		if e == nil || e.Tag == 0 {
			break
		}

		countTag(total, e.Tag)

		if e.Tag == dwarf.TagSubprogram {
			if printed < maxSubprogramsPerCU {
				dumpSubprogram(w, dw, e)
				printed++
			} else {
				skipped++
			}
		}

		if e.Children {
			countChildren(r, total)
		}
	}
	if skipped > 0 {
		_, _ = fmt.Fprintf(w, "  ... and %d more subprograms\n", skipped)
	}
}

func countChildren(r *dwarf.Reader, total *dwarfCounts) {
	for {
		e, err := r.Next()
		if err != nil || e == nil || e.Tag == 0 {
			return
		}
		countTag(total, e.Tag)
		if e.Children {
			countChildren(r, total)
		}
	}
}

func dumpSubprogram(w io.Writer, dw *dwarf.Data, e *dwarf.Entry) {
	name, _ := e.Val(dwarf.AttrName).(string)
	if name == "" {
		name = "<anon>"
	}
	external, _ := e.Val(dwarf.AttrExternal).(bool)
	inline := e.Val(dwarf.AttrInline)

	_, _ = fmt.Fprintf(w, "  subprogram: %s", name)
	if external {
		_, _ = fmt.Fprint(w, " [external]")
	}
	if inline != nil {
		_, _ = fmt.Fprintf(w, " inline=%v", inline)
	}

	if ranges, err := dw.Ranges(e); err == nil && len(ranges) > 0 {
		for i, rg := range ranges {
			if i == 0 {
				_, _ = fmt.Fprintf(w, " pc=[%#x, %#x)", rg[0], rg[1])
			} else {
				_, _ = fmt.Fprintf(w, " +[%#x, %#x)", rg[0], rg[1])
			}
		}
	}
	_, _ = fmt.Fprintln(w)
}
