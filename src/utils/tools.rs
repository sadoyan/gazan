use dashmap::DashMap;
use std::any::type_name;
use std::collections::HashSet;
use std::sync::atomic::AtomicUsize;

#[allow(dead_code)]
pub fn print_upstreams(upstreams: &UpstreamsDashMap) {
    for host_entry in upstreams.iter() {
        let hostname = host_entry.key();
        println!("Hostname: {}", hostname);

        for path_entry in host_entry.value().iter() {
            let path = path_entry.key();
            println!("  Path: {}", path);

            for (ip, port, ssl) in path_entry.value().0.clone() {
                println!("   ===> IP: {}, Port: {}, SSL: {}", ip, port, ssl);
            }
        }
    }
}

pub type UpstreamsDashMap = DashMap<String, DashMap<String, (Vec<(String, u16, bool)>, AtomicUsize)>>;
// pub type HeadersList = DashMap<String, Vec<(String, String)>>;
pub type Headers = DashMap<String, DashMap<String, Vec<(String, String)>>>;
// pub type UpstreamMap = DashMap<String, (Vec<(String, u16)>, AtomicUsize)>;

#[allow(dead_code)]
pub fn typeoff<T>(_: T) {
    let to = type_name::<T>();
    println!("{:?}", to);
}

#[allow(dead_code)]
pub fn string_to_bool(val: Option<&str>) -> Option<bool> {
    match val {
        Some(v) => match v {
            "yes" => Some(true),
            "true" => Some(true),
            _ => Some(false),
        },
        None => Some(false),
    }
}

#[allow(dead_code)]
pub fn clone_dashmap(original: &UpstreamsDashMap) -> UpstreamsDashMap {
    let new_map: UpstreamsDashMap = DashMap::new();

    for outer_entry in original.iter() {
        let hostname = outer_entry.key();
        let inner_map = outer_entry.value();

        let new_inner_map = DashMap::new();

        for inner_entry in inner_map.iter() {
            let path = inner_entry.key();
            let (vec, _) = inner_entry.value();
            let new_vec = vec.clone();
            let new_counter = AtomicUsize::new(0);
            new_inner_map.insert(path.clone(), (new_vec, new_counter));
        }
        new_map.insert(hostname.clone(), new_inner_map);
    }
    new_map
}

pub fn clone_dashmap_into(original: &UpstreamsDashMap, cloned: &UpstreamsDashMap) {
    cloned.clear();
    for outer_entry in original.iter() {
        let hostname = outer_entry.key();
        let inner_map = outer_entry.value();

        let new_inner_map = DashMap::new();

        for inner_entry in inner_map.iter() {
            let path = inner_entry.key();
            let (vec, _) = inner_entry.value();
            let new_vec = vec.clone();
            let new_counter = AtomicUsize::new(0);
            new_inner_map.insert(path.clone(), (new_vec, new_counter));
        }
        cloned.insert(hostname.clone(), new_inner_map);
    }
}

pub fn compare_dashmaps(map1: &UpstreamsDashMap, map2: &UpstreamsDashMap) -> bool {
    let keys1: HashSet<_> = map1.iter().map(|entry| entry.key().clone()).collect();
    let keys2: HashSet<_> = map2.iter().map(|entry| entry.key().clone()).collect();
    if keys1 != keys2 {
        return false;
    }
    for entry1 in map1.iter() {
        let hostname = entry1.key();
        let inner_map1 = entry1.value();
        let Some(inner_map2) = map2.get(hostname) else {
            return false;
        };
        let inner_keys1: HashSet<_> = inner_map1.iter().map(|e| e.key().clone()).collect();
        let inner_keys2: HashSet<_> = inner_map2.iter().map(|e| e.key().clone()).collect();
        if inner_keys1 != inner_keys2 {
            return false;
        }
        for path_entry in inner_map1.iter() {
            let path = path_entry.key();
            let (vec1, _counter1) = path_entry.value();
            let Some(entry2) = inner_map2.get(path) else {
                return false; // Path exists in map1 but not in map2
            };
            let (vec2, _counter2) = entry2.value();
            let set1: HashSet<_> = vec1.iter().collect();
            let set2: HashSet<_> = vec2.iter().collect();
            if set1 != set2 {
                return false;
            }
        }
    }
    true
}
