package main

import "testing"

import "git.ngx.fi/c0mm4nd/tronetl/tron"

// buildABIEntry is a tiny helper for test clarity.
func buildABIEntry(name string, inputs ...string) *tron.HTTPABIEntry {
	params := make([]*tron.HTTPABIParams, 0, len(inputs))
	for _, in := range inputs {
		params = append(params, &tron.HTTPABIParams{Type: in})
	}
	return &tron.HTTPABIEntry{Type: "function", Name: name, Inputs: params}
}

func TestNewCsvContractDetectsERC1155(t *testing.T) {
	abi := []*tron.HTTPABIEntry{
		buildABIEntry("balanceOf", "address", "uint256"),
		buildABIEntry("balanceOfBatch", "address[]", "uint256[]"),
		buildABIEntry("safeTransferFrom", "address", "address", "uint256", "uint256", "bytes"),
		buildABIEntry("safeBatchTransferFrom", "address", "address", "uint256[]", "uint256[]", "bytes"),
		buildABIEntry("supportsInterface", "bytes4"),
	}
	contract := &tron.HTTPContract{
		ContractAddress: tron.EnsureHexAddr("TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t"),
		OriginAddress:   tron.EnsureHexAddr("THPvaUhoh2Qn2y9THCZML3H815hhFhn5YC"),
		Abi: struct {
			Entrys []*tron.HTTPABIEntry `json:"entrys,omitempty"`
		}{Entrys: abi},
	}

	csvCtr := NewCsvContract(contract)

	if !csvCtr.IsErc1155 {
		t.Fatalf("expected IsErc1155 true")
	}
	if csvCtr.IsErc20 || csvCtr.IsErc721 {
		t.Fatalf("ERC1155 detection should not imply ERC20/721")
	}
}

func TestNewCsvContractDoesNotMisdetect1155(t *testing.T) {
	abi := []*tron.HTTPABIEntry{
		buildABIEntry("balanceOf", "address"),
		buildABIEntry("transfer", "address", "uint256"),
		buildABIEntry("approve", "address", "uint256"),
	}
	contract := &tron.HTTPContract{
		ContractAddress: tron.EnsureHexAddr("THPvaUhoh2Qn2y9THCZML3H815hhFhn5YC"),
		OriginAddress:   tron.EnsureHexAddr("THPvaUhoh2Qn2y9THCZML3H815hhFhn5YC"),
		Abi: struct {
			Entrys []*tron.HTTPABIEntry `json:"entrys,omitempty"`
		}{Entrys: abi},
	}

	csvCtr := NewCsvContract(contract)
	if csvCtr.IsErc1155 {
		t.Fatalf("expected IsErc1155 false")
	}
}
