// TODO: Reactive UI - update data on the screen when new data is present
// TODO: Fix cyrillic symbols
// TODO: Upgrade backend and fronted: add users models, make messages model better

use futures_util::{SinkExt, TryStreamExt};
use reqwest_websocket::{Error, Message, Upgrade};
use sola_raylib::{color::Color, drawing::RaylibDraw};
use std::sync::Arc;
use tokio::sync::Mutex;

type ChatMessage = String;

struct AppState {
    messages_history: Vec<ChatMessage>,
}

#[tokio::main]
async fn main() {
    let app_state = Arc::new(Mutex::new(AppState {
        messages_history: Vec::new(),
    }));

    let client_arc = app_state.clone();
    tokio::spawn(async move {
        run_client_thread(client_arc).await;
    });

    let raylib_arc = app_state.clone();
    run_raylib_thread(raylib_arc).await;
}

async fn run_client_thread(state: Arc<Mutex<AppState>>) {
    let result = start_client(state).await;
    match result {
        Ok(()) => println!("Great!"),
        Err(err) => println!("Not great: {}.", err),
    }
}

async fn run_raylib_thread(state: Arc<Mutex<AppState>>) {
    let (mut rl, thread) = sola_raylib::init()
        .size(500, 500)
        .title("Global Chat")
        .build();

    while !rl.window_should_close() {
        let messages_history = {
            let state = state.lock().await;
            state.messages_history.clone()
        };

        let mut dr = rl.begin_drawing(&thread);

        dr.clear_background(Color::WHITE);

        let mut y_pos = 20;
        for msg in messages_history.clone() {
            dr.draw_text(msg.as_str(), 12, y_pos, 20, Color::BLACK);
            y_pos += 20;
        }
    }
}

async fn start_client(state: Arc<Mutex<AppState>>) -> Result<(), String> {
    let client = reqwest::Client::default();

    let response = client
        .get("http://localhost:8080/chat")
        .upgrade()
        .send()
        .await;

    let upgrade_response = response.or_else(e2s)?;
    let mut conn = Box::new(upgrade_response.into_websocket().await.or_else(e2s)?);

    while let Some(msg) = conn.try_next().await.or_else(e2s)? {
        match msg {
            Message::Text(text) => {
                if text == ";;end;;" {
                    continue;
                }

                println!("{}", text);
                let mut state = state.lock().await;
                state.messages_history.push(text);
            }
            Message::Binary(_) => println!("We've got some bytes!!!"),
            Message::Ping(_) => println!("Someone pinged me..."),
            Message::Pong(_) => println!("Someone ponged me..."),
            Message::Close { code: _, reason: _ } => conn.close().await.or_else(e2s)?,
        }
    }

    return Ok(());
}

fn e2s<T>(err: Error) -> Result<T, String> {
    return Err(err.to_string());
}
