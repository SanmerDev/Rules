use std::fs;
use std::path::{Path, PathBuf};

use clap::Parser;

use crate::rule::{BoxRuleSet, ClashRuleSet};

mod rule;

#[derive(Parser, Debug)]
#[command(disable_colored_help = true)]
struct Args {
    /// Source file (or directory)
    #[arg(value_name = "SOURCE")]
    source: PathBuf,

    /// Output file (or directory)
    #[arg(short, long, value_name = "OUTPUT")]
    output: Option<PathBuf>,
}

impl Args {
    fn with_dir<P: AsRef<Path>>(dir: P, output: P) {
        let dir = Path::new(dir.as_ref());
        for file in dir.read_dir().unwrap().flat_map(|e| e.ok()) {
            let file_name = PathBuf::from(file.file_name());
            let path = file.path();
            if path.is_file() {
                let output = PathBuf::from(output.as_ref());
                let output = output.join(file_name.with_extension("json"));
                Self::with_file(path, output);
            }
        }
    }

    fn with_file<P: AsRef<Path>>(file: P, output: P) {
        let rule_set = match fs::read_to_string(&file) {
            Ok(s) => match ClashRuleSet::from_str(&s) {
                Ok(v) => v,
                Err(_) => return,
            },
            Err(_) => return,
        };

        let rule_set = BoxRuleSet::from_clash(&rule_set);
        if let Err(e) = fs::write(output, rule_set.to_string_pretty()) {
            eprintln!("{}", e)
        }
    }

    fn run(self) {
        if self.source.is_file() {
            let output = self
                .output
                .unwrap_or_else(|| self.source.with_extension("json"));
            Self::with_file(self.source, output)
        } else if self.source.is_dir() {
            if let Some(output) = self.output {
                fs::create_dir_all(&output).ok();
                Self::with_dir(&self.source, &output);
            } else {
                Self::with_dir(&self.source, &self.source)
            }
        }
    }
}

fn main() {
    let args = Args::parse();
    args.run()
}
