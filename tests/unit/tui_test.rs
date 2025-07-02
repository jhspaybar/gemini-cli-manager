#[cfg(test)]
mod tests {
    use gemini_cli_manager::tui::{Event, Tui};
    use ratatui::layout::Rect;
    use std::io::IsTerminal;

    // Skip these tests in CI since they require a real TTY
    #[tokio::test]
    async fn test_tui_creation() {
        // Check if we're in CI by checking for TTY
        if std::env::var("CI").is_ok() || !std::io::stdout().is_terminal() {
            // Skip test in CI
            return;
        }

        // Test creating a Tui instance
        let tui = Tui::new();
        assert!(tui.is_ok(), "Tui should be created successfully");
    }

    #[tokio::test]
    async fn test_tui_with_tick_rate() {
        // Check if we're in CI by checking for TTY
        if std::env::var("CI").is_ok() || !std::io::stdout().is_terminal() {
            // Skip test in CI
            return;
        }

        let tui = Tui::new().unwrap().tick_rate(120.0);

        // Just verify we can set tick rate without errors
        drop(tui);
    }

    #[tokio::test]
    async fn test_tui_with_frame_rate() {
        // Check if we're in CI by checking for TTY
        if std::env::var("CI").is_ok() || !std::io::stdout().is_terminal() {
            // Skip test in CI
            return;
        }

        let tui = Tui::new().unwrap().frame_rate(60.0);

        // Just verify we can set frame rate without errors
        drop(tui);
    }

    #[tokio::test]
    async fn test_tui_builder_pattern() {
        // Check if we're in CI by checking for TTY
        if std::env::var("CI").is_ok() || !std::io::stdout().is_terminal() {
            // Skip test in CI
            return;
        }

        let tui = Tui::new()
            .unwrap()
            .tick_rate(120.0)
            .frame_rate(60.0)
            .mouse(true);

        // Just verify we can chain builder methods
        drop(tui);
    }

    #[test]
    fn test_event_creation() {
        use crossterm::event::{KeyCode, KeyEvent, KeyEventKind, KeyEventState, KeyModifiers};

        // Test creating key event
        let key_event = KeyEvent {
            code: KeyCode::Char('q'),
            modifiers: KeyModifiers::NONE,
            kind: KeyEventKind::Press,
            state: KeyEventState::NONE,
        };
        let event = Event::Key(key_event);

        match event {
            Event::Key(k) => assert_eq!(k.code, KeyCode::Char('q')),
            _ => panic!("Expected key event"),
        }
    }

    #[test]
    fn test_mouse_event_creation() {
        use crossterm::event::{MouseButton, MouseEvent, MouseEventKind};

        let mouse_event = MouseEvent {
            kind: MouseEventKind::Down(MouseButton::Left),
            column: 10,
            row: 20,
            modifiers: crossterm::event::KeyModifiers::NONE,
        };
        let event = Event::Mouse(mouse_event);

        match event {
            Event::Mouse(m) => {
                assert_eq!(m.column, 10);
                assert_eq!(m.row, 20);
            }
            _ => panic!("Expected mouse event"),
        }
    }

    #[test]
    fn test_tick_event() {
        let event = Event::Tick;
        match event {
            Event::Tick => { /* ok */ }
            _ => panic!("Expected tick event"),
        }
    }

    #[test]
    fn test_resize_event() {
        let event = Event::Resize(80, 24);
        match event {
            Event::Resize(w, h) => {
                assert_eq!(w, 80);
                assert_eq!(h, 24);
            }
            _ => panic!("Expected resize event"),
        }
    }

    #[tokio::test]
    async fn test_tui_resize() {
        // Check if we're in CI by checking for TTY
        if std::env::var("CI").is_ok() || !std::io::stdout().is_terminal() {
            // Skip test in CI
            return;
        }

        let mut tui = Tui::new().unwrap();
        let rect = Rect::new(0, 0, 100, 50);

        // Test resize method
        let result = tui.resize(rect);
        assert!(result.is_ok(), "Resize should succeed");
    }

    #[test]
    fn test_all_event_variants() {
        // Test that we can create all event variants
        let _events = vec![
            Event::Init,
            Event::Quit,
            Event::Error,
            Event::Closed,
            Event::Tick,
            Event::Render,
            Event::FocusGained,
            Event::FocusLost,
            Event::Paste("test".to_string()),
            Event::Resize(80, 24),
        ];

        // Just verify they can be created
    }

    #[test]
    fn test_event_pattern_matching() {
        // Test that we can match on all event types
        let events = vec![
            Event::Init,
            Event::Quit,
            Event::Error,
            Event::Closed,
            Event::Tick,
            Event::Render,
            Event::FocusGained,
            Event::FocusLost,
            Event::Paste("test".to_string()),
            Event::Key(crossterm::event::KeyEvent {
                code: crossterm::event::KeyCode::Enter,
                modifiers: crossterm::event::KeyModifiers::NONE,
                kind: crossterm::event::KeyEventKind::Press,
                state: crossterm::event::KeyEventState::NONE,
            }),
            Event::Mouse(crossterm::event::MouseEvent {
                kind: crossterm::event::MouseEventKind::Moved,
                column: 0,
                row: 0,
                modifiers: crossterm::event::KeyModifiers::NONE,
            }),
            Event::Resize(80, 24),
        ];

        // Verify we can match on each event type
        for event in events {
            match event {
                Event::Init => {}
                Event::Quit => {}
                Event::Error => {}
                Event::Closed => {}
                Event::Tick => {}
                Event::Render => {}
                Event::FocusGained => {}
                Event::FocusLost => {}
                Event::Paste(_) => {}
                Event::Key(_) => {}
                Event::Mouse(_) => {}
                Event::Resize(_, _) => {}
            }
        }
    }
}
