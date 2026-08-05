package resolver

import (
	"devstack/pkg/config"
	"strings"
)

type keywordRule struct {
	keyword  string
	targetID string
	priority int
}

var keywordRules = []keywordRule{
	{keyword: "fastapi", targetID: "python-fastapi", priority: 10},
	{keyword: "python", targetID: "python-fastapi", priority: 5},
	{keyword: "react", targetID: "go-react", priority: 10},
	{keyword: "go", targetID: "go-react", priority: 5},
	{keyword: "vue", targetID: "node-vue", priority: 10},
	{keyword: "express", targetID: "node-vue", priority: 8},
	{keyword: "node", targetID: "node-vue", priority: 5},
}

type Resolver struct {
	recipes []config.StackRecipe
}

func NewResolver() *Resolver {
	return &Resolver{
		recipes: config.GetPredefinedRecipes(),
	}
}

func (r *Resolver) ResolveStack(input string) (*config.StackRecipe, bool) {
	cleanInput := strings.ToLower(strings.TrimSpace(input))

	for i := range r.recipes {
		if strings.ToLower(r.recipes[i].ID) == cleanInput ||
			strings.ToLower(r.recipes[i].Name) == cleanInput {
			return &r.recipes[i], true
		}
		for _, alias := range r.recipes[i].Aliases {
			if strings.ToLower(alias) == cleanInput {
				return &r.recipes[i], true
			}
		}
	}

	var bestRule *keywordRule
	for idx := range keywordRules {
		rule := &keywordRules[idx]
		if strings.Contains(cleanInput, rule.keyword) {
			if bestRule == nil || rule.priority > bestRule.priority {
				bestRule = rule
			}
		}
	}

	if bestRule != nil {
		if recipe := r.findByID(bestRule.targetID); recipe != nil {
			return recipe, true
		}
	}

	if len(r.recipes) > 0 {
		return &r.recipes[0], true
	}

	return nil, false
}

func (r *Resolver) findByID(id string) *config.StackRecipe {
	for i := range r.recipes {
		if r.recipes[i].ID == id {
			return &r.recipes[i]
		}
	}
	return nil
}
