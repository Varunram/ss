package main

import (
  "os"
  "log"
  "strconv"
  "io/ioutil"
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

    // encrypt using random alphanumeric string
    for i, elem := range a {
      log.Println("SECRET: ", i+1 , elem)
      // encrypt using a random string and print the random string for future use.
      rString := GetRandomString(6)
      byteData, err := Encrypt([]byte(elem), rString)
      if err != nil {
        log.Println("ERR: ", err)
        log.Fatal(err)
      }
      // implementation needs to store this to a file and decrypt this using the random string
      log.Println("Lenght of encrpyted share: ", len(byteData))
      err = ioutil.WriteFile(strconv.Itoa(i+1) + ".secret", byteData, os.ModePerm)
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

  case "encrypt":
    if len(os.Args) < 3 {
      log.Fatal("USAGE: ./ss encrypt secret(s)")
    }
    var inputString string
    // this would ideally be a bip39 mnemonic
    for i := 2 ; i < len(os.Args) ; i ++ {
      inputString = inputString + " " + os.Args[i]
    }
    encoded, err := EncryptQuestions(inputString)
    if err != nil {
      log.Println("ERROR: ", err)

    }
    log.Println("ENCODED: ", encoded)

  case "decrypt":
    if len(os.Args) < 3 {
      log.Fatal("USAGE: ./ss decrypt passphrase")
    }
    secret, err := DecryptQuestions(os.Args[2])
    if err != nil {
      log.Println("ERROR: ", err)
    }
    log.Println("SECRET PHRASE: ", secret)
  }
}
