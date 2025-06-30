use gemini_cli_manager::theme::{self, ThemeFlavour};
use catppuccin::PALETTE;

fn main() {
    // Set Mocha theme
    theme::set_flavour(ThemeFlavour::Mocha);
    
    println!("Catppuccin Mocha Color Values:");
    println!("===============================");
    
    let mocha = &PALETTE.mocha.colors;
    
    // Background colors
    println!("\nBackground Colors:");
    println!("base (background):    {:?} - darkest", mocha.base.rgb);
    println!("mantle:               {:?}", mocha.mantle.rgb);
    println!("crust:                {:?}", mocha.crust.rgb);
    
    // Surface colors
    println!("\nSurface Colors:");
    println!("surface0:             {:?}", mocha.surface0.rgb);
    println!("surface1:             {:?}", mocha.surface1.rgb);
    println!("surface2:             {:?} - current selection", mocha.surface2.rgb);
    
    // Text colors
    println!("\nText Colors:");
    println!("text:                 {:?} - primary", mocha.text.rgb);
    println!("subtext1:             {:?} - secondary", mocha.subtext1.rgb);
    println!("subtext0:             {:?} - muted", mocha.subtext0.rgb);
    println!("overlay2:             {:?}", mocha.overlay2.rgb);
    println!("overlay1:             {:?}", mocha.overlay1.rgb);
    println!("overlay0:             {:?} - disabled", mocha.overlay0.rgb);
    
    // Accent colors
    println!("\nAccent Colors:");
    println!("blue (primary):       {:?}", mocha.blue.rgb);
    println!("yellow (highlight):   {:?}", mocha.yellow.rgb);
    println!("green (success):      {:?}", mocha.green.rgb);
    println!("red (error):          {:?}", mocha.red.rgb);
    println!("peach (warning):      {:?}", mocha.peach.rgb);
    println!("sky (info):           {:?}", mocha.sky.rgb);
    println!("mauve (secondary):    {:?}", mocha.mauve.rgb);
    println!("pink (accent):        {:?}", mocha.pink.rgb);
    
    // Calculate contrast ratios for key combinations
    println!("\nContrast Analysis:");
    println!("==================");
    
    fn luminance(rgb: &catppuccin::Rgb) -> f64 {
        let r = rgb.r as f64 / 255.0;
        let g = rgb.g as f64 / 255.0;
        let b = rgb.b as f64 / 255.0;
        
        let r = if r <= 0.03928 { r / 12.92 } else { ((r + 0.055) / 1.055).powf(2.4) };
        let g = if g <= 0.03928 { g / 12.92 } else { ((g + 0.055) / 1.055).powf(2.4) };
        let b = if b <= 0.03928 { b / 12.92 } else { ((b + 0.055) / 1.055).powf(2.4) };
        
        0.2126 * r + 0.7152 * g + 0.0722 * b
    }
    
    fn contrast_ratio(rgb1: &catppuccin::Rgb, rgb2: &catppuccin::Rgb) -> f64 {
        let l1 = luminance(rgb1);
        let l2 = luminance(rgb2);
        let lighter = l1.max(l2);
        let darker = l1.min(l2);
        (lighter + 0.05) / (darker + 0.05)
    }
    
    // Test key combinations
    println!("text on base:         {:.2}:1", contrast_ratio(&mocha.text.rgb, &mocha.base.rgb));
    println!("text on surface2:     {:.2}:1 (current selection bg)", contrast_ratio(&mocha.text.rgb, &mocha.surface2.rgb));
    println!("subtext0 on base:     {:.2}:1 (muted text)", contrast_ratio(&mocha.subtext0.rgb, &mocha.base.rgb));
    println!("yellow on surface2:   {:.2}:1 (highlight on selection)", contrast_ratio(&mocha.yellow.rgb, &mocha.surface2.rgb));
    
    println!("\nRecommended minimum contrast ratios:");
    println!("Normal text: 4.5:1");
    println!("Large text: 3:1");
    println!("UI components: 3:1");
}