package helpers

import (
	"fmt"

	"github.com/MadManJJ/cms-api/dto"
)

func CategoryTypeResponseEqual(a, b dto.CategoryTypeResponse) bool {
	// Compare IDs, TypeCode, IsActive, CreatedAt, UpdatedAt directly
	if a.ID != b.ID ||
		a.TypeCode != b.TypeCode ||
		a.IsActive != b.IsActive ||
		!a.CreatedAt.Equal(b.CreatedAt) ||
		!a.UpdatedAt.Equal(b.UpdatedAt) {
		fmt.Println("ID, TypeCode, IsActive, CreatedAt, or UpdatedAt mismatch")
		return false
	}

	// Compare Name (string pointer)
	if (a.Name == nil) != (b.Name == nil) {
		fmt.Println("Name pointer mismatch")
		return false
	}
	if a.Name != nil && *a.Name != *b.Name {
		fmt.Println("Name value mismatch:", *a.Name, "!=", *b.Name)
		return false
	}

	// Compare ChildrenCount (map pointer)
	if (a.ChildrenCount == nil) != (b.ChildrenCount == nil) {
		fmt.Println("ChildrenCount pointer mismatch")
		return false
	}
	if a.ChildrenCount != nil {
		if len(*a.ChildrenCount) != len(*b.ChildrenCount) {
			fmt.Println("ChildrenCount length mismatch")
			return false
		}
		for k, v := range *a.ChildrenCount {
			if (*b.ChildrenCount)[k] != v {
				fmt.Println("ChildrenCount value mismatch for key", k, ":", (*b.ChildrenCount)[k], "!=", v)
				return false
			}
		}
	}
	fmt.Println("CategoryTypeResponse objects are equal")
	return true
}
