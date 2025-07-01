use std::collections::HashMap;
use std::env;
use std::fs;
use std::io::Write;
use std::path::{Path, PathBuf};
use std::process::{Command, Stdio};

use color_eyre::{eyre::eyre, Result};
use serde_json::json;

use crate::{models::{Extension, Profile}, storage::Storage};

pub struct Launcher {
    pub storage: Storage,
}

impl Launcher {
    #[allow(dead_code)]
    pub fn new() -> Self {
        Self { 
            storage: Storage::default(),
        }
    }
    
    pub fn with_storage(storage: Storage) -> Self {
        Self {
            storage,
        }
    }
    
    /// Launch Gemini CLI with the specified profile
    pub fn launch_with_profile(&self, profile: &Profile) -> Result<()> {
        // 1. Determine working directory
        let working_dir = if let Some(dir) = &profile.working_directory {
            // Expand ~ to home directory
            let expanded = if dir.starts_with("~") {
                dirs::home_dir()
                    .map(|home| home.join(&dir[2..]))
                    .unwrap_or_else(|| PathBuf::from(dir))
            } else {
                PathBuf::from(dir)
            };
            
            // Create directory if it doesn't exist
            if !expanded.exists() {
                fs::create_dir_all(&expanded)?;
            }
            
            expanded
        } else {
            env::current_dir()?
        };
        
        // 2. Clean existing configuration if requested
        if profile.launch_config.clean_launch {
            println!("🧹 Cleaning existing configuration...");
            self.clean_gemini_directory(&working_dir)?;
        }
        
        // 3. Set up workspace in the working directory
        self.setup_workspace(&working_dir)?;
        
        // 4. Install extensions to the working directory
        self.install_extensions_for_profile(profile, &working_dir)?;
        
        // 5. Set up environment
        let env_vars = self.prepare_environment(profile);
        
        // 6. Launch Gemini CLI
        println!("🚀 Launching Gemini CLI with profile: {}", profile.display_name());
        println!("📂 Working directory: {}", working_dir.display());
        println!("🔧 Extensions: {}", profile.extension_ids.join(", "));
        if profile.launch_config.cleanup_on_exit {
            println!("🧹 Will clean up extensions after exit");
        }
        println!();
        
        // Check if gemini is available (cross-platform)
        let gemini_check = if cfg!(target_os = "windows") {
            Command::new("where")
                .arg("gemini")
                .output()
                .map(|output| output.status.success())
                .unwrap_or(false)
        } else {
            Command::new("which")
                .arg("gemini")
                .output()
                .map(|output| output.status.success())
                .unwrap_or(false)
        };
        
        if !gemini_check {
            return Err(eyre!(
                "Gemini CLI not found. Please ensure 'gemini' is installed and in your PATH."
            ));
        }
        
        // Run gemini
        let mut cmd = Command::new("gemini");
        cmd.current_dir(&working_dir)
            .envs(&env_vars)
            .stdin(Stdio::inherit())
            .stdout(Stdio::inherit())
            .stderr(Stdio::inherit());
        
        let status = cmd.status()?;
        
        // 7. Clean up if requested
        if profile.launch_config.cleanup_on_exit {
            println!("\n🧹 Cleaning up extensions...");
            self.cleanup_extensions(&working_dir)?;
        }
        
        if !status.success() {
            return Err(eyre!("Gemini CLI exited with status: {}", status));
        }
        
        Ok(())
    }
    
    /// Set up the .gemini directory structure in the working directory
    pub fn setup_workspace(&self, working_dir: &Path) -> Result<()> {
        // Create .gemini directory structure
        let gemini_dir = working_dir.join(".gemini");
        fs::create_dir_all(&gemini_dir)?;
        
        let extensions_dir = gemini_dir.join("extensions");
        fs::create_dir_all(&extensions_dir)?;
        
        Ok(())
    }
    
    /// Install extensions to the working directory
    pub fn install_extensions_for_profile(&self, profile: &Profile, working_dir: &Path) -> Result<()> {
        let extensions_dir = working_dir.join(".gemini").join("extensions");
        
        // Load extensions from storage
        for ext_id in &profile.extension_ids {
            match self.storage.load_extension(ext_id) {
                Ok(extension) => {
                    self.install_extension(&extension, &extensions_dir)?;
                }
                Err(e) => {
                    eprintln!("Warning: Failed to load extension '{}': {}", ext_id, e);
                }
            }
        }
        
        Ok(())
    }
    
    /// Install a single extension
    fn install_extension(&self, extension: &Extension, extensions_dir: &Path) -> Result<()> {
        let ext_dir = extensions_dir.join(&extension.id);
        fs::create_dir_all(&ext_dir)?;
        
        // Write gemini-extension.json
        let config = json!({
            "name": extension.name,
            "version": extension.version,
            "mcpServers": extension.mcp_servers,
        });
        
        let config_path = ext_dir.join("gemini-extension.json");
        let mut file = fs::File::create(&config_path)?;
        file.write_all(serde_json::to_string_pretty(&config)?.as_bytes())?;
        
        // Write context file if present
        if let Some(content) = &extension.context_content {
            // Always write as GEMINI.md for Gemini CLI compatibility
            let context_path = ext_dir.join("GEMINI.md");
            let mut file = fs::File::create(&context_path)?;
            file.write_all(content.as_bytes())?;
        }
        
        println!("  ✓ Installed extension: {}", extension.name);
        
        Ok(())
    }
    
    /// Prepare environment variables
    pub fn prepare_environment(&self, profile: &Profile) -> HashMap<String, String> {
        let mut env_vars = env::vars().collect::<HashMap<_, _>>();
        
        // Add profile-specific environment variables
        for (key, value) in &profile.environment_variables {
            // Expand environment variable references
            let expanded_value = if value.starts_with('$') {
                env::var(&value[1..]).unwrap_or_else(|_| value.clone())
            } else {
                value.clone()
            };
            
            env_vars.insert(key.clone(), expanded_value);
        }
        
        // Add Gemini-specific environment variables
        env_vars.insert(
            "GEMINI_PROFILE".to_string(),
            profile.id.clone(),
        );
        
        env_vars
    }
    
    /// Clean the .gemini directory
    fn clean_gemini_directory(&self, working_dir: &Path) -> Result<()> {
        let gemini_dir = working_dir.join(".gemini");
        
        if !gemini_dir.exists() {
            return Ok(());
        }
        
        let extensions_dir = gemini_dir.join("extensions");
        if extensions_dir.exists() {
            // Remove all extensions
            for entry in fs::read_dir(&extensions_dir)? {
                let entry = entry?;
                let path = entry.path();
                
                if path.is_dir() {
                    let dir_name = path.file_name()
                        .and_then(|n| n.to_str())
                        .unwrap_or("");
                    
                    fs::remove_dir_all(&path)?;
                    println!("  ✓ Removed extension: {}", dir_name);
                }
            }
        }
        
        // Clean other Gemini directories if needed
        // (e.g., cache, logs, etc.)
        
        Ok(())
    }
    
    /// Clean up extensions after Gemini exits
    fn cleanup_extensions(&self, working_dir: &Path) -> Result<()> {
        let extensions_dir = working_dir.join(".gemini").join("extensions");
        
        if extensions_dir.exists() {
            // Only remove extensions we created (those with our naming pattern)
            for entry in fs::read_dir(&extensions_dir)? {
                let entry = entry?;
                let path = entry.path();
                
                if path.is_dir() {
                    // Check if this is one of our managed extensions
                    // (You could check for a marker file or naming pattern)
                    let manifest_path = path.join("gemini-extension.json");
                    if manifest_path.exists() {
                        // This is likely one of our extensions
                        fs::remove_dir_all(&path)?;
                        if let Some(name) = path.file_name().and_then(|n| n.to_str()) {
                            println!("  ✓ Removed extension: {}", name);
                        }
                    }
                }
            }
        }
        
        Ok(())
    }
}

/// Launch a profile in a new terminal window (platform-specific)
#[allow(dead_code)]
pub fn launch_in_terminal(profile: &Profile) -> Result<()> {
    let launcher = Launcher::new();
    
    // For now, we'll just launch in the current terminal
    // In a real implementation, we could detect the platform and launch in a new terminal window
    launcher.launch_with_profile(profile)?;
    
    Ok(())
}

/// Create a launch script for a profile
#[allow(dead_code)]
pub fn create_launch_script(profile: &Profile, output_path: &Path) -> Result<()> {
    let script_content = format!(
        r#"#!/bin/bash
# Launch script for Gemini CLI profile: {}
# Generated by Gemini CLI Manager

# Set profile name
export GEMINI_PROFILE="{}"

# Set environment variables
{}

# Change to working directory
{}

# Launch Gemini CLI
echo "🚀 Launching Gemini CLI with profile: {}"
echo "📂 Working directory: $PWD"
echo ""

gemini
"#,
        profile.name,
        profile.id,
        profile.environment_variables
            .iter()
            .map(|(k, v)| format!("export {}=\"{}\"", k, v))
            .collect::<Vec<_>>()
            .join("\n"),
        profile.working_directory
            .as_ref()
            .map(|dir| format!("cd \"{}\"", dir))
            .unwrap_or_else(|| "# No working directory specified".to_string()),
        profile.display_name(),
    );
    
    let mut file = fs::File::create(output_path)?;
    file.write_all(script_content.as_bytes())?;
    
    // Make the script executable on Unix-like systems
    #[cfg(unix)]
    {
        use std::os::unix::fs::PermissionsExt;
        let mut perms = file.metadata()?.permissions();
        perms.set_mode(0o755);
        fs::set_permissions(output_path, perms)?;
    }
    
    Ok(())
}