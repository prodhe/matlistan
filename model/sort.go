package model

// RecipesByTitle sorts recipes by title
type RecipesByTitle []Recipe

func (v RecipesByTitle) Len() int           { return len(v) }
func (v RecipesByTitle) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }
func (v RecipesByTitle) Less(i, j int) bool { return v[i].Title < v[j].Title }

// RecipesByCategory sorts recipes by primary category
type RecipesByCategory []Recipe

func (v RecipesByCategory) Len() int      { return len(v) }
func (v RecipesByCategory) Swap(i, j int) { v[i], v[j] = v[j], v[i] }
func (v RecipesByCategory) Less(i, j int) bool {
	if len(v[i].Categories) == 0 && len(v[j].Categories) > 0 {
		return true
	}
	if len(v[j].Categories) == 0 {
		return false
	}
	return v[i].Categories[0] < v[j].Categories[0]
}
