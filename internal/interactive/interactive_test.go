package interactive

import (
	"errors"
	"testing"
)

// skipIfRace skips tests when running with race detector due to
// known race conditions in the readline library (third-party dependency, not our code).
// The readline library has internal races in its terminal handling that are beyond our control.
func skipIfRace(t *testing.T) {
	if isRaceEnabled {
		t.Skip("Skipping due to known race conditions in readline library (third-party)")
	}
}

func TestNew(t *testing.T) {
	skipIfRace(t)

	session, err := New()
	if err != nil {
		t.Fatalf("New() error = %v, want nil", err)
	}

	if session == nil {
		t.Fatal("New() returned nil session")
	}

	if session.rl == nil {
		t.Error("Session.rl is nil")
	}

	if session.commands == nil {
		t.Error("Session.commands map is nil")
	}

	// Clean up
	session.Close()
}

func TestRegisterCommand(t *testing.T) {
	skipIfRace(t)

	session, err := New()
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	defer session.Close()

	called := false
	handler := func(args []string) error {
		called = true
		return nil
	}

	session.RegisterCommand("test", handler)

	// Verify command was registered
	if session.commands["test"] == nil {
		t.Error("Command 'test' was not registered")
	}

	// Call the handler
	err = session.commands["test"]([]string{})
	if err != nil {
		t.Errorf("Handler returned error: %v", err)
	}

	if !called {
		t.Error("Handler was not called")
	}
}

func TestRegisterMultipleCommands(t *testing.T) {
	skipIfRace(t)

	session, err := New()
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	defer session.Close()

	commands := []string{"cmd1", "cmd2", "cmd3"}

	for _, cmd := range commands {
		session.RegisterCommand(cmd, func(args []string) error {
			return nil
		})
	}

	// Verify all commands were registered
	for _, cmd := range commands {
		if session.commands[cmd] == nil {
			t.Errorf("Command '%s' was not registered", cmd)
		}
	}

	// Verify count
	if len(session.commands) != len(commands) {
		t.Errorf("Got %d commands, want %d", len(session.commands), len(commands))
	}
}

func TestCommandHandlerWithArgs(t *testing.T) {
	skipIfRace(t)

	session, err := New()
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	defer session.Close()

	var receivedArgs []string
	handler := func(args []string) error {
		receivedArgs = args
		return nil
	}

	session.RegisterCommand("test", handler)

	// Call with arguments
	expectedArgs := []string{"arg1", "arg2", "arg3"}
	session.commands["test"](expectedArgs)

	if len(receivedArgs) != len(expectedArgs) {
		t.Errorf("Got %d args, want %d", len(receivedArgs), len(expectedArgs))
	}

	for i, arg := range expectedArgs {
		if receivedArgs[i] != arg {
			t.Errorf("Arg[%d] = %q, want %q", i, receivedArgs[i], arg)
		}
	}
}

func TestCommandHandlerWithError(t *testing.T) {
	skipIfRace(t)

	session, err := New()
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	defer session.Close()

	expectedErr := errors.New("command failed")
	handler := func(args []string) error {
		return expectedErr
	}

	session.RegisterCommand("failing", handler)

	// Call handler and verify error is returned
	err = session.commands["failing"]([]string{})
	if err != expectedErr {
		t.Errorf("Handler returned error = %v, want %v", err, expectedErr)
	}
}

func TestHandleBuiltIn_Exit(t *testing.T) {
	skipIfRace(t)

	session, err := New()
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	defer session.Close()

	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"exit command", "exit", true},
		{"quit command", "quit", true},
		{"help command", "help", true},
		{"? command", "?", true},
		{"clear command", "clear", true},
		{"cls command", "cls", true},
		{"unknown command", "unknown", false},
		{"regular command", "run", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := session.handleBuiltIn(tt.input)
			if got != tt.want {
				t.Errorf("handleBuiltIn(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestHandleBuiltIn_Help(t *testing.T) {
	skipIfRace(t)

	session, err := New()
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	defer session.Close()

	// help and ? should both work
	if !session.handleBuiltIn("help") {
		t.Error("handleBuiltIn('help') should return true")
	}

	if !session.handleBuiltIn("?") {
		t.Error("handleBuiltIn('?') should return true")
	}
}

func TestHandleBuiltIn_Clear(t *testing.T) {
	skipIfRace(t)

	session, err := New()
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	defer session.Close()

	// clear and cls should both work
	if !session.handleBuiltIn("clear") {
		t.Error("handleBuiltIn('clear') should return true")
	}

	if !session.handleBuiltIn("cls") {
		t.Error("handleBuiltIn('cls') should return true")
	}
}

func TestClose(t *testing.T) {
	skipIfRace(t)

	session, err := New()
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	err = session.Close()
	if err != nil {
		t.Errorf("Close() error = %v, want nil", err)
	}

	// Calling Close() again should not panic
	err = session.Close()
	// May return error on second close, but should not panic
}

func TestSession_Structure(t *testing.T) {
	skipIfRace(t)

	session, err := New()
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	defer session.Close()

	// Verify session structure
	if session.rl == nil {
		t.Error("Session.rl should not be nil")
	}

	if session.commands == nil {
		t.Error("Session.commands should not be nil")
	}

	// Commands map should be empty initially
	if len(session.commands) != 0 {
		t.Errorf("Initial commands map should be empty, got %d commands", len(session.commands))
	}
}

func TestRegisterCommand_Override(t *testing.T) {
	skipIfRace(t)

	session, err := New()
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	defer session.Close()

	firstCalled := false
	secondCalled := false

	// Register first handler
	session.RegisterCommand("test", func(args []string) error {
		firstCalled = true
		return nil
	})

	// Override with second handler
	session.RegisterCommand("test", func(args []string) error {
		secondCalled = true
		return nil
	})

	// Call the command
	session.commands["test"]([]string{})

	// Only second handler should be called
	if firstCalled {
		t.Error("First handler should not be called after override")
	}
	if !secondCalled {
		t.Error("Second handler should be called")
	}
}

func TestCommandHandler_NilArgs(t *testing.T) {
	skipIfRace(t)

	session, err := New()
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	defer session.Close()

	var receivedArgs []string
	session.RegisterCommand("test", func(args []string) error {
		receivedArgs = args
		return nil
	})

	// Call with nil args
	session.commands["test"](nil)

	// Should not panic, receivedArgs should be nil
	if receivedArgs != nil {
		t.Errorf("receivedArgs = %v, want nil", receivedArgs)
	}
}

func TestCommandHandler_EmptyArgs(t *testing.T) {
	skipIfRace(t)

	session, err := New()
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	defer session.Close()

	var receivedArgs []string
	session.RegisterCommand("test", func(args []string) error {
		receivedArgs = args
		return nil
	})

	// Call with empty args
	session.commands["test"]([]string{})

	// Should have empty slice
	if len(receivedArgs) != 0 {
		t.Errorf("receivedArgs length = %d, want 0", len(receivedArgs))
	}
}

func BenchmarkRegisterCommand(b *testing.B) {
	session, err := New()
	if err != nil {
		b.Fatalf("New() failed: %v", err)
	}
	defer session.Close()

	handler := func(args []string) error {
		return nil
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		session.RegisterCommand("test", handler)
	}
}

func BenchmarkHandleBuiltIn(b *testing.B) {
	session, err := New()
	if err != nil {
		b.Fatalf("New() failed: %v", err)
	}
	defer session.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		session.handleBuiltIn("unknown")
	}
}

// Test that built-in commands cover expected cases
func TestBuiltInCommands_Coverage(t *testing.T) {
	skipIfRace(t)

	session, err := New()
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	defer session.Close()

	builtInCommands := []string{
		"exit", "quit",  // Exit commands
		"help", "?",     // Help commands
		"clear", "cls",  // Clear commands
	}

	for _, cmd := range builtInCommands {
		if !session.handleBuiltIn(cmd) {
			t.Errorf("Command '%s' should be handled as built-in", cmd)
		}
	}
}

// Test that non-built-in commands return false
func TestHandleBuiltIn_NonBuiltIn(t *testing.T) {
	skipIfRace(t)

	session, err := New()
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	defer session.Close()

	nonBuiltInCommands := []string{
		"run", "list", "compare", "export",
		"stats", "trend", "check", "doctor",
		"flamegraph", "serve", "delete",
		"custom", "test", "benchmark",
	}

	for _, cmd := range nonBuiltInCommands {
		if session.handleBuiltIn(cmd) {
			t.Errorf("Command '%s' should not be handled as built-in", cmd)
		}
	}
}

// Test printWelcome doesn't panic
func TestPrintWelcome(t *testing.T) {
	skipIfRace(t)

	session, err := New()
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	defer session.Close()

	// Should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("printWelcome panicked: %v", r)
		}
	}()

	session.printWelcome()
}

// Test printGoodbye doesn't panic
func TestPrintGoodbye(t *testing.T) {
	skipIfRace(t)

	session, err := New()
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	defer session.Close()

	// Should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("printGoodbye panicked: %v", r)
		}
	}()

	session.printGoodbye()
}

// Test printHelp doesn't panic
func TestPrintHelp(t *testing.T) {
	skipIfRace(t)

	session, err := New()
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	defer session.Close()

	// Should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("printHelp panicked: %v", r)
		}
	}()

	session.printHelp()
}

// Test that CommandHandler is a function type
func TestCommandHandlerType(t *testing.T) {
	var handler CommandHandler = func(args []string) error {
		return nil
	}

	if handler == nil {
		t.Error("CommandHandler should not be nil")
	}

	err := handler([]string{"test"})
	if err != nil {
		t.Errorf("Handler returned error: %v", err)
	}
}

// Test multiple Close() calls
func TestMultipleClose(t *testing.T) {
	skipIfRace(t)

	session, err := New()
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	// Close multiple times
	for i := 0; i < 3; i++ {
		err := session.Close()
		// May return error after first close, but should not panic
		_ = err
	}
}

// Test registering many commands
func TestRegisterManyCommands(t *testing.T) {
	skipIfRace(t)

	session, err := New()
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	defer session.Close()

	numCommands := 100
	for i := 0; i < numCommands; i++ {
		cmdName := "cmd" + string(rune(i))
		session.RegisterCommand(cmdName, func(args []string) error {
			return nil
		})
	}

	if len(session.commands) != numCommands {
		t.Errorf("Got %d commands, want %d", len(session.commands), numCommands)
	}
}

// Test command with special characters in name
func TestRegisterCommand_SpecialChars(t *testing.T) {
	skipIfRace(t)

	session, err := New()
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	defer session.Close()

	specialNames := []string{
		"cmd-with-dash",
		"cmd_with_underscore",
		"cmd.with.dots",
		"cmd123",
		"CmdWithCaps",
	}

	for _, name := range specialNames {
		session.RegisterCommand(name, func(args []string) error {
			return nil
		})

		if session.commands[name] == nil {
			t.Errorf("Command '%s' was not registered", name)
		}
	}
}

// Test handler that modifies args
func TestCommandHandler_ModifiesArgs(t *testing.T) {
	skipIfRace(t)

	session, err := New()
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	defer session.Close()

	session.RegisterCommand("modify", func(args []string) error {
		if len(args) > 0 {
			args[0] = "modified"
		}
		return nil
	})

	args := []string{"original"}
	session.commands["modify"](args)

	// Args should be modified
	if args[0] != "modified" {
		t.Errorf("Args[0] = %q, want 'modified'", args[0])
	}
}

// Test that handleBuiltIn works with whitespace
func TestHandleBuiltIn_Whitespace(t *testing.T) {
	skipIfRace(t)

	session, err := New()
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	defer session.Close()

	// These should not be recognized as built-in (whitespace matters)
	tests := []string{
		" exit",
		"exit ",
		" help ",
		"HELP", // Case sensitive
		"Exit",
	}

	for _, test := range tests {
		if session.handleBuiltIn(test) {
			t.Errorf("handleBuiltIn(%q) should return false (whitespace/case matters)", test)
		}
	}
}

// Integration-like test
func TestSession_FullFlow(t *testing.T) {
	skipIfRace(t)

	session, err := New()
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	defer session.Close()

	// Register some commands
	runCalled := false
	listCalled := false

	session.RegisterCommand("run", func(args []string) error {
		runCalled = true
		return nil
	})

	session.RegisterCommand("list", func(args []string) error {
		listCalled = true
		return nil
	})

	// Execute commands
	session.commands["run"]([]string{})
	session.commands["list"]([]string{})

	// Verify both were called
	if !runCalled {
		t.Error("run command was not called")
	}
	if !listCalled {
		t.Error("list command was not called")
	}

	// Test built-in
	if !session.handleBuiltIn("help") {
		t.Error("help should be handled as built-in")
	}

	// Close
	if err := session.Close(); err != nil {
		t.Errorf("Close() error = %v", err)
	}
}

func TestHandleBuiltIn_CaseInsensitive(t *testing.T) {
	skipIfRace(t)

	session, err := New()
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	defer session.Close()

	// Built-in commands are case sensitive (lowercase only)
	casedCommands := []string{
		"EXIT", "Exit", "eXit",
		"QUIT", "Quit",
		"HELP", "Help",
		"CLEAR", "Clear",
	}

	for _, cmd := range casedCommands {
		if session.handleBuiltIn(cmd) {
			t.Errorf("handleBuiltIn(%q) should be case-sensitive and return false", cmd)
		}
	}
}

// Test error propagation
func TestCommandHandler_ErrorPropagation(t *testing.T) {
	skipIfRace(t)

	session, err := New()
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	defer session.Close()

	testErr := errors.New("specific error")

	session.RegisterCommand("error", func(args []string) error {
		return testErr
	})

	err = session.commands["error"]([]string{})
	if err != testErr {
		t.Errorf("Error = %v, want %v", err, testErr)
	}

	// Verify error message
	if err.Error() != "specific error" {
		t.Errorf("Error message = %q, want 'specific error'", err.Error())
	}
}
