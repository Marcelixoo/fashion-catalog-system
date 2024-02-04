package catalog

import "time"

// e.g. "shoes", "shirt", "dress"
type ProductType string

// e.g. "Women > Dress > Mini Dress"
type ProductCategory string

type Product struct {
	SKU        string
	Price      float64
	Type       ProductType
	Attributes []Attribute
	Categories []ProductCategory
	ModifiedAt time.Time
}

type Schema struct {
	ID                 string
	Name               string
	ProductType        ProductType
	RequiredAttributes []Attribute
}

type Attribute struct {
	Name         string
	DefaultValue any
	Translations []AttributeTranslation
	ModifiedAt   time.Time
}

type AttributeTranslation struct {
	Language string
	Value    string
}

/*

Use cases

- Products can be added/changed/removed.
- Mandatory attributes can be added/changed/removed from products.
- Consumers can fetch products.
*/

type CatalogManagementService interface {
	AddProduct(Product) Product
	UpdateProduct(Product) Product
	RemoveProduct(Product)
	UpdateRequiredAttributes(ProductType, []Attribute) []Attribute
	AddRequiredAttribute(ProductType, Attribute) []Attribute
	RemoveRequiredAttribute(ProductType, Attribute) []Attribute
}

type CatalogListingService interface {
	ListAllProducts() []*Product
}
