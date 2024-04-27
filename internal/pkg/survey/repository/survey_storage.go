package repository

import "github.com/jackc/pgx/v5/pgxpool"

type SurveyStorage struct {
	pool *pgxpool.Pool
}

func NewSurveyStorage(pool *pgxpool.Pool) *SurveyStorage {
	return &SurveyStorage{
		pool: pool,
	}
}

func (survey *SurveyStorage) SaveSurveyResults() {

}

func (survey *SurveyStorage) GetResults() {

}

func (survey *SurveyStorage) GetStatics() {

}
