package swap_sync

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math/big"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/mises-id/mises-swapsvc/app/models"
	"github.com/mises-id/mises-swapsvc/app/models/enum"
	"github.com/mises-id/mises-swapsvc/app/models/search"
	"github.com/mises-id/mises-swapsvc/config/env"
)

type (
	config struct {
		baseSleepTime time.Duration
		offsetNum     uint64 //page offset
	}
	SwapSync struct {
		syncFlag, isRun bool
		swapContractMap map[string]*models.SwapContract
		ctx             context.Context
		config          *config
		wg              sync.WaitGroup
	}
	decodeTxInputData struct {
		SrcToken        common.Address `json:"srcToken"`
		DstToken        common.Address `json:"dstToken"`
		SrcReceiver     common.Address `json:"srcReceiver"`
		DstReceiver     common.Address `json:"dstReceiver"`
		Amount          *big.Int       `json:"amount"`
		MinReturnAmount *big.Int       `json:"minReturnAmount"`
		Flags           *big.Int       `json:"flags"`
	}
)

var (
	swapSyncEntity      *SwapSync
	swapReferrerAddress = env.Envs.SwapReferrerAddress
)

const (
	minSleepTimeMs     = time.Millisecond * 1000
	defaultSleepTimeMs = time.Millisecond * 10000
)

func NewSwapSync() *SwapSync {
	if swapSyncEntity != nil {
		return swapSyncEntity
	}
	swapSyncEntity = &SwapSync{
		syncFlag:        true,
		ctx:             context.TODO(),
		swapContractMap: map[string]*models.SwapContract{},
		config:          &config{baseSleepTime: time.Millisecond * 100, offsetNum: 10000},
		isRun:           false,
	}
	return swapSyncEntity
}

func (ctrl *SwapSync) init() {
	ctrl.initSyncContract()
}

func (ctrl *SwapSync) initSyncContract() {
	params := &search.SwapContractSearch{Status: 1}
	lists, err := models.ListSwapContract(ctrl.ctx, params)
	if err != nil {
		fmt.Printf("[%s] SwapSync ListSwapContract error: %s\n", time.Now().Local().String(), err.Error())
		return
	}
	for _, contract := range lists {
		address := contract.Address
		if address == "" {
			continue
		}
		ctrl.swapContractMap[getContractMapKey(contract)] = contract
	}
}

func getContractMapKey(contract *models.SwapContract) string {
	return fmt.Sprintf("%d&%s", contract.ChainID, contract.Address)
}

func (ctrl *SwapSync) Run() error {
	//check
	err := beforeSwapSync(ctrl.ctx)
	if err != nil {
		return err
	}
	swapSyncEntity.init()
	go ctrl.run()
	fmt.Println("SwapSync Run")
	return nil
}

func (ctrl *SwapSync) run() {
	if ctrl.isRun {
		ctrl.Stop()
		return
	}
	ctrl.isRun = true
	ctrl.wg.Add(1)
	ctrl.startSyncTask()
	ctrl.wg.Wait()
}

func (ctrl *SwapSync) startSyncTask() {
	fmt.Println("SwapSync taskChain Start")
	for _, contract := range ctrl.swapContractMap {
		go ctrl.runSyncByContract(contract)
	}
}

func (ctrl *SwapSync) runSyncByContract(contract *models.SwapContract) {
	if contract == nil || contract.Chain == nil || contract.Chain.ScanApi == "" {
		fmt.Printf("[%s] SwapSync runSyncByContract contract error \n", time.Now().Local().String())
		return
	}
	ctx := ctrl.ctx
	contractName := contract.Name
	contractABI, err := getABIByContractAddress(contract)
	chainID := contract.ChainID
	address := contract.Address
	if err != nil {
		fmt.Printf("[%s][%s] SwapSync getABI error: %s\n", time.Now().Local().String(), contractName, err.Error())
		return
	}
	fmt.Printf("[%s] [%s] SwapSync runSyncByContract start chainID=%d \n", time.Now().Local().String(), contractName, chainID)
	ctrl.wg.Add(1)
	defer ctrl.wg.Done()
	var startBlock, endBlock, stepBlock, safeBlock, defaultSafeBlock, pageNum uint64
	scanEndpoint := contract.Chain.ScanApi
	scanApiKey := contract.Chain.ScanApiKey
	startBlock = contract.SyncBlockNumber
	stepBlock = 50
	defaultSafeBlock = 10
	safeBlock = defaultSafeBlock
	offsetNum := ctrl.config.offsetNum
	pageNum = 1
	sort := "asc"
	sleepTimeMs := time.Millisecond * time.Duration(contract.Chain.SyncInterval)
	if sleepTimeMs == 0 {
		sleepTimeMs = defaultSleepTimeMs
	}
	if sleepTimeMs < minSleepTimeMs {
		sleepTimeMs = minSleepTimeMs
	}
	shouldSleep := false
	for ctrl.isRun && ctrl.checkContractTaskIsRun(contract) {
		if shouldSleep {
			time.Sleep(sleepTimeMs)
		}
		shouldSleep = true
		//maybeHasNextPage := false
		//get txs by contract address
		chainRecentBlockNumber, err := apiEthBlockNumber(ctx, scanEndpoint, scanApiKey)
		//fmt.Printf("[%s] SwapSync apiEthBlockNumber chainID=%d  contract=%s %d-%d\n", time.Now().Local().String(), chainID, address, contract.SyncBlockNumber, chainRecentBlockNumber)
		if err != nil {
			fmt.Printf("[%s] SwapSync apiEthBlockNumber chainID=%d  contract=%s error:%s\n", time.Now().Local().String(), chainID, address, err.Error())
			continue
		}
		//no block number need to sync
		if chainRecentBlockNumber <= contract.SyncBlockNumber {
			continue
		}
		startBlock -= safeBlock
		safeBlock = defaultSafeBlock
		if startBlock <= 0 || startBlock > chainRecentBlockNumber {
			startBlock = chainRecentBlockNumber
		}
		endBlock = startBlock + stepBlock
		if endBlock > chainRecentBlockNumber {
			endBlock = chainRecentBlockNumber
		}
		params := &GetTransactionsByAddressRequest{
			Address:    address,
			StartBlock: startBlock,
			EndBlock:   endBlock,
			Page:       pageNum,
			Offset:     offsetNum,
			Sort:       sort,
		}
		fmt.Printf("[%s] SyncSwap runTaskChain chainID=%d  %d-%d\n", time.Now().Local().String(), chainID, startBlock, endBlock)
		txs, err := apiGetTransactionsByAddress(ctx, scanEndpoint, scanApiKey, params)
		if err != nil {
			if endBlock > 0 && strings.Contains(err.Error(), "No transactions found") {
				//this page no data need to update startBlock
				pageNum = 1
				startBlock = endBlock
				if err := models.ChangeSwapContractBlockNumber(ctx, contract.ID, endBlock); err != nil {
					fmt.Printf("[%s] SwapSync ChangeSwapContractBlockNumber chainID=%d  contract=%s error:%s\n", time.Now().Local().String(), chainID, address, err.Error())
					continue
				}
				contract.SyncBlockNumber = endBlock
				continue
			}
			fmt.Printf("[%s] SyncSwap runTaskChain chainID=%d  %d-%d Error: %s\n", time.Now().Local().String(), chainID, startBlock, endBlock, err.Error())
			safeBlock = 0
			continue
		}
		if txs == nil {
			fmt.Printf("[%s] SyncSwap runTaskChain error txs is nil chainID=%d  %d-%d\n", time.Now().Local().String(), chainID, startBlock, endBlock)
			continue
		}
		txNum := len(txs)
		if txNum == 0 {
			continue
		}
		for _, tx := range txs {
			if !checkIsMisesTx(tx) {
				continue
			}
			//mises tx
			//create swap order
			blockTime, _ := strconv.ParseInt(tx.Timestamp, 10, 64)
			blockAt := time.Unix(blockTime, 0)
			order := &models.SwapOrder{
				ChainID:         chainID,
				FromAddress:     strings.ToLower(tx.From),
				Transaction:     tx,
				BlockAt:         &blockAt,
				ProviderKey:     contract.ProviderKey,
				ContractAddress: strings.ToLower(address),
			}
			status := enum.SwapOrderPending
			if tx.TxreceiptStatus == "1" {
				status = enum.SwapOrderSuccess
			} else {
				status = enum.SwapOrderFail
			}
			order.ReceiptStatus = status
			inputDecode, err := decodeTransactionInputData(contractABI, contract, tx)
			if err != nil {
				fmt.Printf("[%s] SwapSync apiGetTransactionsByAddress chainID=%d  txHash=%s error:%s\n", time.Now().Local().String(), chainID, tx.Hash, err.Error())
				continue
			}
			if inputDecode == nil {
				fmt.Printf("[%s] SwapSync apiGetTransactionsByAddress chainID=%d  txHash=%s input decode is null \n", time.Now().Local().String(), chainID, tx.Hash)
				continue
			}
			srcToken := &models.FromToken{
				Address: strings.ToLower(inputDecode.SrcToken.Hex()),
				Value:   inputDecode.Amount.String(),
			}
			order.FromToken = srcToken
			destToken := &models.ToToken{
				Address: strings.ToLower(inputDecode.DstToken.Hex()),
			}
			order.ToToken = destToken
			order.MinReturnAmount = inputDecode.MinReturnAmount.String()
			order.DestReceiver = strings.ToLower(inputDecode.DstReceiver.Hex())
			_, err = models.CreateSwapOrder(ctx, order)
			if err != nil {
				fmt.Printf("[%s] SwapSync CreateSwapOrder chainID=%d  txHash=%s error:%s\n", time.Now().Local().String(), chainID, tx.Hash, err.Error())
				continue
			}
			fmt.Printf("[%s] SwapSync CreateSwapOrder Success chainID=%d  txHash=%s\n", time.Now().Local().String(), chainID, tx.Hash)
		}
		/* if txNum == int(offsetNum) {
			maybeHasNextPage = true
		}
		//maybeHasNextPage
		if maybeHasNextPage {
			pageNum++
			shouldSleep = false
			continue
		}
		pageNum = 1 */
		startBlock = endBlock
		if err := models.ChangeSwapContractBlockNumber(ctx, contract.ID, endBlock); err != nil {
			fmt.Printf("[%s] SwapSync ChangeSwapContractBlockNumber chainID=%d  contract=%s error:%s\n", time.Now().Local().String(), chainID, address, err.Error())
			continue
		}
		contract.SyncBlockNumber = endBlock
		continue
	}
}

func (ctrl *SwapSync) checkContractTaskIsRun(contract *models.SwapContract) bool {
	_, ok := ctrl.swapContractMap[getContractMapKey(contract)]
	return ok
}

func (ctrl *SwapSync) Stop() error {
	if !ctrl.isRun {
		return nil
	}
	fmt.Println("SwapSync Stop")
	ctrl.syncFlag = false
	ctrl.wg.Done()
	ctrl.isRun = false
	return nil
}

func beforeSwapSync(ctx context.Context) error {
	if swapReferrerAddress == "" {
		return errors.New("Referrer address invalid")
	}
	return nil
}

func decodeTransactionInputData(contractABI *abi.ABI, contract *models.SwapContract, tx *models.Transaction) (decodeInput *decodeTxInputData, err error) {
	if contractABI == nil || contract == nil || tx == nil {
		return nil, errors.New("Error contractABI or contract or transaction is nil")
	}
	data := common.FromHex(tx.Input)
	methodSigData := data[:4]
	method, err := contractABI.MethodById(methodSigData)
	if err != nil {
		return nil, err
	}
	inputsSigData := data[4:]
	inputsMap := make(map[string]interface{})
	if err := method.Inputs.UnpackIntoMap(inputsMap, inputsSigData); err != nil {
		return nil, err
	}
	desc, ok := inputsMap["desc"]
	if !ok {
		return nil, errors.New("Error Unpack inputsMap")
	}
	tempDecodeInput := desc.(struct {
		SrcToken        common.Address `json:"srcToken"`
		DstToken        common.Address `json:"dstToken"`
		SrcReceiver     common.Address `json:"srcReceiver"`
		DstReceiver     common.Address `json:"dstReceiver"`
		Amount          *big.Int       `json:"amount"`
		MinReturnAmount *big.Int       `json:"minReturnAmount"`
		Flags           *big.Int       `json:"flags"`
	})
	decodeInput = &decodeTxInputData{}
	decodeInput.SrcToken = tempDecodeInput.SrcToken
	decodeInput.DstToken = tempDecodeInput.DstToken
	decodeInput.SrcReceiver = tempDecodeInput.SrcReceiver
	decodeInput.DstReceiver = tempDecodeInput.DstReceiver
	decodeInput.Amount = tempDecodeInput.Amount
	decodeInput.MinReturnAmount = tempDecodeInput.MinReturnAmount
	decodeInput.Flags = tempDecodeInput.Flags
	return decodeInput, nil
}

func getABIByContractAddress(contract *models.SwapContract) (*abi.ABI, error) {
	if contract == nil || contract.Address == "" {
		return nil, errors.New("Invalid contract")
	}
	abiFile := fmt.Sprintf("./assets/swap/%s.json", contract.Address)
	abiJson, err := getLocalABI(abiFile)
	if err != nil {
		return nil, err
	}
	contractABI, err := abi.JSON(strings.NewReader(abiJson))
	if err != nil {
		return nil, err
	}
	return &contractABI, nil
}

func getLocalABI(path string) (string, error) {
	abiFile, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer abiFile.Close()

	result, err := io.ReadAll(abiFile)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

// checkIsMisesTx
func checkIsMisesTx(tx *models.Transaction) bool {
	if tx == nil {
		return false
	}
	address := strings.ToLower(swapReferrerAddress)
	if strings.HasPrefix(address, "0x") {
		address = strings.TrimPrefix(address, "0x")
	}
	if strings.Contains(strings.ToLower(tx.Input), address) {
		return true
	}
	return false
}
