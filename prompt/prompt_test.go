package prompt

import (
	"errors"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func newPromptWithHandlers(t *testing.T, handlers ...func(tea.Model)) *Prompt {
	t.Helper()

	var call int

	return New(WithProgramRunner(func(model tea.Model) (tea.Model, error) {
		if call < len(handlers) {
			handlers[call](model)
		}
		call++

		markModelComplete(model)
		return model, nil
	}))
}

func markModelComplete(model tea.Model) {
	switch m := model.(type) {
	case *inputModel:
		m.done = true
	case *confirmModel:
		m.done = true
	case *selectModel:
		m.done = true
	case *multiSelectModel:
		m.done = true
	case *passwordModel:
		m.done = true
	}
}

func TestNew(t *testing.T) {
	p := New()
	if p == nil {
		t.Fatal("expected prompt to be created")
	}
	if p.terminal == nil {
		t.Fatal("expected terminal to be initialised")
	}
}

func TestAskString(t *testing.T) {
	p := newPromptWithHandlers(t, func(model tea.Model) {
		m := model.(*inputModel)
		m.input = "test input"
	})

	result, err := p.AskString("Enter something")
	if err != nil {
		t.Fatalf("AskString returned error: %v", err)
	}
	if result != "test input" {
		t.Fatalf("expected 'test input', got %q", result)
	}
}

func TestAskStringWithDefault(t *testing.T) {
	p := newPromptWithHandlers(t,
		func(model tea.Model) {
			// Simulate pressing enter without typing anything
			_ = model.(*inputModel)
		},
		func(model tea.Model) {
			m := model.(*inputModel)
			m.input = "custom input"
		},
	)

	first, err := p.AskStringWithDefault("Enter name", "default")
	if err != nil {
		t.Fatalf("AskStringWithDefault returned error: %v", err)
	}
	if first != "default" {
		t.Fatalf("expected 'default', got %q", first)
	}

	second, err := p.AskStringWithDefault("Enter name", "default")
	if err != nil {
		t.Fatalf("AskStringWithDefault returned error: %v", err)
	}
	if second != "custom input" {
		t.Fatalf("expected 'custom input', got %q", second)
	}
}

func TestAskInt(t *testing.T) {
	p := newPromptWithHandlers(t,
		func(model tea.Model) {
			m := model.(*inputModel)
			m.input = "42"
		},
		func(model tea.Model) {
			m := model.(*inputModel)
			m.input = "not a number"
		},
	)

	value, err := p.AskInt("Number?")
	if err != nil {
		t.Fatalf("AskInt returned error: %v", err)
	}
	if value != 42 {
		t.Fatalf("expected 42, got %d", value)
	}

	_, err = p.AskInt("Number?")
	if err == nil {
		t.Fatal("expected error for invalid integer")
	}
	var inputErr *InputError
	if !errors.As(err, &inputErr) {
		t.Fatalf("expected InputError, got %T", err)
	}
}

func TestAskIntWithDefault(t *testing.T) {
	p := newPromptWithHandlers(t,
		func(model tea.Model) {
			// simulate enter without input
			_ = model.(*inputModel)
		},
		func(model tea.Model) {
			m := model.(*inputModel)
			m.input = "25"
		},
	)

	value, err := p.AskIntWithDefault("Number?", 10)
	if err != nil {
		t.Fatalf("AskIntWithDefault returned error: %v", err)
	}
	if value != 10 {
		t.Fatalf("expected 10, got %d", value)
	}

	value, err = p.AskIntWithDefault("Number?", 10)
	if err != nil {
		t.Fatalf("AskIntWithDefault returned error: %v", err)
	}
	if value != 25 {
		t.Fatalf("expected 25, got %d", value)
	}
}

func TestAskBool(t *testing.T) {
	yesPrompt := newPromptWithHandlers(t, func(model tea.Model) {
		m := model.(*confirmModel)
		m.choice = 1
	})
	result, err := yesPrompt.AskBool("Continue?")
	if err != nil {
		t.Fatalf("AskBool returned error: %v", err)
	}
	if !result {
		t.Fatal("expected true for affirmative response")
	}

	noPrompt := newPromptWithHandlers(t, func(model tea.Model) {
		m := model.(*confirmModel)
		m.choice = 0
	})
	result, err = noPrompt.AskBool("Continue?")
	if err != nil {
		t.Fatalf("AskBool returned error: %v", err)
	}
	if result {
		t.Fatal("expected false for negative response")
	}
}

func TestAskBoolWithDefault(t *testing.T) {
	truePrompt := newPromptWithHandlers(t, func(model tea.Model) {
		// leave choice unset to use default
		_ = model.(*confirmModel)
	})
	result, err := truePrompt.AskBoolWithDefault("Continue?", true)
	if err != nil {
		t.Fatalf("AskBoolWithDefault returned error: %v", err)
	}
	if !result {
		t.Fatal("expected true when default is true")
	}

	falsePrompt := newPromptWithHandlers(t, func(model tea.Model) {
		_ = model.(*confirmModel)
	})
	result, err = falsePrompt.AskBoolWithDefault("Continue?", false)
	if err != nil {
		t.Fatalf("AskBoolWithDefault returned error: %v", err)
	}
	if result {
		t.Fatal("expected false when default is false")
	}
}

func TestAskSelect(t *testing.T) {
	p := newPromptWithHandlers(t, func(model tea.Model) {
		m := model.(*selectModel)
		m.choice = 2
	})

	result, err := p.AskSelect("Pick one", []string{"a", "b", "c"})
	if err != nil {
		t.Fatalf("AskSelect returned error: %v", err)
	}
	if result != "c" {
		t.Fatalf("expected 'c', got %q", result)
	}
}

func TestAskSelect_NoOptions(t *testing.T) {
	p := New()
	_, err := p.AskSelect("Pick one", nil)
	if err == nil {
		t.Fatal("expected error when no options provided")
	}
}

func TestAskMultiSelect(t *testing.T) {
	p := newPromptWithHandlers(t, func(model tea.Model) {
		m := model.(*multiSelectModel)
		m.choices[0] = true
		m.choices[2] = true
	})

	result, err := p.AskMultiSelect("Pick any", []string{"a", "b", "c"})
	if err != nil {
		t.Fatalf("AskMultiSelect returned error: %v", err)
	}
	if len(result) != 2 || result[0] != "a" || result[1] != "c" {
		t.Fatalf("expected selections [a c], got %v", result)
	}
}

func TestAskPassword(t *testing.T) {
	p := newPromptWithHandlers(t, func(model tea.Model) {
		m := model.(*passwordModel)
		m.input = "secret"
	})

	result, err := p.AskPassword("Enter password")
	if err != nil {
		t.Fatalf("AskPassword returned error: %v", err)
	}
	if result != "secret" {
		t.Fatalf("expected 'secret', got %q", result)
	}
}
