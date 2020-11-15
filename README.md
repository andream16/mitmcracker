# mitmcracker ![](https://stuff.mit.edu/afs/sipb/project/golang/arch/go1.2.1-linux-amd64/favicon.ico)

## Golang Key Cracker Which Implements Meet In the Middle Strategy.

Given a known `cipher-text` and a known `plain-text`, it finds two keys `k1` and `k2` in a way that 
`Ek1(plain) == Dk2(cipher)`
where `E` is the encryption function and `D` is the decryption function.

It tries all the possible keys in the range of `24/28/32 bits` and prints which are the keys to be used for the encryption and decryption.

Keys can only be `6/7/8 Digits HEX` and `Cipher/Plain` can only be `16 Digits HEX`.

## What you need to run the project:

- Go: ~v1.15
- Ocaml: ~v4.08.1

## How to run the app

From the root of the project:

- build the go/ocaml binaries: `make`
- run the command: `./mitmcracker -key=$keyLength -encoded=$encodedMsg -plain=$plainMsg`

## Known Inputs/Outputs

- 24-bit:
    - plain: C330C9CBD01DFBA0
    - encoded: E10C65124518DB05
    - encoding key to be found: 6d3952
    - decoding key to be found: 513346
- 28-bit:
    - plain: 492A5BB83F8A3F95
    - encoded: 47090A6AC4A56798
    - encoding key to be found: ???????
    - decoding key to be found: ???????
- 32-bit:
    - plain: FFC7C9E5694ABFF7
    - encoded: 98AC59F25448FFAC
    - encoding key to be found: ????????
    - decoding key to be found: ????????

When `Ek1(C) == Dk2(P)` you found the right keys `k1, k2` to read all encrypted messages.