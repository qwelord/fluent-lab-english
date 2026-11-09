use serde::Serialize;
use std::collections::{HashMap, HashSet};
use std::io::{self, Read};

#[derive(Serialize)]
struct Analysis { words: usize, sentences: usize, variety: usize, repeated: Vec<String> }

fn analyze(text: &str) -> Analysis {
    let words: Vec<String> = text.split(|c: char| !c.is_alphabetic() && c != '\'' && c != '’')
        .filter(|w| !w.is_empty()).map(|w| w.to_lowercase()).collect();
    let sentences = text.split(|c: char| matches!(c, '.' | '!' | '?'))
        .filter(|s| !s.trim().is_empty()).count();
    let unique: HashSet<&str> = words.iter().map(String::as_str).collect();
    let variety = if words.is_empty() { 0 } else { unique.len() * 100 / words.len() };
    let mut counts: HashMap<&str, usize> = HashMap::new();
    for word in &words { if word.chars().count() > 4 { *counts.entry(word).or_default() += 1; } }
    let mut repeated: Vec<(&str, usize)> = counts.into_iter().filter(|(_, count)| *count >= 3).collect();
    repeated.sort_by(|a,b| b.1.cmp(&a.1).then_with(|| a.0.cmp(b.0)));
    Analysis { words: words.len(), sentences, variety, repeated: repeated.into_iter().take(5).map(|(word,_)| word.to_owned()).collect() }
}

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let mut text = String::new();
    io::stdin().read_to_string(&mut text)?;
    println!("{}", serde_json::to_string(&analyze(&text))?);
    Ok(())
}

#[cfg(test)]
mod tests {
    use super::*;
    #[test] fn empty_text() { let a = analyze(""); assert_eq!((a.words,a.sentences,a.variety),(0,0,0)); }
    #[test] fn repeated_word() { let a = analyze("Focus matters. Focus helps. Focus grows."); assert_eq!(a.words,6); assert_eq!(a.sentences,3); assert_eq!(a.repeated,vec!["focus"]); }
    #[test] fn unicode_words() { let a = analyze("Café café."); assert_eq!(a.words,2); assert_eq!(a.variety,50); }
}
