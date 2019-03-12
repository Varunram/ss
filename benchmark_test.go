package main

import (
  "testing"
)

func BenchmarkTest(b *testing.B) {
  secret := "shoulder artefact abstract position deny example shoulder myth orchard wolf jewel coconut tourist wrong cram"

  for i := 0 ; i < b.N ; i ++ {
    _, err := Create(2, 3, secret)
    if err != nil {
      b.Fatal(err)
    }
  }
}

func BenchmarkAesEncrypt(b *testing.B) {
  data := []byte("shoulder artefact abstract position deny example shoulder myth orchard wolf jewel coconut tourist wrong cram")
  for i := 0 ; i < b.N ; i ++ {
    _, err := Encrypt(data, "password")
    if err != nil {
      b.Fatal(err)
    }
  }
}

func BenchmarkAesDecrypt(b *testing.B) {
  data := []byte("shoulder artefact abstract position deny example shoulder myth orchard wolf jewel coconut tourist wrong cram")
  x, err := Encrypt(data, "password")
  if err != nil {
    b.Fatal(err)
  }
  for i := 0 ; i < b.N ; i ++ {
    _, err := Decrypt(x, "password")
    if err != nil {
      b.Fatal(err)
    }
  }
}
