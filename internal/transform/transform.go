// Package transform provides functionalities to transform raw category data into a structured dataset suitable for machine learning.
// It handles shuffling, splitting into train/validate/test sets, and database interactions.
// It utilizes a CategoryStorer interface to abstract the underlying data storage mechanism.
// The package prepares data specifically for training, validating, and testing machine learning models,
// enabling the development of models that can predict or classify categories based on input features.
package transform

import (
	"log/slog"
	"math/rand"
	"time"

	"github.com/opplieam/bb-transform/.jetgen/postgres/public/model"
)

// Config holds the configuration parameters for the transformation process.
// It includes settings for the dataset version, whether to shuffle the data for randomness,
// and the ratios for splitting the data into train, validate, and test sets, which are crucial for model training and evaluation.
type Config struct {
	Version       string `json:"version"`
	Shuffle       bool   `json:"shuffle"`
	TrainRatio    uint8  `json:"train_ratio"`
	ValidateRatio uint8  `json:"validate_ratio"`
	TestRatio     uint8  `json:"test_ratio"`
}

// CategoryStorer defines the interface for storing and retrieving category data used in the machine learning pipeline.
// It provides methods for accessing original and matched categories,
// cleaning up old datasets based on version, and inserting new datasets, ensuring data consistency and version control.
type CategoryStorer interface {
	OriginalCategory() (Category, error)
	MatchedCategory() ([]model.MatchCategory, error)
	CleanUp(version string) error
	InsertDataset(dataset []model.CategoryDataset) error
}

// CategoryDeepest represents the deepest level of a category, containing its name and its full hierarchical path.
// This information can be used as labels or features in machine learning models.
type CategoryDeepest struct {
	Name string
	Path string
}

// Category represents a mapping of category IDs to their corresponding CategoryDeepest information.
// This mapping is essential for translating between raw data and the structured information used in the dataset.
type Category map[int32]CategoryDeepest

type Transform struct {
	log      *slog.Logger
	catStore CategoryStorer
	config   Config
}

// NewTransform creates a new Transform instance, responsible for orchestrating the dataset generation process.
// It initializes the logger with a "transform" component tag for tracking,
// sets the CategoryStorer for data access, and configures the transformation process using the provided Config.
func NewTransform(cs CategoryStorer, cfg Config) *Transform {
	return &Transform{
		log:      slog.With("component", "transform"),
		catStore: cs,
		config:   cfg,
	}
}

// GenerateDataset generates a dataset specifically designed for training and evaluating machine learning models.
// It retrieves original and matched categories from the CategoryStorer,
// optionally shuffles the matched categories to ensure a random distribution of data,
// splits the data into train, validate, and test sets according to the configured ratios, which is crucial for model training and performance assessment,
// cleans up previous datasets with the same version from the CategoryStorer to avoid data conflicts,
// and inserts the newly generated dataset into the CategoryStorer.
// The resulting dataset contains features (L1-L8) and labels (FullPathOut, NameOut) that can be used to train a model to predict category paths or names.
// Returns an error if any of the steps fail, using specific error variables for clarity.
func (t *Transform) GenerateDataset() error {
	oCat, err := t.catStore.OriginalCategory()
	if err != nil {
		return err
	}
	t.log.Info("get all original category")

	mCat, err := t.catStore.MatchedCategory()
	if err != nil {
		return err
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
		return err
	}
	t.log.Info("cleaned up dataset", "version", t.config.Version)

	if err = t.catStore.InsertDataset(dataset); err != nil {
		return err
	}
	t.log.Info("inserted dataset", "version", t.config.Version)
	return nil
}
