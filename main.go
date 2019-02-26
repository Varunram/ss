package main

import (
  "os"
  "log"
  "strconv"

  qrcode "github.com/skip2/go-qrcode"
)

// TODO: add questions which a user can answer and use the answers to encrypt the secret phrase
// TODO: split the secret in such a way that the data is easily typable without any problems.
// Maybe bech32 or base58 would be a good candidate.
// TODO: remove qr codes
// TODO: add encrpytion scheme which would accept answers for a set of questions
func main() {
  // sample pubkey: 2NFMAdgT7tWdLYbMGhxBvnYQ8h3nLzHDNbq
  // sample privkey: cRGgZNBv9iAokqj64YAzu5PUwydxC272yczQr56Y9KiqrrNvkNR7
  log.Println("OS ARGS: ", os.Args)
  if len(os.Args) < 2 {
    log.Fatal("SUPPORTED COMMANDS: new, reconstruct")
  }
  arg := os.Args[1]
  switch(arg) {
  case "new":
    if len(os.Args) < 5 {
      log.Fatal("USAGE: new <secret> <m> <n> \nGenerate secrets using the m of n Shamir Secret Sharing Scheme")
    }
    secret := os.Args[2]
    m1 := os.Args[3]
    n1 := os.Args[4]
    m, err := strconv.Atoi(m1)
    if err != nil {
      log.Fatal(err)
    }

    n, err := strconv.Atoi(n1)
    if err != nil {
      log.Fatal(err)
    }

    a, err := Create(m, n, secret)
    if err != nil {
      log.Println(err)
    }

    for i, elem := range a {
      log.Println("SECRET: ", i+1 , elem)
    }

    for i, elem := range a {
      sI := strconv.Itoa(i+1)
      err := qrcode.WriteFile(elem, qrcode.Medium, 256, sI + ".png")
      if err != nil {
        log.Fatal(err)
      }
    }

  case "reconstruct":
    if len(os.Args) < 4 {
      log.Fatal("USAGE: reconstruct <secret1> <secret2>\nReconstruct the secret using the Shamir Secret Sharing Scheme")
    }
    var arr []string
    arr = append(arr, os.Args[2])
    arr = append(arr, os.Args[3])
    x, err := Combine(arr)
    if err != nil {
      log.Println(err)
    }
    log.Println("RECONSTRUCTED: ", x)
  }
}
