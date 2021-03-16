// Copyright (c) 2014-2016 The btcsuite developers
// Copyright (c) 2015-2021 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package txscript_test

import (
	"encoding/hex"
	"fmt"

	"github.com/decred/dcrd/chaincfg/chainhash"
	"github.com/decred/dcrd/chaincfg/v3"
	"github.com/decred/dcrd/dcrec"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/decred/dcrd/txscript/v4"
	"github.com/decred/dcrd/txscript/v4/stdaddr"
	"github.com/decred/dcrd/wire"
)

const (
	// noTreasury signifies the treasury agenda should be treated as though
	// it is inactive.  It is used to increase the readability of the
	// tests.
	noTreasury = false
)

// This example demonstrates extracting information from a standard public key
// script.
func ExampleExtractPkScriptAddrs() {
	// Start with a standard pay-to-pubkey-hash script.
	const scriptVersion = 0
	scriptHex := "76a914128004ff2fcaf13b2b91eb654b1dc2b674f7ec6188ac"
	script, err := hex.DecodeString(scriptHex)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Extract and print details from the script.
	mainNetParams := chaincfg.MainNetParams()
	scriptClass, addresses, reqSigs, err := txscript.ExtractPkScriptAddrs(
		scriptVersion, script, mainNetParams, noTreasury)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Script Class:", scriptClass)
	fmt.Println("Addresses:", addresses)
	fmt.Println("Required Signatures:", reqSigs)

	// Output:
	// Script Class: pubkeyhash
	// Addresses: [DsSej1qR3Fyc8kV176DCh9n9cY9nqf9Quxk]
	// Required Signatures: 1
}

// This example demonstrates manually creating and signing a redeem transaction.
func ExampleSignTxOutput() {
	// Ordinarily the private key would come from whatever storage mechanism
	// is being used, but for this example just hard code it.
	privKeyBytes, err := hex.DecodeString("22a47fa09a223f2aa079edf85a7c2" +
		"d4f8720ee63e502ee2869afab7de234b80c")
	if err != nil {
		fmt.Println(err)
		return
	}
	pubKey := secp256k1.PrivKeyFromBytes(privKeyBytes).PubKey()
	pubKeyHash := stdaddr.Hash160(pubKey.SerializeCompressed())
	mainNetParams := chaincfg.MainNetParams()
	addr, err := stdaddr.NewAddressPubKeyHashEcdsaSecp256k1V0(pubKeyHash,
		mainNetParams)
	if err != nil {
		fmt.Println(err)
		return
	}

	// For this example, create a fake transaction that represents what
	// would ordinarily be the real transaction that is being spent.  It
	// contains a single output that pays to address in the amount of 1 DCR.
	originTx := wire.NewMsgTx()
	prevOut := wire.NewOutPoint(&chainhash.Hash{}, ^uint32(0), wire.TxTreeRegular)
	txIn := wire.NewTxIn(prevOut, 100000000, []byte{txscript.OP_0, txscript.OP_0})
	originTx.AddTxIn(txIn)
	pkScriptVer, pkScript := addr.PaymentScript()
	txOut := wire.NewTxOut(100000000, pkScript)
	txOut.Version = pkScriptVer
	originTx.AddTxOut(txOut)
	originTxHash := originTx.TxHash()

	// Create the transaction to redeem the fake transaction.
	redeemTx := wire.NewMsgTx()

	// Add the input(s) the redeeming transaction will spend.  There is no
	// signature script at this point since it hasn't been created or signed
	// yet, hence nil is provided for it.
	prevOut = wire.NewOutPoint(&originTxHash, 0, wire.TxTreeRegular)
	txIn = wire.NewTxIn(prevOut, 100000000, nil)
	redeemTx.AddTxIn(txIn)

	// Ordinarily this would contain that actual destination of the funds,
	// but for this example don't bother.
	txOut = wire.NewTxOut(0, nil)
	redeemTx.AddTxOut(txOut)

	// Sign the redeeming transaction.
	sigType := dcrec.STEcdsaSecp256k1
	lookupKey := func(a stdaddr.Address) ([]byte, dcrec.SignatureType, bool, error) {
		// Ordinarily this function would involve looking up the private
		// key for the provided address, but since the only thing being
		// signed in this example uses the address associated with the
		// private key from above, simply return it with the compressed
		// flag set since the address is using the associated compressed
		// public key.
		//
		// NOTE: If you want to prove the code is actually signing the
		// transaction properly, uncomment the following line which
		// intentionally returns an invalid key to sign with, which in
		// turn will result in a failure during the script execution
		// when verifying the signature.
		//
		// privKey.D.SetInt64(12345)
		//
		return privKeyBytes, sigType, true, nil
	}
	// Notice that the script database parameter is nil here since it isn't
	// used.  It must be specified when pay-to-script-hash transactions are
	// being signed.
	sigScript, err := txscript.SignTxOutput(mainNetParams, redeemTx, 0,
		originTx.TxOut[0].PkScript, txscript.SigHashAll,
		txscript.KeyClosure(lookupKey), nil, nil, noTreasury)
	if err != nil {
		fmt.Println(err)
		return
	}
	redeemTx.TxIn[0].SignatureScript = sigScript

	// Prove that the transaction has been validly signed by executing the
	// script pair.

	flags := txscript.ScriptDiscourageUpgradableNops
	vm, err := txscript.NewEngine(originTx.TxOut[0].PkScript, redeemTx, 0,
		flags, 0, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	if err := vm.Execute(); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Transaction successfully signed")

	// Output:
	// Transaction successfully signed
}

// This example demonstrates creating a script tokenizer instance and using it
// to count the number of opcodes a script contains.
func ExampleScriptTokenizer() {
	// Create a script to use in the example.  Ordinarily this would come from
	// some other source.
	hash160 := stdaddr.Hash160([]byte("example"))
	script, err := txscript.NewScriptBuilder().AddOp(txscript.OP_DUP).
		AddOp(txscript.OP_HASH160).AddData(hash160).
		AddOp(txscript.OP_EQUALVERIFY).AddOp(txscript.OP_CHECKSIG).Script()
	if err != nil {
		fmt.Printf("failed to build script: %v\n", err)
		return
	}

	// Create a tokenizer to iterate the script and count the number of opcodes.
	const scriptVersion = 0
	var numOpcodes int
	tokenizer := txscript.MakeScriptTokenizer(scriptVersion, script)
	for tokenizer.Next() {
		numOpcodes++
	}
	if tokenizer.Err() != nil {
		fmt.Printf("script failed to parse: %v\n", err)
	} else {
		fmt.Printf("script contains %d opcode(s)\n", numOpcodes)
	}

	// Output:
	// script contains 5 opcode(s)
}
