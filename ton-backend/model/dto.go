package model

type TonWalletProofDTO struct {
	Address string                 `json:"address"`
	Network string                 `json:"network"`
	Proof   TonWalletProofDTOProof `json:"proof"`
}

type TonWalletProofDTOProof struct {
	Timestamp int64                   `json:"timestamp"`
	Domain    TonWalletProofDTODomain `json:"domain"`
	Signature string                  `json:"signature"`
	Payload   string                  `json:"payload"`
}

type TonWalletProofDTODomain struct {
	LengthBytes int64  `json:"lengthBytes"`
	Value       string `json:"value"`
}

type TonAccountInfo struct {
	Address struct {
		Bounceable    string `json:"bounceable"`
		NonBounceable string `json:"non_bounceable"`
		Raw           string `json:"raw"`
	} `json:"address"`
	Balance int64  `json:"balance"`
	Status  string `json:"status"`
}

type CreateTelegramGroupDTO struct {
	Id              string `json:"id"`
	TwitterUsername string `json:"twitterUsername"`
	IsSecret        bool   `json:"isSecret"`
}

type TelegramGroupResponseDTO struct {
	TelegramGroup

	IsOwner    bool `json:"isOwner"`
	IsApproved bool `json:"isApproved"`
	IsJoined   bool `json:"isJoined"`
}
