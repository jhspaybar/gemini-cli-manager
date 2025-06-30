use clap::Parser;
use cli::Cli;
use color_eyre::Result;

use crate::app::App;

mod action;
mod app;
mod cli;
mod components;
mod config;
mod errors;
mod launcher;
mod logging;
mod models;
mod storage;
mod tui;
mod view;

#[tokio::main]
async fn main() -> Result<()> {
    crate::errors::init()?;
    crate::logging::init()?;

    let args = Cli::parse();
    
    // Handle list-storage flag
    if args.list_storage {
        list_storage_contents()?;
        return Ok(());
    }
    
    let mut app = App::new(args.tick_rate, args.frame_rate)?;
    app.run().await?;
    Ok(())
}

fn list_storage_contents() -> Result<()> {
    use crate::storage::Storage;
    
    println!("Gemini CLI Manager - Storage Contents");
    println!("=====================================\n");
    
    let storage = Storage::new()?;
    
    println!("Extensions:");
    println!("-----------");
    let extensions = storage.list_extensions()?;
    for ext in extensions {
        println!("- {} v{}: {}", ext.name, ext.version, ext.description.as_deref().unwrap_or(""));
    }
    
    println!("\nProfiles:");
    println!("---------");
    let profiles = storage.list_profiles()?;
    for profile in profiles {
        println!("- {} ({}): {} extensions", 
            profile.name, 
            profile.id,
            profile.extension_ids.len()
        );
        if let Some(desc) = &profile.description {
            println!("  Description: {}", desc);
        }
        if !profile.metadata.tags.is_empty() {
            println!("  Tags: {}", profile.metadata.tags.join(", "));
        }
    }
    
    Ok(())
}
