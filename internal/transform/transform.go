package transform

import (
	"errors"
	"log/slog"
	"math/rand"
	"time"

	"github.com/opplieam/bb-transform/.jetgen/postgres/public/model"
)

type Config struct {
	Version       string
	Shuffle       bool
	TrainRatio    uint8
	ValidateRatio uint8
	TestRatio     uint8
}

type CategoryStorer interface {
	OriginalCategory() (Category, error)
	MatchedCategory() ([]model.MatchCategory, error)
	CleanUp(version string) error
	InsertDataset(dataset []model.CategoryDataset) error
}

type CategoryDeepest struct {
	Name string
	Path string
}

type Category map[int32]CategoryDeepest

var (
	ErrOriginalCategory = errors.New("failed to get original category")
	ErrMatchedCategory  = errors.New("failed to matched category")
	ErrCleanUp          = errors.New("failed to clean up")
	ErrInsertDataset    = errors.New("failed to insert dataset")
)

type Transform struct {
	log      *slog.Logger
	catStore CategoryStorer
	config   Config
}

func NewTransform(cs CategoryStorer, cfg Config) *Transform {
	return &Transform{
		log:      slog.With("component", "transform"),
		catStore: cs,
		config:   cfg,
	}
}

func (t *Transform) GenerateDataset() error {
	oCat, err := t.catStore.OriginalCategory()
	if err != nil {
		return ErrOriginalCategory
	}
	t.log.Info("get all original category")

	mCat, err := t.catStore.MatchedCategory()
	if err != nil {
		return ErrMatchedCategory
	}
	t.log.Info("get all matched category")

	if t.config.Shuffle {
		t.log.Info("shuffle matched category")

		src := rand.NewSource(time.Now().UnixNano())
		rng := rand.New(src)
		rng.Shuffle(len(mCat), func(i, j int) {
			mCat[i], mCat[j] = mCat[j], mCat[i]
		})
	}

	// Calculate the number of samples for each label
	totalSamples := len(mCat)
	numTrain := int(float64(totalSamples) * (float64(t.config.TrainRatio) / 100))
	numValidate := int(float64(totalSamples) * (float64(t.config.ValidateRatio) / 100))

	var dataset []model.CategoryDataset
	var label string
	for i, v := range mCat {
		if i < numTrain {
			label = "train"
		} else if i < numTrain+numValidate {
			label = "validate"
		} else {
			label = "test"
		}

		dataset = append(dataset, model.CategoryDataset{
			L1In:        v.L1,
			L2In:        v.L2,
			L3In:        v.L3,
			L4In:        v.L4,
			L5In:        v.L5,
			L6In:        v.L6,
			L7In:        v.L7,
			L8In:        v.L8,
			FullPathOut: oCat[*v.MatchID].Path,
			NameOut:     oCat[*v.MatchID].Name,
			Version:     t.config.Version,
			Label:       label,
		})
	}

	if err = t.catStore.CleanUp(t.config.Version); err != nil {
		return ErrCleanUp
	}
	t.log.Info("cleaned up dataset", "version", t.config.Version)

	if err = t.catStore.InsertDataset(dataset); err != nil {
		return ErrInsertDataset
	}
	t.log.Info("inserted dataset", "version", t.config.Version)
	return nil
}
