package transport

import (
    "bufio"
    "bytes"
    "encoding/binary"
    "fmt"
    "io"
)

func ReadFrame(r *bufio.Reader) (byte, []byte, error) {
    h1, err := r.ReadByte()
    if err != nil {
        return 0, nil, err
    }
    h2, err := r.ReadByte()
    if err != nil {
        return 0, nil, err
    }

    opcode := h1 & 0x0F
    masked := (h2 & 0x80) != 0
    length := int64(h2 & 0x7F)

    if !masked {
        return 0, nil, fmt.Errorf("client frames must be masked")
    }

    if length == 126 {
        var ext uint16
        binary.Read(r, binary.BigEndian, &ext)
        length = int64(ext)
    } else if length == 127 {
        var ext uint64
        binary.Read(r, binary.BigEndian, &ext)
        length = int64(ext)
    }

    mask := make([]byte, 4)
    io.ReadFull(r, mask)

    buf := make([]byte, length)
    io.ReadFull(r, buf)

    for i := int64(0); i < length; i++ {
        buf[i] ^= mask[i%4]
    }

    return opcode, buf, nil
}

func WriteFrame(w io.Writer, opcode byte, payload []byte) error {
    var header bytes.Buffer

    header.WriteByte(0x80 | opcode)

    l := int64(len(payload))
    switch {
    case l <= 125:
        header.WriteByte(byte(l))
    case l <= 65535:
        header.WriteByte(126)
        binary.Write(&header, binary.BigEndian, uint16(l))
    default:
        header.WriteByte(127)
        binary.Write(&header, binary.BigEndian, uint64(l))
    }

    w.Write(header.Bytes())
    w.Write(payload)
    return nil
}
