package swap

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/mises-id/mises-swapsvc/app/models"
	"github.com/mises-id/mises-swapsvc/app/models/search"
)

type (
	SwapChainJson struct {
		Lists []*models.SwapChain `json:"list"`
	}
	SwapTokenJson struct {
		Lists []*models.SwapToken `json:"list"`
	}
	SwapContractJson struct {
		Lists []*models.SwapContract `json:"list"`
	}
	SwapProviderJson struct {
		Lists []*models.SwapProvider `json:"list"`
	}
)

var (
	tokenListMap                = map[string]*models.SwapToken{}
	last_update_token_list_time time.Time
)

func isNeedUpdateSwapTokenMap() bool {
	end_time := last_update_token_list_time.Add(time.Minute * 10)
	return time.Now().After(end_time)
}

// GetTokenByAddress
func GetTokenByAddress(ctx context.Context, chainID uint64, address string) *models.SwapToken {
	tokenMap := GetSwapTokenMap(ctx)
	tokenKey := getTokenMapKey(chainID, address)
	token, ok := tokenMap[tokenKey]
	if !ok {
		return nil
	}
	return token
}

func GetSwapTokenMap(ctx context.Context) map[string]*models.SwapToken {
	if isNeedUpdateSwapTokenMap() {
		updateSwapToken(ctx)
	}
	return tokenListMap
}

func updateSwapToken(ctx context.Context) error {
	fmt.Printf("[%s] UpdateSwapToken \n", time.Now().Local().String())
	params := &search.SwapTokenSearch{}
	list, err := models.ListSwapToken(context.Background(), params)
	if err != nil {
		fmt.Println("UpdateSwapToken warning: ", err.Error())
		return err
	}
	for _, token := range list {
		key := getTokenMapKey(token.ChainID, token.Address)
		tokenListMap[key] = token
	}
	last_update_token_list_time = time.Now()
	return nil
}

func getTokenMapKey(chainID uint64, address string) string {
	return fmt.Sprintf("%d&%s", chainID, address)
}

// UpdateSwapChain
func UpdateSwapChain(ctx context.Context) error {
	err := runUpdateSwapChain(ctx)
	if err != nil {
		fmt.Printf("Error updating chain list: %s", err.Error())
		return err
	}
	return nil
}

func runUpdateSwapChain(ctx context.Context) error {
	lists, err := getSwapChainByJSON()
	if err != nil {
		return err
	}
	return models.CreateSwapChainMany(ctx, lists)
}

func getSwapChainByJSON() ([]*models.SwapChain, error) {
	//local json
	localfile := path.Join("./assets/swap/chains.json")
	jsonFile, err := os.Open(localfile)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()
	out := &SwapChainJson{}
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(byteValue, out)
	return out.Lists, nil
}

// UpdateSwapToken
func UpdateSwapToken(ctx context.Context) error {
	err := runUpdateSwapToken(ctx)
	if err != nil {
		fmt.Printf("Error updating token list: %s", err.Error())
		return err
	}
	return nil
}

func runUpdateSwapToken(ctx context.Context) error {
	lists, err := getSwapTokenByJSON()
	if err != nil {
		return err
	}
	return models.CreateSwapTokenMany(ctx, lists)
}

func getSwapTokenByJSON() ([]*models.SwapToken, error) {
	//local json
	localfile := path.Join("./assets/swap/tokens.json")
	jsonFile, err := os.Open(localfile)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()
	out := &SwapTokenJson{}
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(byteValue, out)
	return out.Lists, nil
}

// UpdateSwapContract
func UpdateSwapContract(ctx context.Context) error {
	err := runUpdateSwapContract(ctx)
	if err != nil {
		fmt.Printf("Error updating contract list: %s", err.Error())
		return err
	}
	return nil
}

func runUpdateSwapContract(ctx context.Context) error {
	lists, err := getSwapContractByJSON()
	if err != nil {
		return err
	}
	return models.CreateSwapContractMany(ctx, lists)
}

func getSwapContractByJSON() ([]*models.SwapContract, error) {
	//local json
	localfile := path.Join("./assets/swap/contracts.json")
	jsonFile, err := os.Open(localfile)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()
	out := &SwapContractJson{}
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(byteValue, out)
	return out.Lists, nil
}

// UpdateSwapProvider
func UpdateSwapProvider(ctx context.Context) error {
	err := runUpdateSwapProvider(ctx)
	if err != nil {
		fmt.Printf("Error updating contract list: %s", err.Error())
		return err
	}
	return nil
}

func runUpdateSwapProvider(ctx context.Context) error {
	lists, err := getSwapProviderByJSON()
	if err != nil {
		return err
	}
	return models.CreateSwapProviderMany(ctx, lists)
}

func getSwapProviderByJSON() ([]*models.SwapProvider, error) {
	//local json
	localfile := path.Join("./assets/swap/providers.json")
	jsonFile, err := os.Open(localfile)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()
	out := &SwapProviderJson{}
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(byteValue, out)
	return out.Lists, nil
}
