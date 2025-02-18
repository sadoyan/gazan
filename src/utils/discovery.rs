use dashmap::DashMap;
use futures::channel::mpsc::Sender;
use futures::SinkExt;
use std::fs;
use std::sync::atomic::AtomicUsize;
use std::time::{Duration, Instant};

use async_trait::async_trait;
use notify::event::ModifyKind;
use notify::{Config, Event, EventKind, RecommendedWatcher, RecursiveMode, Watcher};
use std::path::Path;
use tokio::task;

pub struct FromFileProvider {
    pub path: String,
}
pub struct APIUpstreamProvider {
    pub api_url: String,
}
#[async_trait]
pub trait Discovery {
    async fn run(&self, tx: Sender<DashMap<String, (Vec<(String, u16)>, AtomicUsize)>>);
}

#[async_trait]
impl Discovery for APIUpstreamProvider {
    async fn run(&self, mut toreturn: Sender<DashMap<String, (Vec<(String, u16)>, AtomicUsize)>>) {
        loop {
            let dm: DashMap<String, (Vec<(String, u16)>, AtomicUsize)> = DashMap::new();
            dm.insert(
                self.api_url.to_string(),
                (vec![("192.168.1.1".parse().unwrap(), 8000), ("192.168.1.10".parse().unwrap(), 8000)], AtomicUsize::new(0)),
            );
            println!("= = = = = = = = = = = = = = = = = = = = = = = = = = = = = = = = = = = = = = = = = = = ");
            let _ = toreturn.send(dm).await.unwrap();
            println!("= = = = = = = = = = = = = = = = = = = = = = = = = = = = = = = = = = = = = = = = = = = ");
            tokio::time::sleep(Duration::from_secs(20)).await;
        }
    }
}

#[async_trait]
impl Discovery for FromFileProvider {
    async fn run(&self, tx: Sender<DashMap<String, (Vec<(String, u16)>, AtomicUsize)>>) {
        tokio::spawn(watch_file(self.path.clone(), tx.clone()));
    }
}
pub async fn watch_file(fp: String, mut toreturn: Sender<DashMap<String, (Vec<(String, u16)>, AtomicUsize)>>) {
    let file_path = fp.as_str();
    let parent_dir = Path::new(file_path).parent().unwrap(); // Watch directory, not file
    let (local_tx, mut local_rx) = tokio::sync::mpsc::channel::<notify::Result<Event>>(1);

    println!("Watching for changes in {:?}", parent_dir);
    let paths = fs::read_dir(parent_dir).unwrap();
    for path in paths {
        println!("  {}", path.unwrap().path().display())
    }

    let snd = read_upstreams_from_file(file_path);
    let _ = toreturn.send(snd).await.unwrap();

    let _watcher_handle = task::spawn_blocking({
        let parent_dir = parent_dir.to_path_buf(); // Move directory path into the closure
        move || {
            let mut watcher = RecommendedWatcher::new(
                move |res| {
                    let _ = local_tx.blocking_send(res);
                },
                Config::default(),
            )
            .unwrap();
            watcher.watch(&parent_dir, RecursiveMode::Recursive).unwrap();

            loop {
                std::thread::sleep(Duration::from_secs(50));
            }
        }
    });
    // loop {
    //     println!(" ---------------------------------------------------------------- ");
    //     thread::sleep(Duration::from_secs(1));
    // }
    let mut start = Instant::now();

    while let Some(event) = local_rx.recv().await {
        match event {
            Ok(e) => match e.kind {
                EventKind::Modify(ModifyKind::Data(_)) | EventKind::Create(..) | EventKind::Remove(..) => {
                    if e.paths[0].to_str().unwrap().ends_with("conf") {
                        // if start.elapsed() > Duration::from_secs(10) {
                        if start.elapsed() > Duration::from_secs(2) {
                            start = Instant::now();
                            println!("Config File changed :=> {:?}", e);

                            let snd = read_upstreams_from_file(file_path);
                            let _ = toreturn.send(snd).await.unwrap();
                        }
                    }
                }
                _ => (), //println!("*  * * * * * * * * * * * * * * * * * * * * * * * * * * * * * *"),
            },
            Err(e) => println!("Watch error: {:?}", e),
        }
    }
}
fn read_upstreams_from_file(path: &str) -> DashMap<String, (Vec<(String, u16)>, AtomicUsize)> {
    let upstreams = DashMap::new();
    let contents = match fs::read_to_string(path) {
        Ok(data) => data,
        Err(e) => {
            eprintln!("Error reading file: {:?}", e);
            return upstreams;
        }
    };

    for line in contents.lines().filter(|line| !line.trim().is_empty()) {
        let mut parts = line.split_whitespace();

        let Some(hostname) = parts.next() else {
            continue;
        };
        let Some(address) = parts.next() else {
            continue;
        };

        let mut addr_parts = address.split(':');
        let Some(ip) = addr_parts.next() else {
            continue;
        };
        let Some(port_str) = addr_parts.next() else {
            continue;
        };

        let Ok(port) = port_str.parse::<u16>() else {
            continue;
        };
        upstreams
            .entry(hostname.to_string()) // Step 1: Find or create entry
            .or_insert_with(|| (Vec::new(), AtomicUsize::new(0))) // Step 2: Insert if missing
            .0 // Step 3: Access the Vec<(String, u16)>
            .push((ip.to_string(), port)); // Step 4: Append new data
    }

    upstreams
}
