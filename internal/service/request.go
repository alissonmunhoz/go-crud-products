package service

import "fmt"

func errParamIsRequired(name_, typ string) error {
	return fmt.Errorf("param: %s (type: %s) is required", name_, typ)
}

type CreateProductRequest struct {
	Name        string `json:"name" binding:"required"`
	Price       int64  `json:"price" binding:"required"`
	Quantity    int32  `json:"quantity" binding:"required"`
	Description string `json:"description" binding:"required"`
}

func (r *CreateProductRequest) Validate() error {
	if r.Name == "" && r.Price <= 0 && r.Quantity <= 0 && r.Description == "" {
		return fmt.Errorf("request body is empty or malformed")
	}

	if r.Name == "" {
		return errParamIsRequired("name", "string")
	}

	if r.Price <= 0 {
		return errParamIsRequired("price", "number")
	}

	if r.Quantity <= 0 {
		return errParamIsRequired("quantity", "number")
	}

	if r.Description == "" {
		return errParamIsRequired("description", "string")
	}

	return nil
}

type UpdateProductRequest struct {
	Name        string `json:"name"`
	Price       int64  `json:"price"`
	Quantity    int32  `json:"quantity"`
	Description string `json:"description"`
}

func (r *UpdateProductRequest) Validate() error {

	if r.Name != "" || r.Price > 0 || r.Quantity >= 0 || r.Description != "" {
		return nil
	}

	return fmt.Errorf("at least one valid field must be provided")
}
