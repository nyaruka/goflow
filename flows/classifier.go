package flows

import "github.com/nyaruka/goflow/assets"

// Classifier represents an NLU classifier.
type Classifier struct {
	assets.Classifier
}

// NewClassifier returns a new classifier object from the given classifier asset
func NewClassifier(asset assets.Classifier) *Classifier {
	return &Classifier{Classifier: asset}
}

// Asset returns the underlying asset
func (c *Classifier) Asset() assets.Classifier { return c.Classifier }

// Reference returns a reference to this classifier
func (c *Classifier) Reference() *assets.ClassifierReference {
	return assets.NewClassifierReference(c.UUID(), c.Name())
}

// ClassifierAssets provides access to all classifier assets
type ClassifierAssets struct {
	byUUID map[assets.ClassifierUUID]*Classifier
}

// NewClassifierAssets creates a new set of classifier assets
func NewClassifierAssets(classifiers []assets.Classifier) *ClassifierAssets {
	s := &ClassifierAssets{
		byUUID: make(map[assets.ClassifierUUID]*Classifier, len(classifiers)),
	}
	for _, asset := range classifiers {
		s.byUUID[asset.UUID()] = NewClassifier(asset)
	}
	return s
}

// Get returns the classifier with the given UUID
func (s *ClassifierAssets) Get(uuid assets.ClassifierUUID) *Classifier {
	return s.byUUID[uuid]
}
