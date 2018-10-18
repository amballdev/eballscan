// Copyright 2018 The go-ecoball Authors
// This file is part of the go-ecoball.
//
// The go-ecoball is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ecoball is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ecoball. If not, see <http://www.gnu.org/licenses/>.

package http

import (
	"net/http"
	"bytes"
	"log"
	"os/exec" 

	"github.com/gin-gonic/gin"
	"github.com/ecoball/eballscan/database"
	"strconv"
)

func StartHttpServer() (err error) {
	//get router instance
	router := gin.Default()

	//register handle
	router.POST("/eballscan/getBlock", getBlock)
	router.POST("/eballscan/getBlockByHeight", getBlockByHeight)
	router.POST("/eballscan/add_block", addBlock)
	router.POST("/eballscan/getTransactionByHash", getTransactionByHash)
	router.POST("/eballscan/add_transaction", addTransaction)
	router.POST("/eballscan/getTransactionByHeight", getTransactionByHeight)
	router.POST("/eballscan/getTransaction", getTransaction)
	router.POST("/eballscan/getTransactionsByAccountName", getTransactionsByAccountName)
	router.POST("/eballscan/getAccounts", getAccounts)
	router.POST("/eballscan/getAccountInfo", getAccountInfo)

	http.ListenAndServe(":20680", router)
	return nil
}

func getBlockByHeight(c *gin.Context) {
	height_str := c.PostForm("height")
	height, err := strconv.Atoi(height_str)
	if nil != err{
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	info, max_height, err := database.QueryOneBlock(height)
	if nil != err{
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"max_height": max_height, "block": info})
}

func getBlock(c *gin.Context) {
	num_str := c.PostForm("num")
	num, err := strconv.Atoi(num_str)
	if nil != err{
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return 
	}

	index_str := c.PostForm("index")
	index, err := strconv.Atoi(index_str)
	if nil != err{
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	info, pageNum, err := database.QueryBlock(index, num)
	if nil != err{
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"pageNum": pageNum, "blocks": info})
}

func getTransaction(c *gin.Context) {
	num_str := c.PostForm("num")
	num, err := strconv.Atoi(num_str)
	if nil != err{
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	index_str := c.PostForm("index")
	index, err := strconv.Atoi(index_str)
	if nil != err{
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	info, pageNum, err := database.QueryTransaction(index, num)
	if nil != err{
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"pageNum": pageNum, "transactions": info})
}

func addBlock(c *gin.Context) {
	height_str := c.PostForm("Height")
	height, err := strconv.Atoi(height_str)
	if nil != err{
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	time_str := c.PostForm("TimeStamp")
	timeStamp, err := strconv.Atoi(time_str)
	if nil != err{
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}
	
	countTxs_str := c.PostForm("CountTxs")
	countTxs, err := strconv.Atoi(countTxs_str)
	if nil != err{
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	hash := c.PostForm("Hash")
	prevHash := c.PostForm("PrevHash")
	merkleHash := c.PostForm("MerkleHash")
	stateHash := c.PostForm("StateHash")
	errcode := database.AddBlock(height, countTxs, timeStamp, hash, prevHash, merkleHash, stateHash)
	if nil != errcode{
		c.JSON(http.StatusBadRequest, gin.H{"result": errcode.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "success"})
}

func getTransactionByHash(c *gin.Context) {
	hash := c.PostForm("hash")
	info, err := database.QueryOneTransaction(hash)
	if nil != err{
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"transaction": info})
}

func addTransaction(c *gin.Context) {
	txType_str := c.PostForm("TxType")
	txType, err := strconv.Atoi(txType_str)
	if nil != err{
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}
	
	timeStamp_str := c.PostForm("TimeStamp")
	timeStamp, err := strconv.Atoi(timeStamp_str)
	if nil != err{
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	blockHeight_str := c.PostForm("BlockHeight")
	blockHeight, err := strconv.Atoi(blockHeight_str)
	if nil != err{
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	hash := c.PostForm("Hash")
	permission := c.PostForm("Permission")
	txFrom := c.PostForm("TxFrom")
	address := c.PostForm("Address")
	errcode := database.AddTransaction(txType, timeStamp, blockHeight, hash, permission, txFrom, address)
	if nil != errcode{
		c.JSON(http.StatusBadRequest, gin.H{"result": errcode.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "success"})
}

func getTransactionByHeight(c *gin.Context) {
	height_str := c.PostForm("blockHeight")
	blockHeight, err := strconv.Atoi(height_str)
	if nil != err{
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}
	
	info, err := database.QueryTransactionsByHeight(blockHeight)
	if nil != err{
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"counts": len(info), "transactions": info})
}

func getTransactionsByAccountName(c *gin.Context) {
	name := c.PostForm("name")

	num_str := c.PostForm("num")
	num, err := strconv.Atoi(num_str)
	if nil != err{
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	index_str := c.PostForm("index")
	index, err := strconv.Atoi(index_str)
	if nil != err{
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	info, pageNum, counts, err := database.QueryTransactionsByAccountName(num, index, name)
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"pageNums": pageNum, "counts": counts, "transactions": info})
}

func getAccounts(c *gin.Context) {
	num_str := c.PostForm("num")
	num, err := strconv.Atoi(num_str)
	if nil != err{
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	index_str := c.PostForm("index")
	index, err := strconv.Atoi(index_str)
	if nil != err{
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	info, pageNum, err := database.QueryAccounts(num, index)
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"pageNum": pageNum, "accounts": info})
}

func getAccountInfo(c *gin.Context) {
	name := c.PostForm("name")

	data, err := database.QueryOneAccount(name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"account": data})
}

func exec_shell(s string) {
    cmd := exec.Command("/bin/bash", "-c", s)
    var out bytes.Buffer

    cmd.Stdout = &out
    err := cmd.Run()
    if err != nil {
        log.Fatal(err)
    }
}

func Compile() {
	exec_shell("cd ../c2wasm-compiler/;./compile.sh api.c api;cd ../eballscan/")
}
