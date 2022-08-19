package main

import (
    "io"
    "os"
    "strings"
)

type rot13Reader struct {
    r io.Reader
}

func(reader *rot13Reader) Read(p []byte) (n int, err error) {
    n, err = reader.r.Read(p)
    for i := range(p) {
      p[i] = rot13(p[i])
    }
    return
}

func rot13(b byte) (c byte) {
    switch {
      case b >= 'A' && b <= 'Z':
        c = (b - 'A' + 13) % 26 + 'A'
      case b >= 'a' && b <= 'z':
        c = (b - 'a' + 13) % 26 + 'a'
      default:
        c = b  
    }
    return
}

func test7() {
    s := strings.NewReader(
        "Lbh penpxrq gur pbqr!")
    r := rot13Reader{s}
    io.Copy(os.Stdout, &r)
}