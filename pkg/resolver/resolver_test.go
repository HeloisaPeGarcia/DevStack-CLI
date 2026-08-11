package resolver_test

import (
	"devstack/pkg/resolver"
	"testing"
)

func TestResolveStack_ByExactID(t *testing.T) {
	r := resolver.NewResolver()

	recipe, found := r.ResolveStack("go-react")
	if !found {
		t.Fatal("esperava encontrar a stack 'go-react' por ID exato, mas não encontrou")
	}
	if recipe.ID != "go-react" {
		t.Errorf("esperava ID 'go-react', obteve '%s'", recipe.ID)
	}
}

func TestResolveStack_ByFullName(t *testing.T) {
	r := resolver.NewResolver()

	recipe, found := r.ResolveStack("Go Backend + React Frontend")
	if !found {
		t.Fatal("esperava encontrar a stack pelo nome completo")
	}
	if recipe.ID != "go-react" {
		t.Errorf("esperava ID 'go-react', obteve '%s'", recipe.ID)
	}
}

func TestResolveStack_ByAlias(t *testing.T) {
	r := resolver.NewResolver()

	cases := []struct {
		alias  string
		wantID string
	}{
		{"fullstack-go", "go-react"},
		{"fastapi", "python-fastapi"},
		{"mevn", "node-vue"},
	}

	for _, tc := range cases {
		t.Run(tc.alias, func(t *testing.T) {
			recipe, found := r.ResolveStack(tc.alias)
			if !found {
				t.Fatalf("alias '%s': stack não encontrada", tc.alias)
			}
			if recipe.ID != tc.wantID {
				t.Errorf("alias '%s': esperava ID '%s', obteve '%s'", tc.alias, tc.wantID, recipe.ID)
			}
		})
	}
}

func TestResolveStack_ByKeyword_DeterministicPriority(t *testing.T) {
	r := resolver.NewResolver()

	cases := []struct {
		keyword string
		wantID  string
	}{
		{"Preciso de um backend em python", "python-fastapi"},
		{"quero uma stack com vue", "node-vue"},
		{"go e react para minha empresa", "go-react"},
		{"fastapi com python", "python-fastapi"},
	}

	for _, tc := range cases {
		t.Run(tc.keyword, func(t *testing.T) {
			recipe, found := r.ResolveStack(tc.keyword)
			if !found {
				t.Fatalf("keyword '%s': stack não encontrada", tc.keyword)
			}
			if recipe.ID != tc.wantID {
				t.Errorf("keyword '%s': esperava ID '%s', obteve '%s'", tc.keyword, tc.wantID, recipe.ID)
			}
		})
	}
}

func TestResolveStack_Fallback(t *testing.T) {
	r := resolver.NewResolver()

	recipe, found := r.ResolveStack("xyz123desconhecido")
	if found {
		t.Error("esperava found = false para entrada não correspondida (fallback)")
	}
	if recipe == nil {
		t.Fatal("recipe retornada no fallback não deve ser nil")
	}
}
