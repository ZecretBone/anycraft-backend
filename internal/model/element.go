package model

type Element struct {
	ID            int     `json:"id"`
	Slug          string  `json:"slug"`
	Name          string  `json:"name"`
	Emoji         *string `json:"emoji,omitempty"`
	IsCharacter   bool    `json:"is_character"`
	IsBaseElement bool    `json:"is_base_element"`
	ImageURL      *string `json:"image_url,omitempty"`
	Rarity        *string `json:"rarity,omitempty"`
	Difficulty    *int    `json:"difficulty,omitempty"`
}

type CombineRequest struct {
	GameCode  string `json:"game_code"`
	ParentAID int    `json:"parent_a_id"`
	ParentBID int    `json:"parent_b_id"`
}

type CombineResponse struct {
	OK     bool     `json:"ok"`
	Error  *string  `json:"error,omitempty"`
	Result *Element `json:"result,omitempty"`
}

type ChallengesRequest struct {
	GameCode               string `json:"game_code"`
	DiscoveredCharacterIDs []int  `json:"discovered_character_ids"`
}

type ChallengeItem struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	ImageURL *string `json:"image_url,omitempty"`
}

type ChallengesResponse struct {
	OK    bool            `json:"ok"`
	Items []ChallengeItem `json:"items"`
}

type BaseElementsResponse struct {
	OK    bool      `json:"ok"`
	Items []Element `json:"items"`
}
