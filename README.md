Go Threshold Secret Sharing
===

[![Build Status](https://travis-ci.org/antik10ud/go-tss.svg?branch=master)](https://travis-ci.org/antik10ud/go-tss)

## Description
Core implementation of Threshold Secret Sharing (TSS) [http://tools.ietf.org/html/draft-mcgrew-tss-03](http://tools.ietf.org/html/draft-mcgrew-tss-03)
  It uses a finite field GF(256) instead of Shamir scheme using large integers modulo a large prime number. 
	
## State
alpha, work in progress

## Usage

	sharesCount := 5 // number of shares
	threshold := 3 // number of requires shares to recover the secret
	shares, err := CreateShares(secret, sharesCount, threshold)
	[...]
	//use 3 shares to recover the secret
	recoveredSecret, _ := RecoverSecret(ShareSet{shares[0], shares[1], shares[4]})
	


## Doc
[![GoDoc](https://godoc.org/github.com/antik10ud/go-tss?status.svg)](https://godoc.org/github.com/antik10ud/go-tss)

## Lint
[http://go-lint.appspot.com/github.com/antik10ud/go-tss] (Run lint) on tss.

## License
MIT License
