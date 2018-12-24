package structs

import (
    "bufio"
    "fmt"
    "strings"
)

/* --- */

// Sequence database entry.
// Stores sequence in string form and data about precalculated dots.
type SequenceEntry struct {
    Name     string
    Sequence string
}

/* --- */

const NameSymbolsCut = 3

// Reads sequence entry from given reader.
func ReadFrom(reader *bufio.Reader, index int) (*SequenceEntry, bool) {
    for {
        firstChar, err := reader.Peek(1)

        if err != nil {
            return nil, false
        }

        if firstChar[0] == '>' {
            break
        }

        _, _, err = reader.ReadLine()

        if err != nil {
            return nil, false
        }
    }

    nameLine, _, err := reader.ReadLine()

    if err != nil {
        return nil, false
    }

    name := fmt.Sprintf("%s", string(nameLine[NameSymbolsCut:]))
    seq  := strings.Builder{}

    for {
        firstChar, err := reader.Peek(1)

        if err != nil || firstChar[0] == '>' {
            break
        }

        part, _, err := reader.ReadLine()

        if err != nil {
            break
        }

        seq.Write(part)
    }

    return &SequenceEntry {
        Name:     name,
        Sequence: seq.String(),
    }, true
}

// Writes sequence into given writer.
func (s *SequenceEntry) WriteInto(writer *bufio.Writer) {
    _, err := writer.WriteString(s.Name)

    if err != nil {
        panic(err)
    }

    err = writer.WriteByte('\n')

    if err != nil {
        panic(err)
    }

    _, err = writer.WriteString(s.Sequence)

    if err != nil {
        panic(err)
    }

    err = writer.WriteByte('\n')

    if err != nil {
        panic(err)
    }
}
