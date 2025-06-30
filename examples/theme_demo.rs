use catppuccin::PALETTE;
use gemini_cli_manager::theme::{self, ThemeFlavour};

fn main() {
    println!("Catppuccin Theme Demo for Gemini CLI Manager");
    println!("=============================================\n");
    
    // Test all theme flavours
    let flavours = [
        ThemeFlavour::Latte,
        ThemeFlavour::Frappe,
        ThemeFlavour::Macchiato,
        ThemeFlavour::Mocha,
    ];
    
    for flavour in flavours {
        theme::set_flavour(flavour);
        
        println!("Theme: {:?}", flavour);
        println!("----------------------------------------");
        
        // Print color values as RGB
        let primary = theme::primary();
        let error = theme::error();
        let success = theme::success();
        let text = theme::text_primary();
        let border = theme::border();
        
        println!("Primary:     {:?}", primary);
        println!("Error:       {:?}", error);
        println!("Success:     {:?}", success);
        println!("Text:        {:?}", text);
        println!("Border:      {:?}", border);
        println!();
    }
    
    // Show the actual Catppuccin palette info
    println!("Direct Catppuccin Palette Access:");
    println!("----------------------------------------");
    let mocha = &PALETTE.mocha;
    println!("Mocha Base: RGB({}, {}, {})", 
        mocha.colors.base.rgb.r, 
        mocha.colors.base.rgb.g, 
        mocha.colors.base.rgb.b
    );
    println!("Mocha Text: RGB({}, {}, {})", 
        mocha.colors.text.rgb.r, 
        mocha.colors.text.rgb.g, 
        mocha.colors.text.rgb.b
    );
}