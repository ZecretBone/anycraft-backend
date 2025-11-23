package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/fpswan/anycraft-backend/internal/model"
)

type ComposeRepository struct {
	DB *pgxpool.Pool
}

func NewComposeRepository(db *pgxpool.Pool) *ComposeRepository {
	return &ComposeRepository{DB: db}
}

func (r *ComposeRepository) GetBaseElements(ctx context.Context, gameCode string) ([]model.Element, error) {
	rows, err := r.DB.Query(ctx, `
		SELECT e.element_id, e.slug, e.name, e.emoji, e.is_character, e.is_base_element, e.image_url, e.rarity, e.difficulty
		FROM element e
		JOIN game g ON e.game_id = g.game_id
		WHERE g.code=$1 AND e.is_base_element=true
		ORDER BY e.name
	`, gameCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Element
	for rows.Next() {
		var e model.Element
		if err := rows.Scan(&e.ID, &e.Slug, &e.Name, &e.Emoji, &e.IsCharacter, &e.IsBaseElement, &e.ImageURL, &e.Rarity, &e.Difficulty); err != nil {
			return nil, err
		}
		list = append(list, e)
	}
	return list, nil
}

func (r *ComposeRepository) Combine(ctx context.Context, gameCode string, a, b int) (*model.Element, error) {
	if a > b { a, b = b, a }
	row := r.DB.QueryRow(ctx, `
		SELECT e.element_id, e.slug, e.name, e.emoji, e.is_character, e.is_base_element, e.image_url, e.rarity, e.difficulty
		FROM recipe r
		JOIN game g ON r.game_id = g.game_id
		JOIN element e ON e.element_id = r.result_id
		WHERE g.code=$1 AND r.parent_a_id=$2 AND r.parent_b_id=$3
	`, gameCode, a, b)

	var e model.Element
	if err := row.Scan(&e.ID, &e.Slug, &e.Name, &e.Emoji, &e.IsCharacter, &e.IsBaseElement, &e.ImageURL, &e.Rarity, &e.Difficulty); err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *ComposeRepository) GetUndiscoveredChallenges(ctx context.Context, gameCode string, exclude []int) ([]model.ChallengeItem, error) {
	if len(exclude) == 0 {
		exclude = []int{-1}
	}
	rows, err := r.DB.Query(ctx, `
		SELECT e.element_id, e.name, e.image_url
		FROM element e
		JOIN game g ON e.game_id = g.game_id
		WHERE g.code=$1 AND e.is_character=true AND e.element_id <> ALL($2)
	`, gameCode, exclude)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.ChallengeItem
	for rows.Next() {
		var c model.ChallengeItem
		if err := rows.Scan(&c.ID, &c.Name, &c.ImageURL); err != nil {
			return nil, err
		}
		list = append(list, c)
	}
	return list, nil
}
